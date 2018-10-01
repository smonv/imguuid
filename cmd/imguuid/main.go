package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"imguuid"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	cPath := os.Args[1]

	if len(cPath) == 0 {
		fmt.Println("Please input check path.")
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

	paths, errc := imguuid.WalkFiles(ctx, root)
	c := make(chan string)

	workers := runtime.NumCPU()

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			imguuid.ContentCheck(ctx, paths, c)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for p := range c {
		if len(p) > 0 {
			newPath := imguuid.ChangeName(p)
			if len(newPath) > 0 {
				fmt.Printf("%s -> %s\n", p, newPath)
			}
		}
	}

	if err := <-errc; err != nil {
		fmt.Println(err)
	}
}
