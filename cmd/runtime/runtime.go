package runtime_cli

import (
	"context"
	"fmt"
	"immodi/lexgo/internal/runtime"
	"immodi/lexgo/internal/watcher"
	"log"
	"path/filepath"
)

func RunAndWatch(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("didn't supply main lua file")
	}

	filePath := args[0]

	rt, err := runtime.New()
	if err != nil {
		return err
	}
	defer rt.Close()

	loadFile := func() error {
		return rt.LoadFile(filePath)
	}

	if err = loadFile(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)

	go func() {
		if err := rt.Start(); err != nil {
			log.Println("runtime error:", err)
			cancel()
			errCh <- err
		}
	}()

	go func() {
		srcPath := filepath.Dir(filePath)
		watcher := watcher.Watcher{
			SrcPath:     srcPath,
			SrcFilesMap: make(map[string][]byte),
			Callback:    loadFile,
			Cancel:      cancel,
		}

		if err := watcher.Watch(ctx); err != nil {
			log.Println("watcher error:", err)
			cancel()
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
