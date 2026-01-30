package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	Router := &Router{
		L:      L,
		Routes: make(map[string]*lua.LFunction),
	}
	defer L.Close()
	Router.registerRoutes()

	if err := L.DoFile("lua/tests/hello.lua"); err != nil {
		panic(err)
	}

	log.Println("http server starting at http://localhost:8080....")
	err := http.ListenAndServe(":8080", Router)

	if err != nil {
		log.Fatal("failed to start server")
	}
}

type Router struct {
	L      *lua.LState
	Routes map[string]*lua.LFunction
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fn, ok := router.Routes[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
		return
	}
	router.L.Push(fn)

	// req table
	luaReq := router.L.NewTable()
	router.L.SetField(luaReq, "method", lua.LString(req.Method))
	router.L.SetField(luaReq, "url", lua.LString(req.URL.Path))

	// res table
	luaRes := router.L.NewTable()
	router.L.SetField(luaRes, "send", router.L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, msg)
		return 0
	}))

	router.L.Push(luaReq)
	router.L.Push(luaRes)

	if err := router.L.PCall(2, 0, nil); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (router *Router) registerRoutes() {
	// app table
	tbl := router.L.NewTable()
	router.L.SetGlobal("app", tbl)

	// add :get method
	router.L.SetField(tbl, "get", router.L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		fn := L.CheckFunction(2)

		router.Routes[path] = fn
		return 0
	}))
}
