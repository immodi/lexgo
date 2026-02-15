package router

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
)

type RouterVmDriver struct {
	Router *Router
}

func (router *RouterVmDriver) RegisterLuaMethodHandler(fn *vm.LuaFunction, path string, method string) {
	router.Router.ConstructRouterNode(&Handler{
		Pattern: path,
		Handler: fn,
		Params:  map[string]string{},
		Method:  framework.HTTPMethod(method),
	})
	// router.Router.Routes[framework.HTTPRoute{Path: path, Method: framework.HTTPMethod(method)}] =
}

func (router *RouterVmDriver) RegisterLuaErrorHandler(fn *vm.LuaFunction) {
	router.Router.ServerErrorFunc = fn
}

func (router *RouterVmDriver) RegisterLuaNotFoundHandler(fn *vm.LuaFunction) {
	router.Router.NotFoundFunc = fn
}

func (router *RouterVmDriver) RegisterLuaMiddleware(fn *vm.LuaFunction) {
	router.Router.MiddleWares = append(router.Router.MiddleWares, fn)
}

func (router *RouterVmDriver) GetAllRegistredRoutes() map[string][]string {
	var routes map[string][]string = map[string][]string{}
	for route := range router.Router.Routes {
		routes[route.Path] = append(routes[route.Path], string(route.Method))
	}
	return routes
}
