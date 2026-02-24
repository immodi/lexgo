package runtime_cli

import (
	"fmt"
	"immodi/lexgo/internal/runtime"
)

func Run(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("didn't supply main lua file")
	}

	mainFile := args[0]
	rt, err := runtime.New()
	if err != nil {
		return err
	}
	defer rt.Close()

	err = rt.LoadFile(mainFile)
	if err != nil {
		return err
	}

	return rt.Start()
}
