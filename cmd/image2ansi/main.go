package main

import (
	"context"
	"errors"
	"flag"
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

	"github.com/hymkor/go-windows1x-virtualterminal"

	"github.com/hymkor/go-image2ansi"
	"github.com/hymkor/go-image2ansi/clipping"
)

var (
	flagTrueColor = flag.Bool("24", false, "Use 24-bit true color")
	flagGray      = flag.Bool("gray", false, "Use 256-color grayscale")

	flagX      = flag.Int("x", 0, "Image start X (0-based)")
	flagY      = flag.Int("y", 0, "Image start Y (0-based)")
	flagWidth  = flag.Int("w", 0, "Maximum output width")
	flagHeight = flag.Int("h", 0, "Maximum output height")

	flagGifFromStdin = flag.Bool("gif-from-stdin", false, "Read a GIF image from standard input")
	flagPngFromStdin = flag.Bool("png-from-stdin", false, "Read a PNG image from standard input")
	flagJpgFromStdin = flag.Bool("jpg-from-stdin", false, "Read a JPEG image from standard input")
)

func drawFile(ctx context.Context, fname string, clipping func(image.Rectangle) image.Rectangle, pen *image2ansi.Pen, w io.Writer) error {
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
		return errors.New("not support format")
	}
	if err != nil {
		return err
	}
	return pen.Draw(ctx, img, clipping(img.Bounds()), w)
}

func clippingTermWidth(fd int) func(image.Rectangle) image.Rectangle {
	width, _, err := term.GetSize(fd)
	if err != nil {
		return clipping.Identity
	}
	return clipping.MaxWidth(width)
}

func mains(args []string) error {
	clipping1 := clipping.From(*flagX, *flagY)
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		closer, err := virtualterminal.EnableStdout()
		if err != nil {
			return err
		}
		defer closer()
		if *flagWidth <= 0 {
			clipping1 = clipping.Compose(clipping1, clippingTermWidth(fd))
		}
	}
	if *flagWidth > 0 {
		clipping1 = clipping.Compose(clipping1, clipping.MaxWidth(*flagWidth))
	}
	if *flagHeight > 0 {
		clipping1 = clipping.Compose(clipping1, clipping.MaxHeight(*flagHeight))
	}
	var pen *image2ansi.Pen = image2ansi.Color256
	if *flagTrueColor {
		pen = image2ansi.TrueColor
	} else if *flagGray {
		pen = image2ansi.Gray256
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	for _, arg := range args {
		if fnames, err := filepath.Glob(arg); err != nil {
			if len(args) > 1 {
				fmt.Fprintf(os.Stdout, "%s:\n", arg)
			}
			err := drawFile(ctx, arg, clipping1, pen, os.Stdout)
			if err != nil {
				return fmt.Errorf("%s: %w", arg, err)
			}
		} else {
			for _, fname := range fnames {
				if len(args) > 1 || len(fnames) > 1 {
					fmt.Fprintf(os.Stdout, "%s:\n", fname)
				}
				err := drawFile(ctx, fname, clipping1, pen, os.Stdout)
				if err != nil {
					return fmt.Errorf("%s: %w", fname, err)
				}
			}
		}
	}
	if *flagGifFromStdin {
		img, err := gif.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return pen.Draw(ctx, img, clipping1(img.Bounds()), os.Stdout)
	}
	if *flagPngFromStdin {
		img, err := png.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return pen.Draw(ctx, img, clipping1(img.Bounds()), os.Stdout)
	}
	if *flagJpgFromStdin {
		img, err := jpeg.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return pen.Draw(ctx, img, clipping1(img.Bounds()), os.Stdout)
	}
	if len(args) <= 0 {
		flag.Usage()
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, "Interrupted")
			os.Exit(130)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
