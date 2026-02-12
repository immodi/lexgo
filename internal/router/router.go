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
}

type Router struct {
	LuaVm           vm.LVm
	Routes          map[framework.HTTPRoute]*Handler
	MiddleWares     []*vm.LuaFunction
	NotFoundFunc    *vm.LuaFunction
	ServerErrorFunc *vm.LuaFunction
}

func MakeRouter(luaVm vm.LVm) (*Router, *RouterVmDriver) {
	router := &Router{
		Routes:      make(map[framework.HTTPRoute]*Handler),
		LuaVm:       luaVm,
		MiddleWares: make([]*vm.LuaFunction, 0),
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
	incomingParts := strings.Split(strings.Trim(incoming.Path, "/"), "/")
	for route, handler := range router.Routes {
		if route.Method != incoming.Method {
			continue
		}

		storedParts := strings.Split(strings.Trim(route.Path, "/"), "/")

		if len(incomingParts) != len(storedParts) {
			continue
		}

		params := make(map[string]string)
		matched := true

		for i := range incomingParts {
			storedSegment := storedParts[i]
			incomingSegment := incomingParts[i]

			if param, ok := strings.CutPrefix(storedSegment, ":"); ok {
				params[param] = incomingSegment
				continue
			}

			if incomingSegment != storedSegment {
				matched = false
				break
			}
		}

		if matched {
			return &Handler{
				Pattern: route.Path,
				Handler: handler.Handler,
				Params:  params,
			}
		}
	}

	return nil
}
