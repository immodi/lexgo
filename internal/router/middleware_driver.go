package router

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
)

type MiddlewareVmDriver struct {
	Router      *Router
	LuaRequest  *framework.LuaRequest
	LuaResponse *framework.LuaResponse
}

func (mwDriver *MiddlewareVmDriver) LuaVm() vm.LVm {
	return mwDriver.Router.LuaVm
}

func (mwDriver *MiddlewareVmDriver) GetLuaRequest() framework.Request {
	return mwDriver.LuaRequest
}

func (mwDriver *MiddlewareVmDriver) GetLuaResponse() framework.Response {
	return mwDriver.LuaResponse
}

func (mwDriver *MiddlewareVmDriver) ExecuteFinal(
	fn *vm.LuaFunction,
) {
	framework.ExecuteLuaHandler(
		mwDriver.Router.LuaVm,
		mwDriver.Router.ServerErrorFunc,
		fn,
		mwDriver.LuaRequest,
		mwDriver.LuaResponse,
	)
}

func (mwDriver *MiddlewareVmDriver) HandleError(msg string) {
	framework.HandleServerError(
		mwDriver.Router.LuaVm,
		mwDriver.Router.ServerErrorFunc,
		msg,
		mwDriver.LuaResponse,
	)
}
