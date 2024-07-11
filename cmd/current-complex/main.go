package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/sfomuseum/go-sfomuseum-architecture/campus"
	"github.com/whosonfirst/go-reader"
)

func main() {

	var iterator_uri string
	var output_mode string
	var reader_uri string
	var verbose bool

	flag.StringVar(&iterator_uri, "iterator-uri", "repo://", "...")
	flag.StringVar(&output_mode, "output-mode", "json", "...")
	flag.StringVar(&reader_uri, "reader-uri", "repo:///usr/local/data/sfomuseum-data-architecture", "...")
	flag.BoolVar(&verbose, "verbose", false, "...")

	flag.Parse()

	ctx := context.Background()

	if verbose {
		// set slog.LogLevel here
		// why is this so hard?
	}

	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stdout)
	wr := io.MultiWriter(writers...)

	paths := flag.Args()

	c, err := campus.MostRecentComplexWithIterator(ctx, iterator_uri, paths...)

	if err != nil {
		log.Fatalf("Failed to derive most recent complex, %v", err)
	}

	switch output_mode {
	case "json":

		err := c.AsJSON(ctx, wr)

		if err != nil {
			log.Fatalf("Failed to encode complex, %v", err)
		}

	case "tree":

		r, err := reader.NewReader(ctx, reader_uri)

		if err != nil {
			log.Fatalf("Failed to create new reader, %v", err)
		}

		err = c.AsTree(ctx, r, wr)

		if err != nil {
			log.Fatalf("Failed to render complex as tree, %v", err)
		}

	default:
		log.Fatalf("Invalid or unsupported output mode, %v", err)
	}
}
