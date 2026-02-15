package router

import (
	"net/http"
	"strings"

	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
)

type Handler struct {
	Pattern string
	Handler *vm.LuaFunction
	Params  map[string]string
	Method  framework.HTTPMethod
}

type RouterTreeNode struct {
	name string

	staticChildren map[string]*RouterTreeNode // literal segments
	param          *RouterTreeNode            // :id
	wildcard       *RouterTreeNode            // *path

	handler *Handler
}

type Router struct {
	LuaVm           vm.LVm
	Routes          map[framework.HTTPRoute]*Handler
	MiddleWares     []*vm.LuaFunction
	NotFoundFunc    *vm.LuaFunction
	ServerErrorFunc *vm.LuaFunction
	RootNode        *RouterTreeNode
}

func MakeRouter(luaVm vm.LVm) (*Router, *RouterVmDriver) {
	router := &Router{
		Routes:      make(map[framework.HTTPRoute]*Handler),
		LuaVm:       luaVm,
		MiddleWares: make([]*vm.LuaFunction, 0),
		RootNode:    &RouterTreeNode{staticChildren: map[string]*RouterTreeNode{}},
	}

	routerDriver := &RouterVmDriver{Router: router}
	return router, routerDriver
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route := framework.HTTPRoute{Path: req.URL.Path, Method: framework.HTTPMethod(req.Method)}
	handler := router.matchRoute(route)
	router.PrintTree()

	if handler == nil {
		handler = &Handler{
			Pattern: route.Path,
			Handler: router.NotFoundFunc,
			Params:  map[string]string{},
			Method:  framework.GET,
		}
	}

	luaReq := &framework.LuaRequest{HttpRequest: req, LuaVm: router.LuaVm, Params: handler.Params}
	luaRes := framework.ConstructResponse(w, router.LuaVm)

	ctx := framework.NewMiddlewaresContext(
		&MiddlewareVmDriver{router, luaReq, luaRes},
		handler.Handler,
	)

	framework.ExecuteMiddlewares(ctx, router.MiddleWares)
}

func (router *Router) matchRoute(incoming framework.HTTPRoute) *Handler {
	var currentNode *RouterTreeNode = router.RootNode
	var wildCardNode *RouterTreeNode = nil
	params := make(map[string]string)
	incomingParts := strings.Split(strings.Trim(incoming.Path, "/"), "/")

	assignWildNode := func(node *RouterTreeNode, index int) {
		if node != nil && node.wildcard != nil {
			wildCardNode = node.wildcard
			params["*"] = strings.Join(incomingParts[index:], "/")
		}
	}

	for i, token := range incomingParts {
		staticNode, ok := currentNode.staticChildren[token]
		if ok {
			currentNode = staticNode
			assignWildNode(currentNode, i)
			continue
		}

		if currentNode.param != nil && currentNode.param.handler != nil {
			currentNode = currentNode.param
			assignWildNode(currentNode, i)
			params[currentNode.name] = token
			continue
		}

		currentNode = wildCardNode
		break
	}

	if currentNode != nil && currentNode.handler != nil && currentNode.handler.Method == incoming.Method {
		return &Handler{
			Pattern: currentNode.handler.Pattern,
			Handler: currentNode.handler.Handler,
			Params:  params,
			Method:  currentNode.handler.Method,
		}
	}

	if wildCardNode != nil && wildCardNode.handler != nil && wildCardNode.handler.Method == incoming.Method {
		return &Handler{
			Pattern: wildCardNode.handler.Pattern,
			Handler: wildCardNode.handler.Handler,
			Params:  params,
			Method:  wildCardNode.handler.Method,
		}
	}

	return nil
}

func (router *Router) ConstructRouterNode(handler *Handler) {
	// the url words after base route [ /test/route --> [test, route] ]
	incomingParts := strings.Split(strings.Trim(handler.Pattern, "/"), "/")
	var node *RouterTreeNode = router.RootNode

	for i, token := range incomingParts {
		var finalHandler *Handler = nil
		if i+1 == len(incomingParts) {
			finalHandler = handler
		}

		switch {
		case strings.Contains(token, ":"):
			paramName, _ := strings.CutPrefix(token, ":")
			node.param = &RouterTreeNode{
				name:           paramName,
				handler:        finalHandler,
				staticChildren: map[string]*RouterTreeNode{}}
			node = node.param
			continue

		case strings.Contains(token, "*"):
			node.wildcard = &RouterTreeNode{
				name:           "*",
				handler:        finalHandler,
				staticChildren: map[string]*RouterTreeNode{}}

			node = node.wildcard
			return

		default:
			currentNode, ok := node.staticChildren[token]
			if !ok {
				currentNode = &RouterTreeNode{
					name:           token,
					handler:        finalHandler,
					staticChildren: map[string]*RouterTreeNode{}}

				node.staticChildren[token] = currentNode
			}
			node = currentNode
			continue

		}
	}
}
