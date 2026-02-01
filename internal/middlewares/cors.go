package middlewares

import (
	lua "github.com/yuin/gopher-lua"
)

func DefaultLuaCORS(L *lua.LState) *lua.LFunction {
	return L.NewFunction(func(L *lua.LState) int {
		req := L.CheckTable(1)
		res := L.CheckTable(2)
		next := L.CheckFunction(3)

		// Set CORS headers
		L.CallByParam(lua.P{
			Fn:      L.GetField(res, "setHeader"),
			NRet:    0,
			Protect: true,
		}, lua.LString("Access-Control-Allow-Origin"), lua.LString("*"))

		L.CallByParam(lua.P{
			Fn:      L.GetField(res, "setHeader"),
			NRet:    0,
			Protect: true,
		}, lua.LString("Access-Control-Allow-Methods"), lua.LString("GET, POST, PUT, PATCH, DELETE, OPTIONS"))

		L.CallByParam(lua.P{
			Fn:      L.GetField(res, "setHeader"),
			NRet:    0,
			Protect: true,
		}, lua.LString("Access-Control-Allow-Headers"), lua.LString("Content-Type, Authorization"))

		// Automatically respond to OPTIONS requests
		method := L.GetField(req, "method").String()
		if method == "OPTIONS" {
			L.CallByParam(lua.P{
				Fn:      L.GetField(res, "status"),
				NRet:    0,
				Protect: true,
			}, lua.LNumber(200))

			L.CallByParam(lua.P{
				Fn:      L.GetField(res, "raw"),
				NRet:    0,
				Protect: true,
			}, lua.LString(""))
			return 0
		}

		// Call next middleware
		L.CallByParam(lua.P{
			Fn:      next,
			NRet:    0,
			Protect: true,
		})
		return 0
	})
}
