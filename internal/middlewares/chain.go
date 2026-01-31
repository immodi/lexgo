package middlewares

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func Execute(ctx *Context, stack []*lua.LFunction) {
	runNext(ctx, stack)
}

func runNext(ctx *Context, stack []*lua.LFunction) {
	if ctx.index >= len(stack) {
		ctx.Driver.ExecuteFinal(
			ctx.FinalHandler,
			ctx.Req,
			ctx.Res,
		)
		return
	}

	current := stack[ctx.index]
	ctx.index++

	L := ctx.Driver.LuaState()

	next := L.NewFunction(func(L *lua.LState) int {
		runNext(ctx, stack)
		return 0
	})

	L.Push(current)
	L.Push(ctx.Req.MakeLuaRequest())
	L.Push(next)

	if err := L.PCall(2, 0, nil); err != nil {
		log.Printf("Lua middleware error: %s", err)
		ctx.Driver.HandleError(err.Error(), ctx.Res)
	}
}
