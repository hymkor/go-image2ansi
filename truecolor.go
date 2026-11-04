package image2ansi

import (
	"fmt"
	"io"
)

const baseA = 255

type trueColor struct{}

var TrueColor = trueColor{}

func (trueColor) Foreground(r, g, b, a uint32, w io.Writer) {
	if a == 0 {
		r = 0
		g = 0
		b = 0
	} else {
		r = r * baseA / a
		g = g * baseA / a
		b = b * baseA / a
	}
	fmt.Fprintf(w, "\x1B[38;2;%d;%d;%dm", r, g, b)
}

func (trueColor) Background(r, g, b, a uint32, w io.Writer) {
	if a == 0 {
		r = 0
		g = 0
		b = 0
	} else {
		r = r * baseA / a
		g = g * baseA / a
		b = b * baseA / a
	}
	fmt.Fprintf(w, "\x1B[48;2;%d;%d;%dm", r, g, b)
}
