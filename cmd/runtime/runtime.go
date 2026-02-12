package runtime_cli

import (
	"fmt"
	"immodi/lexgo/internal/runtime"
)

func Runtime(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("didn't supply main lua file")
	}

	mainFile := args[0]
	rt, err := runtime.New(mainFile)
	if err != nil {
		return err
	}
	defer rt.Close()

	return rt.Listen(fmt.Sprintf(":%d", rt.Port))
}
