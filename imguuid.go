package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	args := os.Args

	if len(args) == 0 {
		fmt.Println("Please input path")
		return
	}

	cPath := args[1]

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

	done := make(chan struct{})
	paths, errc := walkFiles(done, root)
	c := make(chan string)

	// workers := runtime.NumCPU()
	workers := 20

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			contentCheck(done, paths, c)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for p := range c {
		if len(p) > 0 {
			fmt.Println(p)
		}
	}

	if err := <-errc; err != nil {
		fmt.Println(err)
	}

	defer close(done)
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
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
			case <-done:
				return errors.New("walker canceled")
			}
			return nil
		})
	}()

	return paths, errc
}

func contentCheck(done <-chan struct{}, paths <-chan string, c chan<- string) {
	for path := range paths {
		select {
		case c <- detectContectType(path):
		case <-done:
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
