package defaultmiddlewares

import (
	"fmt"
	"immodi/lexgo/internal/vm"
)

func DefaultLuaLogger(LVm vm.LVm) *vm.LuaFunction {
	return LVm.NewFunction(func(l vm.LVm) vm.LuaValue {
		req, err := l.CheckTable(1)
		if err != nil {
			return nil
		}

		next, err := l.CheckFunction(3)
		if err != nil {
			return nil
		}

		method := req.GetField("method").String()
		url := req.GetField("url").String()

		fmt.Printf("[LOG] %s %s\n", method, url)

		// // Optionally print query parameters
		// query := req.GetField("query")
		// if tbl, ok := query.(*vm.LuaTable); ok {
		// 	fmt.Println("[LOG] Query params:")
		// 	tbl.ForEach(func(k, v vm.LuaValue) {
		// 		key := k.String()
		// 		switch val := v.(type) {
		// 		case *vm.LuaTable:
		// 			arr := []string{}
		// 			val.ForEach(func(_, item vm.LuaValue) {
		// 				arr = append(arr, item.String())
		// 			})
		// 			fmt.Printf("  %s: %v\n", key, arr)
		// 		case vm.LuaString:
		// 			fmt.Printf("  %s: %s\n", key, string(val))
		// 		default:
		// 			fmt.Printf("  %s: %v\n", key, val)
		// 		}
		// 	})
		// }

		l.RunFunction(next)
		return nil
	})
}
