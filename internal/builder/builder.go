package builder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	OS_CREATE_FILE = 0644
)

type ApplicationBuilder struct {
	SrcDir string
}

func (b *ApplicationBuilder) checkSrcDir() error {
	info, err := os.Stat(b.SrcDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("source directory '%v' does not exist", b.SrcDir)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("'%v' exists but is not a directory", b.SrcDir)
	}
	return nil
}

func (b *ApplicationBuilder) checkMainLua() error {
	info, err := os.Stat(fmt.Sprintf("%v/main.lua", b.SrcDir))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("main Lua file '%v/main.lua' does not exist", b.SrcDir)
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("'%v/main.lua' exists but is a directory", b.SrcDir)
	}
	return nil
}

func (b *ApplicationBuilder) buildMain() error {
	cmd := exec.Command("go", "build", "-o", "app", "main.go")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Join(err, fmt.Errorf(string(output)))
	}

	err = os.Remove("main.go")
	return err
}

func (b *ApplicationBuilder) Build() error {
	if err := os.WriteFile("main.go", []byte(b.generateMain()), OS_CREATE_FILE); err != nil {
		return err
	}

	if err := b.checkSrcDir(); err != nil {
		return err
	}

	if err := b.checkMainLua(); err != nil {
		return err
	}

	if err := b.buildMain(); err != nil {
		return err
	}

	return nil
}
