package app

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/term"

	"github.com/biessek/golang-ico"
	"rsc.io/qr"

	"github.com/hymkor/go-windows1x-virtualterminal"

	"github.com/hymkor/go-image2ansi"
	"github.com/hymkor/go-image2ansi/clipping"
)

var (
	ErrNotSupportFormat = errors.New("not support format")
	ErrNothingToDo      = errors.New("nothing to do")
)

type Config struct {
	TrueColor   bool `flag:"24,Use 24-bit true color"`
	Gray        bool `flag:"gray,Use 256-color grayscale"`
	Mono        bool `flag:"mono,Use monochrome output with Unicode block characters (no ANSI escape sequences"`
	MonoReverse bool `flag:"mono-reverse,Reverse monochrome output for light backgrounds"`
	Animation   bool `flag:"animation,Enable experimental animated GIF support"`

	X      int `flag:"x,Image start X (0-based)"`
	Y      int `flag:"y,Image start Y (0-based)"`
	Width  int `flag:"w,Maximum output width"`
	Height int `flag:"h,Maximum output height"`

	GifFromStdin bool `flag:"gif-from-stdin,Read a GIF image from standard input"`
	PngFromStdin bool `flag:"png-from-stdin,Read a PNG image from standard input"`
	JpgFromStdin bool `flag:"jpg-from-stdin,Read a JPEG image from standard input"`

	QR string `flag:"qr,Display the given string as a QR code"`
}

func (cfg *Config) drawFile(ctx context.Context, fname string, clipping func(image.Rectangle) image.Rectangle, pen *image2ansi.Pen, w io.Writer) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fd.Close()

	ext := filepath.Ext(fname)

	var img image.Image
	if strings.EqualFold(ext, ".png") {
		img, err = png.Decode(fd)
	} else if strings.EqualFold(ext, ".jpg") {
		img, err = jpeg.Decode(fd)
	} else if strings.EqualFold(ext, ".jpeg") {
		img, err = jpeg.Decode(fd)
	} else if strings.EqualFold(ext, ".gif") {
		if cfg.Animation {
			gifimg, err := gif.DecodeAll(fd)
			if err != nil {
				return fmt.Errorf(`"image/gif".DecodeAll: %w`, err)
			}
			return gifDecode(ctx, gifimg, clipping, pen, w)
		}
		img, err = gif.Decode(fd)
	} else if strings.EqualFold(ext, ".bmp") {
		img, err = bmp.Decode(fd)
	} else if strings.EqualFold(ext, ".ico") {
		images, err := ico.DecodeAll(fd)
		if err != nil {
			return err
		}
		for i, img := range images {
			if i > 0 {
				w.Write([]byte{'\n'})
			}
			err := pen.Draw(ctx, img, clipping(img.Bounds()), w)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		return ErrNotSupportFormat
	}
	if err != nil {
		return err
	}
	return pen.Draw(ctx, img, clipping(img.Bounds()), w)
}

func (cfg Config) Run(ctx context.Context, in, out *os.File, args []string) error {
	clipping1 := clipping.From(cfg.X, cfg.Y)
	fd := int(out.Fd())
	if term.IsTerminal(fd) {
		closer, err := virtualterminal.EnableStdout()
		if err != nil {
			return err
		}
		defer closer()
		if cfg.Width <= 0 {
			clipping1 = clipping.Compose(clipping1, clippingTermWidth(fd))
		}
		if cfg.Animation && cfg.Height <= 0 {
			clipping1 = clipping.Compose(clipping1, clippingTermHeight(fd))
		}
	}
	if cfg.Width > 0 {
		clipping1 = clipping.Compose(clipping1, clipping.MaxWidth(cfg.Width))
	}
	if cfg.Height > 0 {
		clipping1 = clipping.Compose(clipping1, clipping.MaxHeight(cfg.Height))
	}
	var pen *image2ansi.Pen = image2ansi.Color256
	if cfg.TrueColor {
		pen = image2ansi.TrueColor
	} else if cfg.Gray {
		pen = image2ansi.Gray256
	} else if cfg.Mono {
		pen = image2ansi.MonoChrome
	} else if cfg.MonoReverse {
		pen = image2ansi.MonoChromeReverse
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	for _, arg := range args {
		if fnames, err := filepath.Glob(arg); err != nil {
			if len(args) > 1 {
				fmt.Fprintf(out, "%s:\n", arg)
			}
			err := cfg.drawFile(ctx, arg, clipping1, pen, out)
			if err != nil {
				return fmt.Errorf("%s: %w", arg, err)
			}
		} else {
			for _, fname := range fnames {
				if len(args) > 1 || len(fnames) > 1 {
					fmt.Fprintf(out, "%s:\n", fname)
				}
				err := cfg.drawFile(ctx, fname, clipping1, pen, out)
				if err != nil {
					return fmt.Errorf("%s: %w", fname, err)
				}
			}
		}
	}
	if cfg.GifFromStdin {
		if cfg.Animation {
			gifimg, err := gif.DecodeAll(os.Stdin)
			if err != nil {
				return err
			}
			return gifDecode(ctx, gifimg, clipping1, pen, out)
		}
		img, err := gif.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return pen.Draw(ctx, img, clipping1(img.Bounds()), out)
	}
	if cfg.PngFromStdin {
		img, err := png.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return pen.Draw(ctx, img, clipping1(img.Bounds()), out)
	}
	if cfg.JpgFromStdin {
		img, err := jpeg.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return pen.Draw(ctx, img, clipping1(img.Bounds()), out)
	}
	if cfg.QR != "" {
		code, err := qr.Encode(cfg.QR, qr.L)
		if err != nil {
			return err
		}
		img := code.Image()
		d := code.Size + 1
		return pen.Draw(ctx, img, clipping1(
			image.Rect(0, 0, d, d)), out)
	}
	if len(args) <= 0 {
		return ErrNothingToDo
	}
	return nil
}
