package framework

import (
	"immodi/lexgo/internal/vm"
	"strings"
)

type AppData struct {
	Port           int32
	Env            AppEnv
	AllowedOrigins []string
}

func RegisterFramework(luaVm vm.LVm, routerDriver RouterDriver) (*AppData, error) {
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

			data.Port = portInt32
			return nil
		})

		app.SetField("listen", listenFn)
		return app
	})

	tbl.SetField("new", newFn)
	RegisterDefaultMiddlewares(luaVm, tbl, routerDriver.GetAllRegistredRoutes, data)

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

func (d *AppData) GetAllowedMethods(url string, getRoutes func() map[string][]string) string {
	const (
		DEFAULT_METHODS_DEV  = "GET, POST, PUT, DELETE, OPTIONS"
		DEFAULT_METHODS_PROD = ""
	)
	var allowedMethodsString string = DEFAULT_METHODS_DEV

	registerdRotues := getRoutes()
	allowedMethods, ok := registerdRotues[url]

	if d.IsProduction() {
		allowedMethodsString = DEFAULT_METHODS_PROD
	}

	if ok {
		allowedMethodsString = strings.Join(allowedMethods, ", ")
	}

	return allowedMethodsString
}

func (d AppData) GetRegisterdRoutes(getRoutes func() map[string][]string) map[string][]string {
	return getRoutes()
}
