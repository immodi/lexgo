package router

import (
	"fmt"
	"net/http"
	"strings"

	"immodi/lexgo/internal/framework"
)

type Handler struct {
	Pattern string
	Params  map[string]string
	Handler framework.RouterServerHandler
	Method  framework.HTTPMethod
	Mws     []framework.RouterServerHandler
}

type RouterTreeNode struct {
	name           string
	staticChildren map[string]*RouterTreeNode
	param          *RouterTreeNode
	wildcard       *RouterTreeNode
	handler        *Handler
}

type Router struct {
	Routes          map[framework.HTTPRoute]*Handler
	NotFoundFunc    framework.RouterServerHandler
	ServerErrorFunc framework.RouterServerHandler
	Mws             []framework.RouterServerHandler
	RootNode        *RouterTreeNode
}

func MakeRouter() (*Router, *RouterVmDriver) {
	router := &Router{
		Routes:   make(map[framework.HTTPRoute]*Handler),
		RootNode: &RouterTreeNode{staticChildren: map[string]*RouterTreeNode{}},
	}

	routerDriver := &RouterVmDriver{Router: router}
	return router, routerDriver
}

func (router *Router) GetHTTPRoute(req *http.Request) *framework.HTTPRoute {
	route := framework.HTTPRoute{Path: req.URL.Path, Method: framework.HTTPMethod(req.Method)}
	return &route
}

func (router *Router) Match(incoming *framework.HTTPRoute) (
	fn *Handler,
	notFoundFn framework.RouterServerHandler,
	serverErrorFn framework.RouterServerHandler,
) {
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
			Mws:     append(currentNode.handler.Mws, router.Mws...),
		}, router.NotFoundFunc, router.ServerErrorFunc
	}

	if wildCardNode != nil && wildCardNode.handler != nil && wildCardNode.handler.Method == incoming.Method {
		return &Handler{
			Pattern: wildCardNode.handler.Pattern,
			Handler: wildCardNode.handler.Handler,
			Params:  params,
			Method:  wildCardNode.handler.Method,
			Mws:     append(wildCardNode.handler.Mws, router.Mws...),
		}, router.NotFoundFunc, router.ServerErrorFunc
	}

	return &Handler{
		Pattern: incoming.Path,
		Handler: router.NotFoundFunc,
		Params:  map[string]string{},
		Method:  incoming.Method,
		Mws:     make([]framework.RouterServerHandler, 0),
	}, router.NotFoundFunc, router.ServerErrorFunc
}

func (router *Router) AppendRoute(handler *Handler) {
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
