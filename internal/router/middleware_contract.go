package router

import (
	"immodi/lexgo/internal/middlewares"

	lua "github.com/yuin/gopher-lua"
)

func (r *Router) LuaState() *lua.LState {
	return r.LuaVm.L
}

func (r *Router) ExecuteFinal(
	fn *lua.LFunction,
	req middlewares.Request,
	res middlewares.Response,
) {
	r.executeLuaHandler(
		fn,
		req.(*LuaRequest),
		res.(*LuaResponse),
	)
}

func (r *Router) HandleError(msg string, res any) {
	r.handleServerError(msg, res.(*LuaResponse))
}
