package gifanimation

import (
	"image"
	"image/draw"
	"image/gif"
	"time"
)

type View struct {
	canvas draw.Image
	frame  *image.Paletted
	delay  int
}

func (v *View) Delay() time.Duration {
	return time.Duration(v.delay) * time.Second / 100
}

func (v *View) Image() image.Image {
	draw.Draw(
		v.canvas,
		v.frame.Bounds(),
		v.frame,
		v.frame.Bounds().Min,
		draw.Over,
	)
	return v.canvas
}

func Each(img *gif.GIF) func(yield func(int, *View) bool) {
	canvas := image.NewRGBA(
		image.Rect(
			0, 0,
			img.Config.Width,
			img.Config.Height))

	return func(yield func(int, *View) bool) {
		for i, frame := range img.Image {
			v := &View{
				canvas: canvas,
				frame:  frame,
				delay:  img.Delay[i],
			}
			if !yield(i, v) {
				break
			}
		}
	}
}
