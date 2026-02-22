package router

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
	"net/http"
)

type LuaHandler struct {
	luaVm vm.LVm
	luaFn *vm.LuaFunction
}

func (h *LuaHandler) Handle(
	luaReq *framework.LuaRequest,
	luaRes *framework.LuaResponse,
	args ...vm.LuaValue,
) error {

	values := []vm.LuaValue{
		luaReq.MakeLuaRequest(),
		luaRes.MakeLuaResponse(),
	}

	values = append(values, args...)

	return h.luaVm.RunFunction(
		h.luaFn,
		values...,
	)
}

type GoHandler struct {
	fn func(w http.ResponseWriter, r *http.Request)
}

func (h *GoHandler) Handle(
	luaReq *framework.LuaRequest,
	luaRes *framework.LuaResponse,
	_ ...vm.LuaValue,
) error {
	h.fn(luaRes.HttpWriter, luaReq.HttpRequest)
	return nil
}

type RouterVmDriver struct {
	Router *Router
}

func (router *RouterVmDriver) MakeLuaHandler(luaVm vm.LVm, fn *vm.LuaFunction) framework.RouterServerHandler {
	return &LuaHandler{luaVm: luaVm, luaFn: fn}
}

func (router *RouterVmDriver) MakeGoHandler(fn func(w http.ResponseWriter, r *http.Request)) framework.RouterServerHandler {
	return &GoHandler{fn: fn}
}

func (router *RouterVmDriver) RegisterHandler(fn framework.RouterServerHandler, path string, method string, mws []framework.RouterServerHandler) {
	router.Router.AppendRoute(
		&Handler{
			Pattern: path,
			Handler: fn,
			Params:  map[string]string{},
			Method:  framework.HTTPMethod(method),
			Mws:     mws,
		},
	)
}

func (router *RouterVmDriver) RegisterErrorHandler(fn framework.RouterServerHandler) {
	router.Router.ServerErrorFunc = fn
}

func (router *RouterVmDriver) RegisterNotFoundHandler(fn framework.RouterServerHandler) {
	router.Router.NotFoundFunc = fn
}

func (router *RouterVmDriver) RegisterMiddleware(fn framework.RouterServerHandler) {
	router.Router.Mws = append(router.Router.Mws, fn)
}

func (router *RouterVmDriver) GetAllRegistredRoutes() map[string][]string {
	var routes map[string][]string = map[string][]string{}
	for route := range router.Router.Routes {
		routes[route.Path] = append(routes[route.Path], string(route.Method))
	}
	return routes
}
