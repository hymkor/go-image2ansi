package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hymkor/struct2flag"

	"github.com/hymkor/go-image2ansi/app"
)

var version string

func main() {
	var cfg app.Config
	struct2flag.BindDefault(&cfg)
	flag.Parse()
	if err := cfg.Run(context.Background(), os.Stdin, os.Stdout, flag.Args()); err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, "Interrupted")
			os.Exit(130)
		} else if errors.Is(err, app.ErrNothingToDo) {
			fmt.Fprintf(os.Stderr, "%s %s-%s-%s\n",
				filepath.Base(os.Args[0]),
				version,
				runtime.GOOS,
				runtime.GOARCH)
			flag.Usage()
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
