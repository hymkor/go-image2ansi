package image2ansi

import (
	"context"
	"image"
	"io"
)

const blockChar = "\u2580"

type Pen interface {
	Foreground(r, g, b, a uint32, w io.Writer)
	Background(r, g, b, a uint32, w io.Writer)
}

type penAndColor struct {
	Pen
	r1, g1, b1, a1 uint32
	r2, g2, b2, a2 uint32
	init1, init2   bool
}

func (p *penAndColor) draw(r1, g1, b1, a1, r2, g2, b2, a2 uint32, w io.Writer) {
	if !p.init1 || p.r1 != r1 || p.g1 != g1 || p.b1 != b1 || p.a1 != a1 {
		p.Foreground(r1, g1, b1, a1, w)
		p.r1 = r1
		p.g1 = g1
		p.b1 = b1
		p.a1 = a1
		p.init1 = true
	}
	if !p.init2 || p.r2 != r2 || p.g2 != g2 || p.b2 != b2 || p.a2 != a2 {
		p.Background(r2, g2, b2, a2, w)
		p.r2 = r2
		p.g2 = g2
		p.b2 = b2
		p.a2 = a2
		p.init2 = true
	}
	io.WriteString(w, blockChar)
}

func Draw(ctx context.Context, img image.Image, rec image.Rectangle, pen Pen, w io.Writer) error {
	for y := rec.Min.Y; y < rec.Max.Y; y += 2 {
		pc := &penAndColor{Pen: pen}
		for x := rec.Min.X; x < rec.Max.X; x++ {
			var r2, g2, b2, a2 uint32
			if y+1 < rec.Max.Y {
				r2, g2, b2, a2 = img.At(x, y+1).RGBA()
			}
			r1, g1, b1, a1 := img.At(x, y).RGBA()
			pc.draw(r1, g1, b1, a1, r2, g2, b2, a2, w)
		}
		io.WriteString(w, "\x1B[0m\n")
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return nil
}

func NoClipping(img image.Image) image.Rectangle {
	return img.Bounds()
}

func ClippingWidth(width int) func(image.Image) image.Rectangle {
	return func(img image.Image) image.Rectangle {
		rec := img.Bounds()
		if rec.Max.X-rec.Min.X > width {
			rec.Max.X = rec.Min.X + width
		}
		return rec
	}
}
