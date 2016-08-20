package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/satori/go.uuid"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	var cPath string
	flag.StringVar(&cPath, "p", "", "check path")
	flag.Parse()

	if len(cPath) == 0 {
		fmt.Println("Please input path. Use -p flag")
		return
	}

	src, err := os.Stat(cPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !src.IsDir() {
		fmt.Println("check path is not directory")
		return
	}

	root, err := filepath.Abs(cPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("checking %s\n", root)

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	paths, errc := walkFiles(ctx, root)
	c := make(chan string)

	workers := runtime.NumCPU()

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			contentCheck(ctx, paths, c)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for p := range c {
		if len(p) > 0 {
			newPath := changeName(p)
			if len(newPath) > 0 {
				fmt.Printf("%s -> %s\n", p, newPath)
			}
		}
	}

	if err := <-errc; err != nil {
		fmt.Println(err)
	}
}

func walkFiles(ctx context.Context, root string) (<-chan string, <-chan error) {
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

func contentCheck(ctx context.Context, paths <-chan string, c chan<- string) {
	for path := range paths {
		select {
		case c <- detectContectType(path):
		case <-ctx.Done():
			return
		}
	}
}

func detectContectType(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	filetype := http.DetectContentType(buf)
	switch filetype {
	case "image/jpeg", "image/jpg":
		return path
	case "image/png":
		return path
	default:
	}
	return ""
}

func changeName(path string) string {
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
