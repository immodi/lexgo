package vm

import (
	"log"
)

type MiddlewaresDriver interface {
	ExecuteFinal(
		fn *LuaFunction,
	)
	HandleError(msg string)
	LuaVm() LVm
	GetLuaResponse() Response
	GetLuaRequest() Request
}

type MiddlewaresContext struct {
	MiddlewaresDriver MiddlewaresDriver
	FinalHandler      *LuaFunction
	index             int
}

type Request interface {
	MakeLuaRequest() *LuaTable
}

type Response interface {
	MakeLuaResponse() *LuaTable
}

func ExecuteMiddlewares(ctx *MiddlewaresContext, stack []*LuaFunction) {
	runNext(ctx, stack)
}

func RegisterDefaultMiddlewares(luaVm LVm, tbl *LuaTable) {
	mwTbl := luaVm.NewTable()
	tbl.SetField("middlewares", mwTbl)

	// L.SetField(mwTbl, "logger", middlewares.DefaultLuaLogger(L))
	// L.SetField(mwTbl, "cors", middlewares.DefaultLuaCORS(L))
}

func runNext(ctx *MiddlewaresContext, stack []*LuaFunction) {
	if ctx.index >= len(stack) {
		ctx.MiddlewaresDriver.ExecuteFinal(
			ctx.FinalHandler,
		)
		return
	}

	current := stack[ctx.index]
	ctx.index++

	luaVm := ctx.MiddlewaresDriver.LuaVm()

	next := luaVm.NewFunction(func(l LVm) LuaValue {
		runNext(ctx, stack)
		return nil
	})

	err := luaVm.RunFunction(
		current,
		ctx.MiddlewaresDriver.GetLuaRequest().MakeLuaRequest(),
		ctx.MiddlewaresDriver.GetLuaResponse().MakeLuaResponse(),
		next,
	)

	if err != nil {
		log.Printf("Lua middleware error: %s", err)
		ctx.MiddlewaresDriver.HandleError(err.Error())
	}
}

func NewMiddlewaresContext(
	driver MiddlewaresDriver,
	final *LuaFunction,
) *MiddlewaresContext {
	return &MiddlewaresContext{
		MiddlewaresDriver: driver,
		FinalHandler:      final,
	}
}
