package lib2

import (
	"github.com/gographics/gmagick"
)

func Resize() {
	mw := gmagick.NewMagickWand()
	defer mw.Destroy()

}
