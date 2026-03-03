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
		if err := runtime_cli.RunAndWatch(args[1:]); err != nil {
			log.Fatalf(err.Error())
		}
	case "build":
		if err := build.Build(args[1:]); err != nil {
			log.Fatalf(err.Error())
		}
	default:
		log.Fatalf("supplied command is unknown")
	}
}
