package router

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
	"net/http"
	"strings"
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

func (router *RouterVmDriver) MakeLuaHandler(luaVm vm.LVm, fn *vm.LuaFunction) framework.RouterHandler {
	return &LuaHandler{luaVm: luaVm, luaFn: fn}
}

func (router *RouterVmDriver) MakeGoHandler(fn func(w http.ResponseWriter, r *http.Request)) framework.RouterHandler {
	return &GoHandler{fn: fn}
}

func (router *RouterVmDriver) RegisterHandler(fn framework.RouterHandler, path string, method string, mws []framework.RouterHandler) {
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

func (router *RouterVmDriver) RegisterErrorHandler(fn framework.RouterHandler) {
	router.Router.ServerErrorFunc = fn
}

func (router *RouterVmDriver) RegisterNotFoundHandler(fn framework.RouterHandler) {
	router.Router.NotFoundFunc = fn
}

func (router *RouterVmDriver) RegisterMiddleware(fn framework.RouterHandler) {
	router.Router.Mws = append(router.Router.Mws, fn)
}

func (router *RouterVmDriver) GetAllRegistredRoutes() map[string][]string {
	routes := make(map[string][]string)

	var traverse func(node *RouterTreeNode, pathParts []string)
	traverse = func(node *RouterTreeNode, pathParts []string) {
		if node == nil {
			return
		}

		// If this node has a handler, record it
		if node.handler != nil {
			fullPath := "/" + strings.Join(pathParts, "/")
			method := string(node.handler.Method)
			routes[fullPath] = append(routes[fullPath], method)
		}

		// Traverse static children
		for _, child := range node.staticChildren {
			traverse(child, append(pathParts, child.name))
		}

		// Traverse param child
		if node.param != nil {
			traverse(node.param, append(pathParts, ":"+node.param.name))
		}

		// Traverse wildcard child
		if node.wildcard != nil {
			traverse(node.wildcard, append(pathParts, "*"+node.wildcard.name))
		}
	}

	traverse(router.Router.RootNode, []string{})
	return routes
}

func (router *RouterVmDriver) ClearRoutes() {
	router.Router.ClearRoutes()
}
