package image2ansi

import (
	"fmt"
	"io"
)

type color256 struct{}

var Color256 = color256{}

func (color256) Foreground(r, g, b, a uint32, w io.Writer) {
	if a == 0 {
		r = 0
		g = 0
		b = 0
	} else {
		r = r * 5 / a
		g = g * 5 / a
		b = b * 5 / a
	}
	fmt.Fprintf(w, "\x1B[38;5;%dm", 16+r*36+g*6+b)
}

func (color256) Background(r, g, b, a uint32, w io.Writer) {
	if a == 0 {
		r = 0
		g = 0
		b = 0
	} else {
		r = r * 5 / a
		g = g * 5 / a
		b = b * 5 / a
	}
	fmt.Fprintf(w, "\x1B[48;5;%dm", 16+r*36+g*6+b)
}
