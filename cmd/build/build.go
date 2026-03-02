package build

import (
	"fmt"
	"immodi/lexgo/internal/builder"
)

func Build(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("didn't supply the lua source directory")
	}

	builder := builder.ApplicationBuilder{SrcDir: args[0]}
	if err := builder.Build(); err != nil {
		return err
	}

	return nil
}
