package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/image-transform-test/go/graphicsmagick"
	"github.com/image-transform-test/go/vips"
)

const (
	width   = 700
	height  = 700
	quality = 100
	//filename = "file1.jpeg"
	filename = "file2.jpg"
)

var (
	rootDir   = "images/"
	MAX_BATCH = 10
)

func main2() {
	/* http.HandleFunc("/upload", UploadHandler)

	log.Fatal(grace.Serve(":9000", nil)) */
	vips.Resize(filename, rootDir, "results/vips/", width, height, quality)
	graphicsmagick.Resize(filename, rootDir, "results/gmagick/", width, height, quality)
}

func main() {
	var wg sync.WaitGroup

	chLimiter := make(chan struct{}, MAX_BATCH)

	for i := 0; i < MAX_BATCH; i++ {
		chLimiter <- struct{}{}
	}

	files := make([]string, 0, 100)

	if err := filepath.Walk("images/100-images", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		files = append(files, path)

		return nil
	}); err != nil {
		panic(err)
	}

	// Delete root files
	files = files[1:]

	start := time.Now()
	// Do Resize
	for _, url := range files {
		<-chLimiter
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			vips.Resize(url[18:], url[:18], "results/vips/", width, height, quality)
			//graphicsmagick.Resize(url[18:], url[:18], "results/gmagick/", width, height, quality)

			chLimiter <- struct{}{}
		}(url)
	}

	wg.Wait()
	fmt.Println(time.Since(start))
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
