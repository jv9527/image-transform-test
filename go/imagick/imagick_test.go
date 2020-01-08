package imagick

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
	Convert(filename, "../../images/png/", "../../images/gmagick/")
}
