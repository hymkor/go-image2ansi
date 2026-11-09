package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"time"

	"github.com/hymkor/go-image2ansi"
)

func gifDecode(ctx context.Context, img *gif.GIF, clip func(image.Rectangle) image.Rectangle, pen *image2ansi.Pen, w io.Writer) error {
	// fmt.Fprintf(w, "canvas: (%d,%d)-(%d,%d)\n", img.Config.Width, img.Config.Height)
	rec := image.Rect(
		0, 0,
		img.Config.Width,
		img.Config.Height)
	canvas := image.NewRGBA(rec)
	rec = clip(rec)
	up := 0
	var delay <-chan time.Time
	for i, frame := range img.Image {
		// fmt.Fprintf(w, "%d: (%d,%d)-(%d,%d)\n", i, b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
		if delay != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-delay:
				break
			}
		}
		delay = time.After(time.Duration(img.Delay[i]) * time.Second / 100)
		if up > 0 {
			fmt.Fprintf(w, "\x1B[%dA", up)
		}
		draw.Draw(
			canvas,
			frame.Bounds(),
			frame,
			frame.Bounds().Min,
			draw.Over,
		)
		err := pen.Draw(ctx, canvas, rec, w)
		if err != nil {
			return err
		}
		up = (rec.Max.Y - rec.Min.Y + 1) / 2
	}
	return nil
}
