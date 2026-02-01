package router

import (
	"immodi/lexgo/internal/vm"

	lua "github.com/yuin/gopher-lua"
)

type RouterVmDriver struct {
	Router *Router
}

func (router *RouterVmDriver) RegisterLuaMethodHandler(fn *lua.LFunction, path string, method string) {
	router.Router.Routes[vm.HTTPHandler{Path: path, Method: vm.HTTPMethod(method)}] = fn
}

func (router *RouterVmDriver) ResgisterLuaErrorHandler(fn *lua.LFunction) {
	router.Router.ServerErrorFunc = fn
}

func (router *RouterVmDriver) ResgisterLuaNotFoundHandler(fn *lua.LFunction) {
	router.Router.NotFoundFunc = fn
}

func (router *RouterVmDriver) RegisterLuaMiddleware(fn *lua.LFunction) {
	router.Router.MiddleWares = append(router.Router.MiddleWares, fn)
}
