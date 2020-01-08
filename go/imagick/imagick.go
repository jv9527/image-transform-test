package imagick

import (
	"log"
	"time"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func Convert(filename, rootDir, rootOut string) {
	start := time.Now()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// pw := imagick.NewPixelWand()
	// pw.SetRed(255)
	// pw.SetGreen(255)
	//pw.SetBlue(255)
	//mw.SetImageBackgroundColor(pw)

	mw.ReadImage(rootDir + filename)
	mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN)
	mw.SetFormat("jpg")

	mw.WriteImage(rootOut + filename + ".jpg")
	log.Printf("Imagick convert time to %.4f", time.Since(start).Seconds())
}
