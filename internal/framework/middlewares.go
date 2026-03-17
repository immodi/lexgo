package framework

import (
	"immodi/lexgo/internal/middlewares"
	"immodi/lexgo/internal/vm"
	"log"
	"strings"
)

type ExecutionDriver interface {
	ExecuteFinal(
		handler RouterHandler,
		serverErr RouterHandler,
	)
	HandleError(msg string, serverErr RouterHandler)
	LuaVm() vm.LVm
	GetLuaResponse() *LuaResponse
	GetLuaRequest() *LuaRequest
}

type ExecutionContext struct {
	ExecutionDriver  ExecutionDriver
	FinalHandler     RouterHandler
	NotFoundHandler  RouterHandler
	ServerErrHandler RouterHandler
	MiddleWares      []RouterHandler
	index            int
}

type CORSRuntime struct {
	appData      *AppData
	routerDriver RouterDriver
}

func (cr *CORSRuntime) GetRegisterdRoutes() map[string][]string {
	return cr.routerDriver.GetAllRegistredRoutes()
}

func (cr *CORSRuntime) GetAllowedMethods(url string) string {
	const (
		DEFAULT_METHODS_DEV  = "GET, POST, PUT, DELETE, OPTIONS"
		DEFAULT_METHODS_PROD = ""
	)
	var allowedMethodsString string = DEFAULT_METHODS_DEV

	registerdRotues := cr.GetRegisterdRoutes()
	allowedMethods, ok := registerdRotues[url]

	if cr.appData.IsProduction() {
		allowedMethodsString = DEFAULT_METHODS_PROD
	}

	if ok {
		allowedMethodsString = strings.Join(allowedMethods, ", ")
	}

	return allowedMethodsString
}

func (cr *CORSRuntime) GetAllowedOrigin(requestOrigin string) string {
	return cr.appData.GetAllowedOrigin(requestOrigin)
}

func Execute(ctx *ExecutionContext) {
	stack := ctx.MiddleWares
	runNext(ctx, stack)
}

func RegisterDefaultMiddlewares(luaVm vm.LVm, tbl *vm.LuaTable, cors *CORSRuntime) {
	mwTbl := luaVm.NewTable()

	tbl.SetField("middlewares", mwTbl)
	mwTbl.SetField("logger", middlewares.DefaultLuaLogger(luaVm))
	mwTbl.SetField("cors", middlewares.DefaultLuaCORS(luaVm, cors))
}

func runNext(ctx *ExecutionContext, stack []RouterHandler) {
	if ctx.index >= len(stack) {
		ctx.ExecutionDriver.ExecuteFinal(
			ctx.FinalHandler,
			ctx.ServerErrHandler,
		)
		return
	}

	current := stack[ctx.index]
	ctx.index++

	// luaVm := ctx.ExecutionDriver.LuaVm()
	next := func() {
		runNext(ctx, stack)
	}

	err := current.Handle(
		ctx.ExecutionDriver.GetLuaRequest(),
		ctx.ExecutionDriver.GetLuaResponse(),
		next,
	)

	if err != nil {
		log.Printf("Lua middleware error: %s", err)
		ctx.ExecutionDriver.HandleError(err.Error(), ctx.ServerErrHandler)
	}
}

func NewExecutionContext(
	driver ExecutionDriver,
	final RouterHandler,
	notFound RouterHandler,
	serverErr RouterHandler,
	middlewares []RouterHandler,
) *ExecutionContext {
	return &ExecutionContext{
		ExecutionDriver:  driver,
		FinalHandler:     final,
		NotFoundHandler:  notFound,
		ServerErrHandler: serverErr,
		MiddleWares:      middlewares,
	}
}
