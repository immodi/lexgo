package engine

import (
	"immodi/lexgo/internal/framework"
	"immodi/lexgo/internal/vm"
)

type ExecutionDriver struct {
	LVm         vm.LVm
	LuaRequest  *framework.LuaRequest
	LuaResponse *framework.LuaResponse
}

func (mwDriver *ExecutionDriver) LuaVm() vm.LVm {
	return mwDriver.LVm
}

func (mwDriver *ExecutionDriver) GetLuaRequest() framework.Request {
	return mwDriver.LuaRequest
}

func (mwDriver *ExecutionDriver) GetLuaResponse() framework.Response {
	return mwDriver.LuaResponse
}

func (mwDriver *ExecutionDriver) ExecuteFinal(
	fn framework.RouterServerHandler,
	serverErr framework.RouterServerHandler,
) {
	framework.ExecuteLuaHandler(
		mwDriver.LVm,
		serverErr,
		fn,
		mwDriver.LuaRequest,
		mwDriver.LuaResponse,
	)
}

func (mwDriver *ExecutionDriver) HandleError(msg string, serverErr framework.RouterServerHandler) {
	framework.HandleServerError(
		mwDriver.LVm,
		serverErr,
		msg,
		mwDriver.LuaRequest,
		mwDriver.LuaResponse,
	)
}
