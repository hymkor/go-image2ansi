package image2ansi

import (
	"fmt"
	"io"
)

type color256 struct{}

var Color256 = color256{}

func id(r, g, b, a uint32) uint32 {
	if a == 0 {
		r = 0
		g = 0
		b = 0
	} else {
		r = r * 5 / a
		g = g * 5 / a
		b = b * 5 / a
	}
	return 16 + r*36 + g*6 + b
}

func (color256) Foreground(r, g, b, a uint32, w io.Writer) {
	fmt.Fprintf(w, "\x1B[38;5;%dm", id(r, g, b, a))
}

func (color256) Background(r, g, b, a uint32, w io.Writer) {
	fmt.Fprintf(w, "\x1B[48;5;%dm", id(r, g, b, a))
}

func (color256) Newline(w io.Writer) {
	io.WriteString(w, "\x1B[0m\n")
}
