package graphicsmagick

import (
	"testing"
)

const (
	width    = 200
	height   = 200
	quality  = 70
	filename = "file2.jpg"
)

func BenchmarkResizeGmagick(b *testing.B) {
	Resize(filename, "../../images/", "../../images/gmagick/", width, height, quality)
}
