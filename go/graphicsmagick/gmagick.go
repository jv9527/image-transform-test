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

	mw.ReadImage(rootDir + filename)
	/* if err := mw.ReadImageBlob(orig); err != nil {
		return nil, err
	} */

	// Strip Image
	/* if err := mw.StripImage(); err != nil {
		return nil, err
	} */

	// Resize Image
	if err := mw.ResizeImage(uint(width), uint(height), gmagick.FILTER_LANCZOS, 1); err != nil {
		log.Println(err)
	}

	mw.SetCompressionQuality(uint(quality))

	// If Quality is Defined, Set Compression Quality
	/* if quality > 0 {
		if quality > 100 {
			quality = 100
		}

		// Set Compression Quality
		if err := mw.SetCompressionQuality(uint(quality)); err != nil {
			return nil, err
		}
	} */

	mw.WriteImage(rootOut + filename)
	log.Printf("Gmagick resize time to %dx%d with quaility %d : %.4f", width, height, quality, time.Since(start).Seconds())
}

func Convert(filename, rootDir, rootOut string) {
	start := time.Now()
	mw := gmagick.NewMagickWand()
	defer mw.Destroy()

	pw := gmagick.NewPixelWand()
	pw.SetRed(255)
	pw.SetGreen(255)
	pw.SetBlue(255)

	mw.SetImageBackgroundColor(pw)
	mw.ReadImage(rootDir + filename)

	//mw.SetFormat("jpg")

	mw.WriteImage(rootOut + filename + ".jpg")
	log.Printf("Gmagick convert time to %.4f", time.Since(start).Seconds())
}
