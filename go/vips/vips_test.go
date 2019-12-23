package vips

import (
	"testing"
)

const (
	width    = 200
	height   = 200
	quality  = 70
	filename = "file2.jpg"
)

func BenchmarkResizeVips(b *testing.B) {
	Resize(filename, "../../images/", "../../images/vips/", width, height, quality)
}
