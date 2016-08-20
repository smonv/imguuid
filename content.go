package imguuid

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// ContentCheck ...
func ContentCheck(ctx context.Context, paths <-chan string, c chan<- string) {
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
