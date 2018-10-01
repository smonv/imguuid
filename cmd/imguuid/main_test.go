package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"imguuid"
)

func TestDetectContentType(t *testing.T) {
	path := createPNG("testdct.png")
	dPath := imguuid.DetectContentType(path)
	if len(dPath) == 0 {
		t.Errorf("Cannot detect png file")
	}
	os.Remove(path)

	path = createJPG("testdct.jpg")
	dPath = imguuid.DetectContentType(path)
	if len(dPath) == 0 {
		t.Errorf("Cannot detect jpg file")
	}
	os.Remove(path)

	path = createText("testdct.txt")
	dPath = imguuid.DetectContentType(path)
	if len(dPath) > 0 {
		t.Errorf("Wrong detect, txt is not image file")
	}
	os.Remove(path)
}

func TestChangeName(t *testing.T) {
	path := createPNG("testCN.png")
	newPath := imguuid.ChangeName(path)

	_, err := os.Stat(newPath)
	if err != nil {
		t.Errorf(err.Error())
	}
	os.Remove(newPath)
}

func createPNG(filename string) string {
	path := filepath.Join("/tmp", filename)

	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))

	// Draw a red dot at (2, 3)
	img.Set(2, 3, color.RGBA{255, 0, 0, 255})

	// Save to out.png
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	png.Encode(f, img)
	return path
}

func createJPG(filename string) string {
	path := filepath.Join("/tmp", filename)

	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))

	// Draw a red dot at (2, 3)
	img.Set(2, 3, color.RGBA{255, 0, 0, 255})

	// Save to out.png
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, &jpeg.Options{Quality: 80})
	return path
}

func createText(filename string) string {
	path := filepath.Join("/tmp", filename)
	content := []byte("TestDCT")
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		fmt.Println(err)
	}
	return path
}
