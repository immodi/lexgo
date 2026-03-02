package watcher

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Watcher struct {
	SrcPath     string
	SrcFilesMap map[string][]byte
	Callback    func() error
	Cancel      context.CancelFunc
}

func (w *Watcher) Watch(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			files, err := os.ReadDir(w.SrcPath)
			if err != nil {
				return err
			}

			for _, entry := range files {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".lua") {
					continue
				}

				path := filepath.Join(w.SrcPath, entry.Name())
				content, err := w.readLuaFile(path)
				if err != nil {
					log.Println("error reading file:", err)
					continue
				}

				if w.compareFiles(entry.Name(), content) {
					log.Printf("reloading file %v...", entry.Name())
					if err := w.Callback(); err != nil {
						log.Println("callback error:", err)
						continue
					}
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (w *Watcher) readLuaFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func (w *Watcher) compareFiles(filePath string, newContent []byte) bool {
	oldContent, ok := w.SrcFilesMap[filePath]

	if ok && bytes.Equal(oldContent, newContent) {
		return false
	}

	w.SrcFilesMap[filePath] = newContent
	return true

}
