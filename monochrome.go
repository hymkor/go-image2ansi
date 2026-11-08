package image2ansi

import (
	"io"
)

type monoChrome struct {
	value   int
	pattern []string
}

const (
	space          = " "
	upperHalfBlock = "\u2580"
	lowerHalfBlock = "\u2584"
	fullBlock      = "\u2588"
)

var (
	MonoChrome = &Pen{
		Ink: &monoChrome{
			pattern: []string{space, upperHalfBlock, lowerHalfBlock, fullBlock},
		},
	}

	MonoChromeReverse = &Pen{
		Ink: &monoChrome{
			pattern: []string{fullBlock, lowerHalfBlock, upperHalfBlock, space},
		},
	}
)

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

func (m *monoChrome) Foreground(r, g, b, a uint32, w io.Writer) {
	if isWhite(r, g, b, a) {
		m.value |= 1
	} else {
		m.value &^= 1
	}
}

func (m *monoChrome) Background(r, g, b, a uint32, w io.Writer) {
	if isWhite(r, g, b, a) {
		m.value |= 2
	} else {
		m.value &^= 2
	}
}

func (m *monoChrome) DrawBlock(w io.Writer) {
	io.WriteString(w, m.pattern[m.value])
}

func (m *monoChrome) Newline(w io.Writer) {
	io.WriteString(w, "\n")
}
