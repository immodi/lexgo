package router

import (
	"net/http"

	"immodi/lexgo/internal/vm"

	lua "github.com/yuin/gopher-lua"
)

type Router struct {
	LuaVm           *vm.LuaVm
	Routes          map[vm.HTTPHandler]*lua.LFunction
	MiddleWares     []*lua.LFunction
	NotFoundFunc    *lua.LFunction
	ServerErrorFunc *lua.LFunction
}

func MakeRouter(luaVm *vm.LuaVm) (*Router, *RouterVmDriver) {
	router := &Router{
		Routes:      make(map[vm.HTTPHandler]*lua.LFunction),
		LuaVm:       luaVm,
		MiddleWares: make([]*lua.LFunction, 0),
	}

	routerDriver := &RouterVmDriver{Router: router}
	return router, routerDriver
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := vm.HTTPHandler{Path: req.URL.Path, Method: vm.HTTPMethod(req.Method)}
	fn, ok := router.Routes[handler]

	luaReq := &vm.LuaRequest{HttpRequest: req, LuaVm: router.LuaVm}
	luaRes := &vm.LuaResponse{HttpWriter: w, LuaVm: router.LuaVm, Written: false}

	if !ok {
		fn = router.NotFoundFunc
	}

	if len(router.MiddleWares) > 0 {
		ctx := vm.NewMiddlewaresContext(
			&MiddlewareVmDriver{router, luaReq, luaRes},
			fn,
		)

		vm.ExecuteMiddlewares(ctx, router.MiddleWares)
	} else {
		vm.ExecuteLuaHandler(router.LuaVm.L, router.ServerErrorFunc, fn, luaReq, luaRes)
	}
}
