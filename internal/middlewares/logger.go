package middlewares

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

// DefaultLuaLogger returns a Lua function that logs request info
func DefaultLuaLogger(L *lua.LState) *lua.LFunction {
	return L.NewFunction(func(L *lua.LState) int {
		// Arguments: req, res, next
		req := L.CheckTable(1)
		next := L.CheckFunction(3)

		// Get request method and url
		method := L.GetField(req, "method").String()
		url := L.GetField(req, "url").String()

		// Print basic info
		fmt.Printf("[LOG] %s %s\n", method, url)

		// // Optionally print query parameters
		// query := L.GetField(req, "query")
		// if tbl, ok := query.(*lua.LTable); ok {
		// 	fmt.Println("[LOG] Query params:")
		// 	tbl.ForEach(func(k, v lua.LValue) {
		// 		key := k.String()
		// 		switch val := v.(type) {
		// 		case *lua.LTable:
		// 			arr := []string{}
		// 			val.ForEach(func(_, item lua.LValue) {
		// 				arr = append(arr, item.String())
		// 			})
		// 			fmt.Printf("  %s: %v\n", key, arr)
		// 		case lua.LString:
		// 			fmt.Printf("  %s: %s\n", key, string(val))
		// 		default:
		// 			fmt.Printf("  %s: %v\n", key, val)
		// 		}
		// 	})
		// }

		// Call next middleware
		L.CallByParam(lua.P{
			Fn:      next,
			NRet:    0,
			Protect: true,
		})

		return 0
	})
}
