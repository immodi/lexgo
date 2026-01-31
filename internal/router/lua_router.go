package router

import lua "github.com/yuin/gopher-lua"

func (router *Router) RegisterFramework() {
	L := router.LuaVm.L

	tbl := L.NewTable()
	L.SetGlobal("lexgo", tbl)

	L.SetField(tbl, "new", L.NewFunction(func(l *lua.LState) int {
		app := router.MakeAppTable()
		l.Push(app)
		return 1
	}))
}

func (router *Router) MakeAppTable() *lua.LTable {
	L := router.LuaVm.L

	appTbl := L.NewTable()

	L.SetField(appTbl, "get", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Handler{Path: path, Method: GET}] = fn
		return 0
	}))

	L.SetField(appTbl, "post", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Handler{Path: path, Method: POST}] = fn
		return 0
	}))

	L.SetField(appTbl, "put", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Handler{Path: path, Method: PUT}] = fn
		return 0
	}))

	L.SetField(appTbl, "delete", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Handler{Path: path, Method: DELETE}] = fn
		return 0
	}))

	L.SetField(appTbl, "patch", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Handler{Path: path, Method: PATCH}] = fn
		return 0
	}))

	L.SetField(appTbl, "options", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[Handler{Path: path, Method: OPTIONS}] = fn
		return 0
	}))

	L.SetField(appTbl, "notFound", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		router.NotFoundFunc = fn
		return 0
	}))

	L.SetField(appTbl, "error", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		router.ServerErrorFunc = fn
		return 0
	}))

	L.SetField(appTbl, "use", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		router.MiddleWares = append(router.MiddleWares, fn)
		return 0
	}))

	return appTbl
}
