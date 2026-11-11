package main

import (
	"context"
	"fmt"
	"image"
	"image/gif"
	"io"
	"time"

	"github.com/hymkor/go-image2ansi"
	"github.com/hymkor/go-image2ansi/gifanimation"
)

func gifDecode(ctx context.Context, img *gif.GIF, clip func(image.Rectangle) image.Rectangle, pen *image2ansi.Pen, w io.Writer) error {
	line := 0
	defer pen.Newline(w)

	var delay <-chan (time.Time)
	for _, v := range gifanimation.Each(img) {
		if delay != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-delay:
			}
		}
		delay = time.After(v.Delay())
		if line > 1 {
			fmt.Fprintf(w, "\x1B[0m\x1B[%dF", line-1)
		}
		view := v.Image()
		rec := clip(view.Bounds())
		err := pen.DrawNoNewline(ctx, view, rec, w)
		if err != nil {
			return err
		}
		line = (rec.Max.Y - rec.Min.Y + 1) / 2
	}
	return nil
}
