package router

import (
	lua "github.com/yuin/gopher-lua"
)

type Handler struct {
	Pattern string
	Handler *lua.LFunction
	Params  map[string]string
}
