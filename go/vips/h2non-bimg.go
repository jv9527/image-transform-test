package vips

import (
	"fmt"
	"log"
	"os"
	"time"

	bimg "gopkg.in/h2non/bimg.v1"
)

func Resize(filename string, width, height, quality int) {
	start := time.Now()
	buffer, err := bimg.Read("images/" + filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newImage, err := bimg.Resize(buffer, bimg.Options{
		Width:   width,
		Height:  height,
		Quality: quality,
		Crop:    true,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	bimg.Write("results/vips/"+filename, newImage)
	log.Printf("Vips resize time to %dx%d with quaility %d : %.4f", width, height, quality, time.Since(start).Seconds())
}
