package framework

import "immodi/lexgo/internal/vm"

type AppData struct {
	Port int32
}

func RegisterFramework(luaVm vm.LVm, routerDriver RouterDriver) (*AppData, error) {
	tbl := luaVm.NewTable()
	data := &AppData{}

	luaVm.SetGlobal("lexgo", tbl)
	newFn := luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		app := ResgisterRouter(l, routerDriver)

		listenFn := l.NewFunction(func(l vm.LVm) vm.LuaValue {
			port, err := l.CheckNumber(1)
			if err != nil {
				l.Error(err.Error())
				return nil
			}
			app.SetField("_port", vm.LuaNumber(port))
			data.Port = int32(port)
			return nil
		})
		app.SetField("listen", listenFn)
		return app
	})

	tbl.SetField("new", newFn)
	RegisterDefaultMiddlewares(luaVm, tbl, routerDriver.GetAllRegistredRoutes)

	return data, nil
}
