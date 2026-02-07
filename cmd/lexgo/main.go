package main

import (
	"fmt"
	"immodi/lexgo/internal/app"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		log.Fatal("didn't supply main lua file")
	}

	mainFile := args[0]
	app, err := app.New(mainFile)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	err = app.Listen(fmt.Sprintf(":%d", app.Port))
	if err != nil {
		log.Fatal(err)
	}
}
