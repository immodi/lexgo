package main

import (
	"immodi/lexgo/cmd/build"
	runtime_cli "immodi/lexgo/cmd/runtime"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		log.Fatal("didn't supply command")
	}

	switch args[0] {
	case "run":
		if err := runtime_cli.Run(args[1:]); err != nil {
			log.Fatalf(err.Error())
		}

	case "build":
		bArgs := args[1:]
		if len(bArgs) <= 0 {
			log.Fatalf("didn't supply the lua source directory")
		}

		builder := build.ApplicationBuilder{SrcDir: bArgs[0]}
		if err := builder.Build(); err != nil {
			log.Fatalf(err.Error())
		}
	default:
		log.Fatalf("supplied command is unknown")
	}
}
