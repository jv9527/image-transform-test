package vips

import (
	"log"
	"time"

	"gopkg.in/h2non/bimg.v1"
)

func Resize(filename, rootDir, rootOut string, width, height, quality int) {
	start := time.Now()
	buffer, err := bimg.Read(rootDir + filename)
	if err != nil {
		log.Printf("error while read image: %v\n", err)
	}

	newImage, err := bimg.Resize(buffer, bimg.Options{
		Width:   width,
		Height:  height,
		Quality: quality,
		//Crop:    true,
	})
	if err != nil {
		log.Printf("error while Resize image: %v\n", err)
		return
	}

	if err := bimg.Write(rootOut+filename, newImage); err != nil {
		log.Printf("error while write image: %v\n", err)
	}

	log.Printf("Vips resize time to %dx%d with quaility %d : %.4f", width, height, quality, time.Since(start).Seconds())
}
