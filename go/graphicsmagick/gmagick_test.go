package graphicsmagick

import (
	"testing"
)

const (
	width    = 200
	height   = 200
	quality  = 70
	filename = "f1.png"
)

func BenchmarkResizeGmagick(b *testing.B) {
	Resize(filename, "../../images/png/", "../../images/gmagick/", width, height, quality)
}
