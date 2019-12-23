package main

import (
	"log"

	"github.com/image-transform-test/go/lib1"
	grace "gopkg.in/tokopedia/grace.v1"

	"net/http"
)

func main() {
	http.HandleFunc("/upload", UploadHandler)

	log.Fatal(grace.Serve(":9000", nil))
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	lib1.Resize()
	w.WriteHeader(http.StatusOK)
}
