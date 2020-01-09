package main

import (
	"net/http"
	"sync"

	"github.com/jv9527/image-transform-test/go/graphicsmagick"
)

const (
	width    = 200
	height   = 200
	quality  = 70
	filename = "test5.png"
)

var (
	rootDir = "../images/png/"
)

func main() {
	//imagick.Initialize()
	//defer imagick.Terminate()
	/* http.HandleFunc("/upload", UploadHandler)
	log.Fatal(grace.Serve(":9000", nil)) */
	// vips.Resize("file2.jpg", "images/", "results/vips/", width, height, quality)
	// graphicsmagick.Resize("file2.jpg", "images/", "results/gmagick/", width, height, quality)

	var wg sync.WaitGroup
	defer wg.Wait()
	for x := 0; x < 100; x++ {
		go func() {
			wg.Add(1)
			graphicsmagick.Convert("test1.png", rootDir, "results/gmagick/")
			wg.Done()
		}()
		// go func() {k
		// 	wg.Add(1)
		// 	vips.Convert("test2.png", rootDir, "results/vips/")
		// 	wg.Done()
		// }()
		// go func() {
		// 	wg.Add(1)
		// 	vips.Convert("test3.png", rootDir, "results/vips/")
		// 	wg.Done()
		// }()
		// go func() {
		// 	wg.Add(1)
		// 	vips.Convert("test4.png", rootDir, "results/vips/")
		// 	wg.Done()
		// }()
		// go func() {
		// 	wg.Add(1)
		// 	vips.Convert("test5.png", rootDir, "results/vips/")
		// 	wg.Done()
		// }()

	}
	// vips.Convert(filename, rootDir, "results/vips/")
	// graphicsmagick.Convert(filename, rootDir, "results/gmagick/")
	// imgsosick.Convert(filename, rootDir, "results/imagick/")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
