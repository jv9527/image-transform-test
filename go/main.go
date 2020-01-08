package main

import (
	"net/http"

	"github.com/jv9527/image-transform-test/go/vips"
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
	//vips.Resize("file2.jpg", "images/", "results/vips/", width, height, quality)
	//graphicsmagick.Resize("file2.jpg", "images/", "results/gmagick/", width, height, quality)

	for x := 0; x < 20; x++ {
		go func() {
			vips.Convert("test1.png", rootDir, "results/vips/")

			vips.Convert("test2.png", rootDir, "results/vips/")

			vips.Convert("test3.png", rootDir, "results/vips/")

			vips.Convert("test4.png", rootDir, "results/vips/")

			vips.Convert("test5.png", rootDir, "results/vips/")

		}()

	}
	//vips.Convert(filename, rootDir, "results/vips/")
	//graphicsmagick.Convert(filename, rootDir, "results/gmagick/")
	//imgsosick.Convert(filename, rootDir, "results/imagick/")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
