package router

import (
	"immodi/lexgo/internal/vm"
)

type Handler struct {
	Pattern string
	Handler *vm.LuaFunction
	Params  map[string]string
}
