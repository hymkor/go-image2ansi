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

	"golang.org/x/term"

	"github.com/hymkor/go-windows1x-virtualterminal"

	"github.com/hymkor/go-image2ansi"
)

var (
	flagTrueColor = flag.Bool("24", false, "Use 24-bit true color")
	flagWidth     = flag.Int("w", 0, "Maximum output width")

	flagGifFromStdin = flag.Bool("gif-from-stdin", false, "Read a GIF image from standard input")
	flagPngFromStdin = flag.Bool("png-from-stdin", false, "Read a PNG image from standard input")
	flagJpgFromStdin = flag.Bool("jpg-from-stdin", false, "Read a JPEG image from standard input")
)

func drawFile(ctx context.Context, fname string, clipping func(image.Image) image.Rectangle, pen image2ansi.Pen, w io.Writer) error {
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
	} else {
		return errors.New("not support format")
	}
	if err != nil {
		return err
	}
	return image2ansi.Draw(ctx, img, clipping(img), pen, w)
}

func clippingTermWidth(fd int) func(image.Image) image.Rectangle {
	width, _, err := term.GetSize(fd)
	if err != nil {
		return image2ansi.NoClipping
	}
	return image2ansi.ClippingWidth(width)
}

func mains(args []string) error {
	clipping := image2ansi.NoClipping
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		closer, err := virtualterminal.EnableStdout()
		if err != nil {
			return err
		}
		defer closer()
		clipping = clippingTermWidth(fd)
	}
	if *flagWidth > 0 {
		clipping = image2ansi.ClippingWidth(*flagWidth)
	}
	var pen image2ansi.Pen = image2ansi.Color256
	if *flagTrueColor {
		pen = image2ansi.TrueColor
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	for _, fname := range args {
		err := drawFile(ctx, fname, clipping, pen, os.Stdout)
		if err != nil {
			return fmt.Errorf("%s: %w", fname, err)
		}
	}
	if *flagGifFromStdin {
		img, err := gif.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return image2ansi.Draw(ctx, img, clipping(img), pen, os.Stdout)
	}
	if *flagPngFromStdin {
		img, err := png.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return image2ansi.Draw(ctx, img, clipping(img), pen, os.Stdout)
	}
	if *flagJpgFromStdin {
		img, err := jpeg.Decode(os.Stdin)
		if err != nil {
			return err
		}
		return image2ansi.Draw(ctx, img, clipping(img), pen, os.Stdout)
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
