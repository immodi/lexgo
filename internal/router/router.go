package router

import (
	"net/http"
	"strings"

	"immodi/lexgo/internal/vm"

	lua "github.com/yuin/gopher-lua"
)

type Router struct {
	LuaVm           *vm.LuaVm
	Routes          map[vm.HTTPRoute]*Handler
	MiddleWares     []*lua.LFunction
	NotFoundFunc    *lua.LFunction
	ServerErrorFunc *lua.LFunction
}

func MakeRouter(luaVm *vm.LuaVm) (*Router, *RouterVmDriver) {
	router := &Router{
		Routes:      make(map[vm.HTTPRoute]*Handler),
		LuaVm:       luaVm,
		MiddleWares: make([]*lua.LFunction, 0),
	}

	routerDriver := &RouterVmDriver{Router: router}
	return router, routerDriver
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route := vm.HTTPRoute{Path: req.URL.Path, Method: vm.HTTPMethod(req.Method)}
	handler := router.matchRoute(route)

	if handler == nil {
		handler = &Handler{
			Pattern: route.Path,
			Handler: router.NotFoundFunc,
			Params:  map[string]string{},
		}
	}

	luaReq := &vm.LuaRequest{HttpRequest: req, LuaVm: router.LuaVm, Params: handler.Params}
	luaRes := &vm.LuaResponse{HttpWriter: w, LuaVm: router.LuaVm, Written: false}

	if len(router.MiddleWares) > 0 {
		ctx := vm.NewMiddlewaresContext(
			&MiddlewareVmDriver{router, luaReq, luaRes},
			handler.Handler,
		)

		vm.ExecuteMiddlewares(ctx, router.MiddleWares)
	} else {
		vm.ExecuteLuaHandler(router.LuaVm.L, router.ServerErrorFunc, handler.Handler, luaReq, luaRes)
	}
}

func (router *Router) matchRoute(incoming vm.HTTPRoute) *Handler {
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
