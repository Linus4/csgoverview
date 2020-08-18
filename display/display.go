package display

import (
	"github.com/kbinani/screenshot"
	"image"
)

func GetBounds() image.Rectangle {
	return screenshot.GetDisplayBounds(0)
}

func GetWidth() int {
	bounds := GetBounds()
	return bounds.Dx()
}

func GetHeight() int {
	bounds := GetBounds()
	return bounds.Dy()
}
