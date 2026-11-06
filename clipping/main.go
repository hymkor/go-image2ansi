package clipping

import (
	"image"
)

func Identity(rec image.Rectangle) image.Rectangle {
	return rec
}

func MaxWidth(width int) func(rec image.Rectangle) image.Rectangle {
	return func(rec image.Rectangle) image.Rectangle {
		if rec.Max.X-rec.Min.X > width {
			rec.Max.X = rec.Min.X + width
		}
		return rec
	}
}

func MaxHeight(height int) func(image.Rectangle) image.Rectangle {
	return func(rec image.Rectangle) image.Rectangle {
		if rec.Max.Y-rec.Min.Y > height {
			rec.Max.Y = rec.Min.Y + height
		}
		return rec
	}
}

func From(x, y int) func(rec image.Rectangle) image.Rectangle {
	return func(rec image.Rectangle) image.Rectangle {
		rec.Min.X += x
		rec.Min.Y += y
		rec.Max.X += x
		rec.Max.Y += y
		return rec
	}
}

func Compose(f ...func(image.Rectangle) image.Rectangle) func(image.Rectangle) image.Rectangle {
	return func(rec image.Rectangle) image.Rectangle {
		for _, f1 := range f {
			if f1 != nil {
				rec = f1(rec)
			}
		}
		return rec
	}
}
