package vm

import (
	"immodi/lexgo/internal/middlewares"
	"log"

	lua "github.com/yuin/gopher-lua"
)

type MiddlewaresDriver interface {
	ExecuteFinal(
		fn *lua.LFunction,
	)
	HandleError(msg string)
	LuaState() *lua.LState
	GetLuaResponse() Response
	GetLuaRequest() Request
}

type MiddlewaresContext struct {
	MiddlewaresDriver MiddlewaresDriver
	FinalHandler      *lua.LFunction
	index             int
}

type Request interface {
	MakeLuaRequest() *lua.LTable
}

type Response interface {
	MakeLuaResponse() *lua.LTable
}

func ExecuteMiddlewares(ctx *MiddlewaresContext, stack []*lua.LFunction) {
	runNext(ctx, stack)
}

func RegisterDefaultMiddlewares(L *lua.LState, tbl *lua.LTable) {
	mwTbl := L.NewTable()
	L.SetField(tbl, "middlewares", mwTbl)

	L.SetField(mwTbl, "logger", middlewares.DefaultLuaLogger(L))
	L.SetField(mwTbl, "cors", middlewares.DefaultLuaCORS(L))
}

func runNext(ctx *MiddlewaresContext, stack []*lua.LFunction) {
	if ctx.index >= len(stack) {
		ctx.MiddlewaresDriver.ExecuteFinal(
			ctx.FinalHandler,
		)
		return
	}

	current := stack[ctx.index]
	ctx.index++

	L := ctx.MiddlewaresDriver.LuaState()

	next := L.NewFunction(func(L *lua.LState) int {
		runNext(ctx, stack)
		return 0
	})

	L.Push(current)
	L.Push(ctx.MiddlewaresDriver.GetLuaRequest().MakeLuaRequest())
	L.Push(ctx.MiddlewaresDriver.GetLuaResponse().MakeLuaResponse())
	L.Push(next)

	if err := L.PCall(3, 0, nil); err != nil {
		log.Printf("Lua middleware error: %s", err)
		ctx.MiddlewaresDriver.HandleError(err.Error())
	}
}

func NewMiddlewaresContext(
	driver MiddlewaresDriver,
	final *lua.LFunction,
) *MiddlewaresContext {
	return &MiddlewaresContext{
		MiddlewaresDriver: driver,
		FinalHandler:      final,
	}
}
