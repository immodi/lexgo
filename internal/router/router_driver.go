package router

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
)

type RouterVmDriver struct {
	Router *Router
}

func (router *RouterVmDriver) RegisterLuaMethodHandler(fn framework.RouterServerHandler, path string, method string, mws []*vm.LuaFunction) {
	router.Router.AppendRoute(&Handler{
		Pattern: path,
		Handler: fn,
		Params:  map[string]string{},
		Method:  framework.HTTPMethod(method),
		Mws:     mws,
	})
}

func (router *RouterVmDriver) RegisterLuaErrorHandler(fn framework.RouterServerHandler) {
	router.Router.ServerErrorFunc = fn
}

func (router *RouterVmDriver) RegisterLuaNotFoundHandler(fn framework.RouterServerHandler) {
	router.Router.NotFoundFunc = fn
}

func (router *RouterVmDriver) RegisterLuaMiddleware(fn framework.RouterServerHandler) {
	// router.Router.MiddleWares = append(router.Router.MiddleWares, fn)
}

func (router *RouterVmDriver) GetAllRegistredRoutes() map[string][]string {
	var routes map[string][]string = map[string][]string{}
	for route := range router.Router.Routes {
		routes[route.Path] = append(routes[route.Path], string(route.Method))
	}
	return routes
}
