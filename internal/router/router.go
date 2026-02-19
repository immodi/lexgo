package router

import (
	"fmt"
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
	Mws     []*vm.LuaFunction
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

	if handler == nil {
		handler = &Handler{
			Pattern: route.Path,
			Handler: router.NotFoundFunc,
			Params:  map[string]string{},
			Method:  framework.GET,
			Mws:     make([]*vm.LuaFunction, 0),
		}
	}

	luaReq := &framework.LuaRequest{HttpRequest: req, LuaVm: router.LuaVm, Params: handler.Params}
	luaRes := framework.ConstructResponse(w, router.LuaVm)

	ctx := framework.NewMiddlewaresContext(
		&MiddlewareVmDriver{router, luaReq, luaRes},
		handler.Handler,
		handler.Mws,
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
			params["*"] = strings.Join(incomingParts[index+1:], "/")
		}
	}

	for i, token := range incomingParts {
		staticNode, ok := currentNode.staticChildren[fmt.Sprintf("%s:%s", incoming.Method, token)]
		if ok {
			currentNode = staticNode
			assignWildNode(currentNode, i)
			continue
		}

		if currentNode.param != nil {
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
			Mws:     currentNode.handler.Mws,
		}
	}

	if wildCardNode != nil && wildCardNode.handler != nil && wildCardNode.handler.Method == incoming.Method {
		return &Handler{
			Pattern: wildCardNode.handler.Pattern,
			Handler: wildCardNode.handler.Handler,
			Params:  params,
			Method:  wildCardNode.handler.Method,
			Mws:     currentNode.handler.Mws,
		}
	}

	return nil
}

// registers a lua router handler into the go router
func (router *Router) ConstructRouterNode(handler *Handler) {
	parts := strings.Split(strings.Trim(handler.Pattern, "/"), "/")
	node := router.RootNode

	for i, part := range parts {
		isLast := i+1 == len(parts)
		var finalHandler *Handler
		if isLast {
			finalHandler = handler
		}

		switch {
		case strings.HasPrefix(part, ":"):
			paramName := part[1:]
			if node.param == nil {
				node.param = &RouterTreeNode{
					name:           paramName,
					staticChildren: map[string]*RouterTreeNode{},
				}
			}
			if finalHandler != nil {
				node.param.handler = finalHandler
			}
			node = node.param

		case strings.HasPrefix(part, "*"):
			node.wildcard = &RouterTreeNode{
				name:           "*",
				handler:        finalHandler,
				staticChildren: map[string]*RouterTreeNode{},
			}
			node = node.wildcard
			return

		default:
			key := fmt.Sprintf("%s:%s", handler.Method, part)
			child, exists := node.staticChildren[key]
			if !exists {
				child = &RouterTreeNode{
					name:           part,
					handler:        finalHandler,
					staticChildren: map[string]*RouterTreeNode{},
				}
				node.staticChildren[key] = child
			}
			node = child
		}
	}
}
