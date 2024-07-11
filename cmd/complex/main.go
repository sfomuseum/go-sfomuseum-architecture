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
	var complex_id int64
	var dsn string

	flag.StringVar(&iterator_uri, "iterator-uri", "repo://", "...")
	flag.StringVar(&output_mode, "output-mode", "json", "...")
	flag.StringVar(&reader_uri, "reader-uri", "repo:///usr/local/data/sfomuseum-data-architecture", "...")
	flag.BoolVar(&verbose, "verbose", false, "...")
	flag.Int64Var(&complex_id, "complex-id", 0, "If 0 then the most recent (current) complex ID will be used.")
	flag.StringVar(&dsn, "dsn", ":memory:", "...")

	flag.Parse()

	ctx := context.Background()

	if verbose {
		// set slog.LogLevel here
		// why is this so hard?
	}

	paths := flag.Args()

	aa_db, err := campus.NewDatabaseWithIterator(ctx, dsn, iterator_uri, paths...)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	db_conn, err := aa_db.Conn()

	if err != nil {
		log.Fatalf("Failed to create database connection, %v", err)
	}

	c, err := campus.DeriveComplex(ctx, db_conn, complex_id)

	if err != nil {
		log.Fatalf("Failed to derive complex, %v", err)
	}

	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stdout)
	wr := io.MultiWriter(writers...)

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

		err = c.AsTree(ctx, r, wr, 0)

		if err != nil {
			log.Fatalf("Failed to render complex as tree, %v", err)
		}

	default:
		log.Fatalf("Invalid or unsupported output mode, %v", err)
	}
}
