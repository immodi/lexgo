package builder

import "fmt"

const mainSrc string = `
package main

import (
	"context"
	"embed"
	"immodi/lexgo/internal/runtime"
	"log"
)

//go:embed %v
var src embed.FS

func main() {
	mainFile, err := src.ReadFile("%v/main.lua")
	if err != nil {
		log.Fatal(err.Error())
	}

	rt, err := runtime.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rt.Close()

	err = rt.LoadString(string(mainFile))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = rt.Start(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
}
`

func (b *ApplicationBuilder) generateMain() string {
	return fmt.Sprintf(mainSrc, b.SrcDir, b.SrcDir)
}
