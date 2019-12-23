package main

import (
	"net/http"

	"github.com/image-transform-test/go/graphicsmagick"
	"github.com/image-transform-test/go/vips"
)

const (
	width    = 200
	height   = 200
	quality  = 70
	filename = "file2.jpg"
)

func main() {
	/* http.HandleFunc("/upload", UploadHandler)

	log.Fatal(grace.Serve(":9000", nil)) */
	vips.Resize(filename, width, height, quality)
	graphicsmagick.Resize(filename, width, height, quality)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
