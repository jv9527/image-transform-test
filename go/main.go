package main

import (
	"net/http"

	"github.com/image-transform-test/go/graphicsmagick"
	"github.com/image-transform-test/go/vips"
)

const (
	width    = 700
	height   = 700
	quality  = 70
	filename = "file2.jpg"
)

var (
	rootDir = "images/"
)

func main() {
	/* http.HandleFunc("/upload", UploadHandler)

	log.Fatal(grace.Serve(":9000", nil)) */
	vips.Resize(filename, rootDir, "results/vips/", width, height, quality)
	graphicsmagick.Resize(filename, rootDir, "results/gmagick/", width, height, quality)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
