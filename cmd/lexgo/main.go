package main

import (
	"immodi/lexgo/internal/app"
	"log"
)

func main() {
	app, err := app.New("examples/hello/hello.lua")
	if err != nil {
		log.Fatal(err)
	}

	defer app.Close()

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
