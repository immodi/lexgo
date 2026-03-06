package runtime_cli

import (
	"context"
	"fmt"
	"immodi/lexgo/internal/runtime"
	"immodi/lexgo/internal/watcher"
	"log"
	"path/filepath"
	"time"
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
		return rt.DoFile(filePath)
	}

	if err = loadFile(); err != nil {
		return err
	}

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	serverErrCh := make(chan error, 1)
	watcherErrCh := make(chan error, 1)

	go func() {
		for {
			serverCtx, serverCancel := context.WithCancel(context.Background())

			go func(ctx context.Context) {
				if err := rt.Start(ctx); err != nil {
					serverErrCh <- err
				}
			}(serverCtx)

			select {
			case <-rt.RestartCh:
				log.Println("restarting server...")
				serverCancel()
				// wait a short moment to avoid race on port
				time.Sleep(100 * time.Millisecond)
				continue
			case err := <-serverErrCh:
				log.Println("runtime error:", err)
				serverCancel()
				rootCancel() // stop everything
				return
			case <-rootCtx.Done():
				serverCancel()
				return
			}
		}
	}()

	go func() {
		srcPath := filepath.Dir(filePath)
		watcher := watcher.Watcher{
			SrcPath:     srcPath,
			SrcFilesMap: make(map[string][]byte),
			Callback:    loadFile,
		}

		if err := watcher.Watch(rootCtx); err != nil {
			log.Println("watcher error:", err)
			watcherErrCh <- err
			rootCancel()
		}
	}()

	select {
	case err := <-watcherErrCh:
		return err
	case <-rootCtx.Done():
		return rootCtx.Err()
	}
}
