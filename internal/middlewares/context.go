package middlewares

import lua "github.com/yuin/gopher-lua"

type Context struct {
	Driver       Driver
	Req          Request
	Res          Response
	FinalHandler *lua.LFunction
	index        int
}

type Request interface {
	MakeLuaRequest() *lua.LTable
}

type Response interface {
	MakeLuaResponse() *lua.LTable
}

func NewContext(
	driver Driver,
	req Request,
	res Response,
	final *lua.LFunction,
) *Context {
	return &Context{
		Driver:       driver,
		Req:          req,
		Res:          res,
		FinalHandler: final,
	}
}
