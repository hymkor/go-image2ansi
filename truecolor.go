package image2ansi

import (
	"fmt"
	"io"
)

type trueColor struct{}

var TrueColor = trueColor{}

func rgb(r, g, b, a uint32) (uint32, uint32, uint32) {
	const baseA = 255

	if a == 0 {
		r = 0
		g = 0
		b = 0
	} else {
		r = r * baseA / a
		g = g * baseA / a
		b = b * baseA / a
	}
	return r, g, b
}

func (trueColor) Foreground(r, g, b, a uint32, w io.Writer) {
	r, g, b = rgb(r, g, b, a)
	fmt.Fprintf(w, "\x1B[38;2;%d;%d;%dm", r, g, b)
}

func (trueColor) Background(r, g, b, a uint32, w io.Writer) {
	r, g, b = rgb(r, g, b, a)
	fmt.Fprintf(w, "\x1B[48;2;%d;%d;%dm", r, g, b)
}
