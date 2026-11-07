package image2ansi

import (
	"io"
)

type monoChrome struct{}

var MonoChrome = &Pen{Ink: monoChrome{}}

func isWhite(r, g, b, a uint32) bool {
	if a == 0 {
		return false
	}

	r = r * 255 / a
	g = g * 255 / a
	b = b * 255 / a

	y := (r*299 + g*587 + b*114) / 1000
	return y >= 128
}

func (monoChrome) Foreground(r, g, b, a uint32, w io.Writer) {
	if isWhite(r, g, b, a) {
		io.WriteString(w, "\x1B[37m")
	} else {
		io.WriteString(w, "\x1B[30m")
	}
}

func (monoChrome) Background(r, g, b, a uint32, w io.Writer) {
	if isWhite(r, g, b, a) {
		io.WriteString(w, "\x1B[47m")
	} else {
		io.WriteString(w, "\x1B[40m")
	}
}

func (monoChrome) Newline(w io.Writer) {
	io.WriteString(w, "\x1B[0m\n")
}
