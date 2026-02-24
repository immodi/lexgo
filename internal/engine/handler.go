package engine

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/router"
	"immodi/lexgo/internal/vm"
	"net/http"
)

type RouterInterface interface {
	Match(incoming *framework.HTTPRoute) (
		fn *router.Handler,
		notFoundFn framework.RouterHandler,
		serverErrorFn framework.RouterHandler,
	)
	GetHTTPRoute(req *http.Request) *framework.HTTPRoute
	AppendRoute(handler *router.Handler)
}

type Engine struct {
	LuaVm  vm.LVm
	Router RouterInterface
}

func MakeEngine(router RouterInterface, vm vm.LVm) http.Handler {
	return &Engine{LuaVm: vm, Router: router}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route := engine.Router.GetHTTPRoute(req)
	handler, notFoundFn, serverErrorFn := engine.Router.Match(route)

	luaReq := framework.ConstructRequest(req, engine.LuaVm, handler.Params)
	luaRes := framework.ConstructResponse(w, engine.LuaVm)

	ctx := framework.NewExecutionContext(
		&ExecutionDriver{
			engine.LuaVm,
			luaReq,
			luaRes,
		},
		handler.Handler,
		notFoundFn,
		serverErrorFn,
		handler.Mws,
	)

	framework.Execute(ctx)
}
