package framework

import "immodi/lexgo/internal/vm"

func parseEnv(l vm.LVm, config *vm.LuaTable, data *AppData) {
	env := config.GetField("env")

	switch env.Type() {

	case vm.LTString:
		parsed, err := ParseAppEnv(env.String())
		if err != nil {
			l.Error(err.Error())
			return
		}
		data.Env = parsed

	case vm.LTNil:
		l.Warn("'config.env' not set, defaulting to 'dev'")
		data.Env = ENVDev

	default:
		l.Error("unable to parse 'config.env', required type is 'string | nil'")
	}
}

func parseAllowedOrigins(l vm.LVm, config *vm.LuaTable, data *AppData) {

	allowedOrigins, ok := vm.GenericGetField[*vm.LuaTable](config, "allowedOrigins")

	if !ok {
		switch data.Env {
		case ENVDev:
			data.AllowedOrigins = []string{"*"}
			l.Warn("'config.allowedOrigins' not set correctly, defaulting to '*'")

		case ENVProd:
			l.Error("'config.allowedOrigins' not set correctly in production")
		}
		return
	}

	allowedOrigins.ForEach(func(_, value vm.LuaValue) {
		origin, ok := value.(vm.LuaString)
		if !ok {
			l.Error("'config.allowedOrigins' should be []string")
			return
		}

		data.AllowedOrigins = append(data.AllowedOrigins, origin.String())
	})
}
