package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	defer os.Exit(1)
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

	fullPath, err := filepath.Abs(cPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("checking %s", fullPath)
	if err = filepath.Walk(fullPath, walker); err != nil {
		fmt.Println(err)
	}
}

func walker(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		buf := make([]byte, 512)
		_, err = file.Read(buf)
		if err != nil {
			return err
		}

		filetype := http.DetectContentType(buf)
		switch filetype {
		case "image/jpeg", "image/jpg":
		case "image/png":
		default:
		}
	}
	return nil
}
