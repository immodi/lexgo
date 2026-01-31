package middlewares

import lua "github.com/yuin/gopher-lua"

type Driver interface {
	ExecuteFinal(
		fn *lua.LFunction,
		req Request,
		res Response,
	)
	HandleError(msg string, res any)
	LuaState() *lua.LState
}
