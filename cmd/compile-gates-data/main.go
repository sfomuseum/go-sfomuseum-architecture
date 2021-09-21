package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/sfomuseum/go-sfomuseum-architecture/gates"
	"io"
	"log"
	"os"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "repo://?include=properties.sfomuseum:placetype=gate", "A valid whosonfirst/go-whosonfirst-iterate URI")
	iterator_source := flag.String("iterator-source", "/usr/local/data/sfomuseum-data-architecture", "The URI containing documents to iterate.")

	target := flag.String("target", "data/gates.json", "The path to write SFO Museum gates data.")
	stdout := flag.Bool("stdout", false, "Emit SFO Museum gates data to SDOUT.")

	flag.Parse()

	ctx := context.Background()

	writers := make([]io.Writer, 0)

	fh, err := os.OpenFile(*target, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open '%s', %v", *target, err)
	}

	writers = append(writers, fh)

	if *stdout {
		writers = append(writers, os.Stdout)
	}

	wr := io.MultiWriter(writers...)

	lookup, err := gates.CompileGatesData(ctx, *iterator_uri, *iterator_source)

	if err != nil {
		log.Fatalf("Failed to compile gates data, %v", err)
	}

	enc := json.NewEncoder(wr)
	err = enc.Encode(lookup)

	if err != nil {
		log.Fatalf("Failed to marshal results, %v", err)
	}
}
