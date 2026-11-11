package gifanimation

import (
	"image"
	"image/color"
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

	bg := color.Color(color.Transparent)
	if palette, ok := img.Config.ColorModel.(color.Palette); ok {
		if int(img.BackgroundIndex) < len(palette) {
			bg = palette[img.BackgroundIndex]
		}
	}
	background := image.NewUniform(bg)

	return func(yield func(int, *View) bool) {
		var count int
		switch img.LoopCount {
		case 0:
			count = 32767
		case -1:
			count = 1
		default:
			count = img.LoopCount + 1
		}
		var prev *image.Paletted
		var prevDisposal byte
		for i := 0; i < count; i++ {
			for j, frame := range img.Image {
				if prev != nil && prevDisposal == gif.DisposalBackground {
					draw.Draw(
						canvas,
						prev.Bounds(),
						background,
						image.Point{},
						draw.Src,
					)
				}
				v := &View{
					canvas: canvas,
					frame:  frame,
					delay:  img.Delay[j],
				}
				if !yield(j, v) {
					return
				}
				prev = frame
				if j < len(img.Disposal) {
					prevDisposal = img.Disposal[j]
				} else {
					prevDisposal = 0
				}
			}
		}
	}
}
