package vm

type AppData struct {
	Port int32
}

func RegisterFramework(luaVm LVm, routerDriver RouterDriver) (*AppData, error) {
	tbl := luaVm.NewTable()
	data := &AppData{}

	luaVm.SetGlobal("lexgo", tbl)
	newFn := luaVm.NewFunction(func(l LVm) LuaValue {
		app := ResgisterRouter(l, routerDriver)

		listenFn := l.NewFunction(func(l LVm) LuaValue {
			port, err := l.CheckNumber(1)
			if err != nil {
				l.Error(err.Error())
				return nil
			}
			// L.SetField(app, "_port", lua.LNumber(port))
			app.SetField("_port", LuaNumber(port))
			data.Port = int32(port)
			return nil
		})
		app.SetField("listen", listenFn)
		return app
	})

	tbl.SetField("new", newFn)
	RegisterDefaultMiddlewares(luaVm, tbl)

	return data, nil
}
