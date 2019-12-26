package graphicsmagick

import (
	"log"
	"time"

	"github.com/gographics/gmagick"
)

func Resize(filename, rootDir, rootOut string, width, height, quality int) {
	start := time.Now()
	mw := gmagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImage(rootDir + filename); err != nil {
		log.Printf("error while read image: %v\n", err)
		return
	}

	// Strip Image
	if err := mw.StripImage(); err != nil {
		log.Printf("error while strip image: %v\n", err)
		return
	}

	// Resize Image
	if err := mw.ResizeImage(uint(width), uint(height), gmagick.FILTER_LANCZOS, 1); err != nil {
		log.Printf("error while resize image: %v\n", err)
		return
	}

	if err := mw.SetCompressionQuality(uint(quality)); err != nil {
		log.Printf("error while set compression quality: %v\n", err)
	}

	// write image
	if err := mw.WriteImage(rootOut + filename); err != nil {
		log.Printf("error while write image: %v\n", err)
		return
	}

	log.Printf("Gmagick resize time to %dx%d with quaility %d : %.4f", width, height, quality, time.Since(start).Seconds())
}
