package framework

import "immodi/lexgo/internal/vm"

func createNewAppFunction(
	luaVm vm.LVm,
	data *AppData,
	routerDriver RouterDriver,
	restartServerChannel chan struct{},
) *vm.LuaFunction {

	return luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		config, err := l.CheckTable(1)
		if err != nil {
			l.Error("unable to parse 'lexgo.new(config)'")
			return nil
		}

		parseEnv(l, config, data)
		parseAllowedOrigins(l, config, data)

		routerDriver.ClearRoutes()

		app := RegisterApp(l, routerDriver)
		registerListenFunction(l, app, data, restartServerChannel)

		return app
	})
}

func registerListenFunction(
	l vm.LVm,
	app *vm.LuaTable,
	data *AppData,
	restartServerChannel chan struct{},
) {

	listenFn := l.NewFunction(func(l vm.LVm) vm.LuaValue {
		port, err := l.CheckNumber(1)
		if err != nil {
			l.Error("unable to parse port in 'app.listen(port)'")
			return nil
		}

		app.SetField("_port", vm.LuaNumber(port))

		if data.Port != 0 {
			restartServerChannel <- struct{}{}
		}

		data.Port = int32(port)

		return nil
	})

	app.SetField("listen", listenFn)
}
