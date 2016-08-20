package imguuid

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// WalkFiles ...
func WalkFiles(ctx context.Context, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	}()

	return paths, errc
}

// ChangeName ...
func ChangeName(path string) string {
	basename := filepath.Base(path)
	fileExt := filepath.Ext(path)
	fileDir := filepath.Dir(path)

	filename := strings.TrimSuffix(basename, fileExt)
	_, err := uuid.FromString(filename)
	if err == nil {
		return ""
	}
	u := uuid.NewV4()
	newFilename := u.String() + fileExt
	newPath := filepath.Join(fileDir, newFilename)
	err = os.Rename(path, newPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return newPath
}
