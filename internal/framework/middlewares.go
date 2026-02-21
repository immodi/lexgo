package framework

import (
	"immodi/lexgo/internal/middlewares"
	"immodi/lexgo/internal/vm"
	"log"
)

type ExecutionDriver interface {
	ExecuteFinal(
		handler RouterServerHandler,
		serverErr RouterServerHandler,
	)
	HandleError(msg string, serverErr RouterServerHandler)
	LuaVm() vm.LVm
	GetLuaResponse() Response
	GetLuaRequest() Request
}

type MiddlewaresAppProvider interface {
	GetAllowedOrigin(requestOrigin string) string
	GetRegisterdRoutes(getRoutes func() map[string][]string) map[string][]string
	GetAllowedMethods(url string, getRoutes func() map[string][]string) string
}

type Request interface {
	MakeLuaRequest() *vm.LuaTable
}

type Response interface {
	MakeLuaResponse() *vm.LuaTable
}

type ExecutionContext struct {
	ExecutionDriver  ExecutionDriver
	FinalHandler     RouterServerHandler
	NotFoundHandler  RouterServerHandler
	ServerErrHandler RouterServerHandler
	MiddleWares      []*vm.LuaFunction
	index            int
}

type AppProviderMiddlewaresImplementation struct {
	appProvider        MiddlewaresAppProvider
	getRegisterdRoutes func() map[string][]string
}

func (p *AppProviderMiddlewaresImplementation) GetRegisterdRoutes() map[string][]string {
	return p.appProvider.GetRegisterdRoutes(p.getRegisterdRoutes)
}

func (p *AppProviderMiddlewaresImplementation) GetAllowedMethods(url string) string {
	return p.appProvider.GetAllowedMethods(url, p.getRegisterdRoutes)
}

func (p *AppProviderMiddlewaresImplementation) GetAllowedOrigin(requestOrigin string) string {
	return p.appProvider.GetAllowedOrigin(requestOrigin)
}

func Execute(ctx *ExecutionContext) {
	stack := ctx.MiddleWares
	runNext(ctx, stack)
}

func RegisterDefaultMiddlewares(luaVm vm.LVm, tbl *vm.LuaTable, getRoutes func() map[string][]string, appProvider MiddlewaresAppProvider) {
	mwTbl := luaVm.NewTable()
	middlewaresAppProvider := &AppProviderMiddlewaresImplementation{appProvider: appProvider, getRegisterdRoutes: getRoutes}
	tbl.SetField("middlewares", mwTbl)

	mwTbl.SetField("logger", middlewares.DefaultLuaLogger(luaVm))
	mwTbl.SetField("cors", middlewares.DefaultLuaCORS(luaVm, middlewaresAppProvider))
}

func runNext(ctx *ExecutionContext, stack []*vm.LuaFunction) {
	if ctx.index >= len(stack) {
		ctx.ExecutionDriver.ExecuteFinal(
			ctx.FinalHandler,
			ctx.ServerErrHandler,
		)
		return
	}

	current := stack[ctx.index]
	ctx.index++

	luaVm := ctx.ExecutionDriver.LuaVm()
	next := luaVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		runNext(ctx, stack)
		return nil
	})

	err := luaVm.RunFunction(
		current,
		ctx.ExecutionDriver.GetLuaRequest().MakeLuaRequest(),
		ctx.ExecutionDriver.GetLuaResponse().MakeLuaResponse(),
		next,
	)

	if err != nil {
		log.Printf("Lua middleware error: %s", err)
		ctx.ExecutionDriver.HandleError(err.Error(), ctx.ServerErrHandler)
	}
}

func NewExecutionContext(
	driver ExecutionDriver,
	final RouterServerHandler,
	notFound RouterServerHandler,
	serverErr RouterServerHandler,
	middlewares []*vm.LuaFunction,
) *ExecutionContext {
	return &ExecutionContext{
		ExecutionDriver:  driver,
		FinalHandler:     final,
		NotFoundHandler:  notFound,
		ServerErrHandler: serverErr,
		MiddleWares:      middlewares,
	}
}
