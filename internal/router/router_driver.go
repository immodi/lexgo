package router

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
)

type RouterVmDriver struct {
	Router *Router
}

func (router *RouterVmDriver) RegisterLuaMethodHandler(fn *vm.LuaFunction, path string, method string) {
	router.Router.Routes[framework.HTTPRoute{Path: path, Method: framework.HTTPMethod(method)}] = &Handler{
		Pattern: path,
		Handler: fn,
		Params:  map[string]string{},
	}
}

func (router *RouterVmDriver) ResgisterLuaErrorHandler(fn *vm.LuaFunction) {
	router.Router.ServerErrorFunc = fn
}

func (router *RouterVmDriver) ResgisterLuaNotFoundHandler(fn *vm.LuaFunction) {
	router.Router.NotFoundFunc = fn
}

func (router *RouterVmDriver) RegisterLuaMiddleware(fn *vm.LuaFunction) {
	router.Router.MiddleWares = append(router.Router.MiddleWares, fn)
}
