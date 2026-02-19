package framework

import (
	"immodi/lexgo/internal/middlewares"
	"immodi/lexgo/internal/vm"
	"log"
)

type MiddlewaresContext struct {
	MiddlewaresDriver MiddlewaresDriver
	FinalHandler      *vm.LuaFunction
	RouteSpecificMws  []*vm.LuaFunction
	index             int
}

type AppProviderMiddlewaresImplementation struct {
	appProvider        MiddlewaresAppProvider
	getRegisterdRoutes func() map[string][]string
}

type MiddlewaresDriver interface {
	ExecuteFinal(
		fn *vm.LuaFunction,
	)
	HandleError(msg string)
	LuaVm() vm.LVm
	GetLuaResponse() Response
	GetLuaRequest() Request
}

type MiddlewaresAppProvider interface {
	GetAllowedOrigin(requestOrigin string) string
	GetRegisterdRoutes(getRoutes func() map[string][]string) map[string][]string
	GetAllowedMethods(url string, getRoutes func() map[string][]string) string
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

type Request interface {
	MakeLuaRequest() *vm.LuaTable
}

type Response interface {
	MakeLuaResponse() *vm.LuaTable
}

func ExecuteMiddlewares(ctx *MiddlewaresContext, stack []*vm.LuaFunction) {
	if ctx.RouteSpecificMws != nil {
		stack = append(stack, ctx.RouteSpecificMws...)
	}

	runNext(ctx, stack)
}

func RegisterDefaultMiddlewares(luaVm vm.LVm, tbl *vm.LuaTable, getRoutes func() map[string][]string, appProvider MiddlewaresAppProvider) {
	mwTbl := luaVm.NewTable()
	middlewaresAppProvider := &AppProviderMiddlewaresImplementation{appProvider: appProvider, getRegisterdRoutes: getRoutes}
	tbl.SetField("middlewares", mwTbl)

	mwTbl.SetField("logger", middlewares.DefaultLuaLogger(luaVm))
	mwTbl.SetField("cors", middlewares.DefaultLuaCORS(luaVm, middlewaresAppProvider))
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
	routeSpecificMws []*vm.LuaFunction,
) *MiddlewaresContext {
	return &MiddlewaresContext{
		MiddlewaresDriver: driver,
		FinalHandler:      final,
		RouteSpecificMws:  routeSpecificMws,
	}
}
