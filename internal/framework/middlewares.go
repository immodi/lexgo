package framework

import (
	defaultmiddlewares "immodi/lexgo/internal/framework/def_libs/middlewares"
	"immodi/lexgo/internal/vm"
	"log"
)

type MiddlewaresDriver interface {
	ExecuteFinal(
		fn *vm.LuaFunction,
	)
	HandleError(msg string)
	LuaVm() vm.LVm
	GetLuaResponse() Response
	GetLuaRequest() Request
}

type MiddlewaresContext struct {
	MiddlewaresDriver MiddlewaresDriver
	FinalHandler      *vm.LuaFunction
	index             int
}

type Request interface {
	MakeLuaRequest() *vm.LuaTable
}

type Response interface {
	MakeLuaResponse() *vm.LuaTable
}

func ExecuteMiddlewares(ctx *MiddlewaresContext, stack []*vm.LuaFunction) {
	runNext(ctx, stack)
}

func RegisterDefaultMiddlewares(luaVm vm.LVm, tbl *vm.LuaTable, getRoutes func() map[string][]string) {
	mwTbl := luaVm.NewTable()
	tbl.SetField("middlewares", mwTbl)

	mwTbl.SetField("logger", defaultmiddlewares.DefaultLuaLogger(luaVm))
	mwTbl.SetField("cors", defaultmiddlewares.DefaultLuaCORS(luaVm, getRoutes))
}

func runNext(ctx *MiddlewaresContext, stack []*vm.LuaFunction) {
	if ctx.index >= len(stack) {
		ctx.MiddlewaresDriver.ExecuteFinal(
			ctx.FinalHandler,
		)
		return
	}

	current := stack[ctx.index]
	ctx.index++

	luaVm := ctx.MiddlewaresDriver.LuaVm()
	next := luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
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
	final *vm.LuaFunction,
) *MiddlewaresContext {
	return &MiddlewaresContext{
		MiddlewaresDriver: driver,
		FinalHandler:      final,
	}
}
