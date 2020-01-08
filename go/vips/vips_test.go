package vips

import (
	"testing"
)

const (
	width   = 200
	height  = 200
	quality = 70
	//filename = "file2.jpg"
	filename = "f1.png"
)

func BenchmarkResizeVips(b *testing.B) {
	Resize(filename, "../../images/png/", "../../images/vips/", width, height, quality)
	Convert(filename, "../../images/png/", "../../images/vips/")
}
