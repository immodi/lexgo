package framework

import (
	"fmt"
	"immodi/lexgo/internal/vm"
)

type AppData struct {
	Port           int32
	Env            AppEnv
	AllowedOrigins []string
}

func RegisterFramework(luaVm vm.LVm, routerDriver RouterDriver, restartServerChannel chan struct{}) (*AppData, error) {
	tbl := luaVm.NewTable()
	data := &AppData{}

	luaVm.SetGlobal("lexgo", tbl)
	newFn := luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		config, err := l.CheckTable(1)
		if err != nil {
			msg := "unable to parse the 'config' argument table in 'lexgo.new(config)'"
			l.Error(msg)
			return nil
		}

		env := config.GetField("env")
		switch env.Type() {
		case vm.LTString:
			env, err := ParseAppEnv(env.String())
			if err != nil {
				l.Error(err.Error())
				return nil
			}
			data.Env = env
		case vm.LTNil:
			l.Warn("'config.env' not set, defaulting to 'dev'")
			data.Env = ENVDev
		default:
			msg := "unable to parse 'config.env', required type is 'string | nil'"
			l.Error(msg)
			return nil
		}

		allowedOrigins, ok := vm.GenericGetField[*vm.LuaTable](config, "allowedOrigins")
		if !ok {
			switch data.Env {
			case ENVDev:
				data.AllowedOrigins = []string{"*"}
				l.Warn("'config.allowedOrigins' of type []string not set correctly, defaulting to '*'")
			case ENVProd:
				l.Error("'config.allowedOrigins' of type []string not set correctly in 'production' mode")
			}
		}

		allowedOrigins.ForEach(func(key, value vm.LuaValue) {
			origin, ok := value.(vm.LuaString)
			if !ok {
				l.Error("unable to parse, 'config.allowedOrigins' should be of type []string")
				return
			}

			data.AllowedOrigins = append(data.AllowedOrigins, origin.String())
		})

		routerDriver.ClearRoutes()
		app := RegisterApp(l, routerDriver)

		listenFn := l.NewFunction(func(l vm.LVm) vm.LuaValue {
			port, err := l.CheckNumber(1)
			if err != nil {
				msg := "unable to parse the port number argument in 'app.listen(port)'"
				l.Error(msg)
				return nil
			}
			app.SetField("_port", vm.LuaNumber(port))
			portInt32 := int32(port)

			if data.Port != 0 {
				restartServerChannel <- struct{}{}
			}

			data.Port = portInt32
			return nil
		})

		app.SetField("listen", listenFn)
		return app
	})

	tbl.SetField("new", newFn)

	lx, err := LxTable(luaVm)
	if err != nil {
		return nil, fmt.Errorf("failed to load 'lx' library into runtime")
	}
	tbl.SetField("lx", lx)

	RegisterDefaultMiddlewares(luaVm, tbl, &CORSRuntime{
		appData:      data,
		routerDriver: routerDriver,
	})

	return data, nil
}

func (d *AppData) IsProduction() bool {
	if d.Env == ENVProd {
		return true
	}

	return false
}

func (d *AppData) GetAllowedOrigin(requestOrigin string) string {
	const (
		DEFAULT_ALLOWED_ORIGIN_DEV  = "*"
		DEFAULT_ALLOWED_ORIGIN_PROD = ""
	)
	var allowedOrigin string = DEFAULT_ALLOWED_ORIGIN_DEV
	if d.IsProduction() {
		allowedOrigin = DEFAULT_ALLOWED_ORIGIN_PROD
	}

	for _, origin := range d.AllowedOrigins {
		if requestOrigin == origin {
			return origin
		}
	}

	return allowedOrigin
}
