package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/sfomuseum/go-sfomuseum-architecture/campus"
	"io"
	"log"
	"os"
)

func main() {

	mode := flag.String("mode", "repo://", "...")
	flag.Parse()

	ctx := context.Background()

	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stdout)
	wr := io.MultiWriter(writers...)

	paths := flag.Args()

	c, err := campus.MostRecentComplexWithIterator(ctx, *mode, paths...)

	if err != nil {
		log.Fatalf("Failed to derive most recent complex, %v", err)
	}

	enc := json.NewEncoder(wr)
	err = enc.Encode(c)

	if err != nil {
		log.Fatalf("Failed to encode complex, %v", err)
	}
}
