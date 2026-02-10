package router

import (
	"immodi/lexgo/internal/vm"
)

type MiddlewareVmDriver struct {
	Router      *Router
	LuaRequest  *vm.LuaRequest
	LuaResponse *vm.LuaResponse
}

func (mwDriver *MiddlewareVmDriver) LuaVm() vm.LVm {
	return mwDriver.Router.LuaVm
}

func (mwDriver *MiddlewareVmDriver) GetLuaRequest() vm.Request {
	return mwDriver.LuaRequest
}

func (mwDriver *MiddlewareVmDriver) GetLuaResponse() vm.Response {
	return mwDriver.LuaResponse
}

func (mwDriver *MiddlewareVmDriver) ExecuteFinal(
	fn *vm.LuaFunction,
) {
	vm.ExecuteLuaHandler(
		mwDriver.Router.LuaVm,
		mwDriver.Router.ServerErrorFunc,
		fn,
		mwDriver.LuaRequest,
		mwDriver.LuaResponse,
	)
}

func (mwDriver *MiddlewareVmDriver) HandleError(msg string) {
	vm.HandleServerError(
		mwDriver.Router.LuaVm,
		mwDriver.Router.ServerErrorFunc,
		msg,
		mwDriver.LuaResponse,
	)
}
