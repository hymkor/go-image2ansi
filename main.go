package image2ansi

import (
	"context"
	"image"
	"io"
)

const blockChar = "\u2580"

type Ink interface {
	Foreground(r, g, b, a uint32, w io.Writer)
	Background(r, g, b, a uint32, w io.Writer)
	DrawBlock(w io.Writer)
	Newline(w io.Writer)
}

type defaultInk struct{}

func (defaultInk) DrawBlock(w io.Writer) {
	io.WriteString(w, blockChar)
}

func (defaultInk) Newline(w io.Writer) {
	io.WriteString(w, "\x1B[0m\n")
}

type Pen struct {
	Ink
}

type session struct {
	Ink
	r1, g1, b1, a1 uint32
	r2, g2, b2, a2 uint32
	init1, init2   bool
}

func (s *session) draw(r1, g1, b1, a1, r2, g2, b2, a2 uint32, w io.Writer) {
	if !s.init1 || s.r1 != r1 || s.g1 != g1 || s.b1 != b1 || s.a1 != a1 {
		s.Foreground(r1, g1, b1, a1, w)
		s.r1 = r1
		s.g1 = g1
		s.b1 = b1
		s.a1 = a1
		s.init1 = true
	}
	if !s.init2 || s.r2 != r2 || s.g2 != g2 || s.b2 != b2 || s.a2 != a2 {
		s.Background(r2, g2, b2, a2, w)
		s.r2 = r2
		s.g2 = g2
		s.b2 = b2
		s.a2 = a2
		s.init2 = true
	}
	s.DrawBlock(w)
}

func (p *Pen) Draw(ctx context.Context, img image.Image, rec image.Rectangle, w io.Writer) error {
	for y := rec.Min.Y; y < rec.Max.Y; y += 2 {
		s := &session{Ink: p.Ink}
		for x := rec.Min.X; x < rec.Max.X; x++ {
			var r2, g2, b2, a2 uint32
			if y+1 < rec.Max.Y {
				r2, g2, b2, a2 = img.At(x, y+1).RGBA()
			}
			r1, g1, b1, a1 := img.At(x, y).RGBA()
			s.draw(r1, g1, b1, a1, r2, g2, b2, a2, w)
		}
		s.Newline(w)
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return nil
}
