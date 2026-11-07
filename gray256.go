package image2ansi

import (
	"fmt"
	"io"
)

type gray256 struct{ defaultInk }

var Gray256 = &Pen{Ink: gray256{}}

// grayID converts RGBA values (as returned by color.Color.RGBA())
// to an xterm 256-color grayscale palette index.
func grayID(r, g, b, a uint32) uint32 {
	if a == 0 {
		return 232
	}

	// Convert to non-premultiplied RGB (0..255).
	r = r * 255 / a
	g = g * 255 / a
	b = b * 255 / a

	// Convert RGB to luminance.
	// Integer approximation of Rec. 601 luma.
	y := (r*299 + g*587 + b*114) / 1000

	// Map 0..255 to 24 grayscale levels (232..255).
	return 232 + y*23/255
}

func (gray256) Foreground(r, g, b, a uint32, w io.Writer) {
	fmt.Fprintf(w, "\x1B[38;5;%dm", grayID(r, g, b, a))
}

func (gray256) Background(r, g, b, a uint32, w io.Writer) {
	fmt.Fprintf(w, "\x1B[48;5;%dm", grayID(r, g, b, a))
}
