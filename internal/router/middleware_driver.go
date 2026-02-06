package router

import (
	"immodi/lexgo/internal/vm"

	lua "github.com/yuin/gopher-lua"
)

type MiddlewareVmDriver struct {
	Router      *Router
	LuaRequest  *vm.LuaRequest
	LuaResponse *vm.LuaResponse
}

func (mwDriver *MiddlewareVmDriver) LuaState() *lua.LState {
	return mwDriver.Router.LuaVm.L
}

func (mwDriver *MiddlewareVmDriver) GetLuaRequest() vm.Request {
	return mwDriver.LuaRequest
}

func (mwDriver *MiddlewareVmDriver) GetLuaResponse() vm.Response {
	return mwDriver.LuaResponse
}

func (mwDriver *MiddlewareVmDriver) ExecuteFinal(
	fn *lua.LFunction,
) {
	vm.ExecuteLuaHandler(
		mwDriver.Router.LuaVm.L,
		mwDriver.Router.ServerErrorFunc,
		fn,
		mwDriver.LuaRequest,
		mwDriver.LuaResponse,
	)
}

func (mwDriver *MiddlewareVmDriver) HandleError(msg string) {
	vm.HandleServerError(
		mwDriver.Router.LuaVm.L,
		mwDriver.Router.ServerErrorFunc,
		msg,
		mwDriver.LuaResponse,
	)
}
