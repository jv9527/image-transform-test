package vips

import (
	"fmt"
	"log"
	"os"
	"time"

	bimg "gopkg.in/h2non/bimg.v1"
)

func Resize(filename, rootDir, rootOut string, width, height, quality int) {
	start := time.Now()
	buffer, err := bimg.Read(rootDir + filename)
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

	bimg.Write(rootOut+filename, newImage)
	log.Printf("Vips resize time to %dx%d with quaility %d : %.4f", width, height, quality, time.Since(start).Seconds())
}

func Convert(filename, rootDir, rootOut string) {
	start := time.Now()
	buffer, err := bimg.Read(rootDir + filename)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	bgWhite, err := bimg.NewImage(buffer).Process(bimg.Options{
		Background: bimg.Color{
			R: 255,
			G: 255,
			B: 255,
		},
		StripMetadata: true,
	})

	if err != nil {
		log.Println(err)
	}

	newImg, err := bimg.NewImage(bgWhite).Convert(bimg.JPEG)

	if err != nil {
		log.Println(err)
	}

	bimg.Write(rootOut+filename+".jpg", newImg)
	log.Printf("Vips convert time %.4f", time.Since(start).Seconds())
}
