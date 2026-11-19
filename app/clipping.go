package app

import (
	"image"

	"golang.org/x/term"

	"github.com/hymkor/go-image2ansi/clipping"
)

func clippingTermWidth(fd int) func(image.Rectangle) image.Rectangle {
	width, _, err := term.GetSize(fd)
	if err != nil {
		return clipping.Identity
	}
	return clipping.MaxWidth(width)
}

func clippingTermHeight(fd int) func(image.Rectangle) image.Rectangle {
	_, height, err := term.GetSize(fd)
	if err != nil {
		return clipping.Identity
	}
	return clipping.MaxHeight(height*2 - 1)
}
