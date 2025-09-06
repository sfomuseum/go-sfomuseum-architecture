package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	_ "github.com/sfomuseum/go-sfomuseum-architecture/galleries"
	_ "github.com/sfomuseum/go-sfomuseum-architecture/gates"
	_ "github.com/sfomuseum/go-sfomuseum-architecture/terminals"

	"github.com/sfomuseum/go-sfomuseum-architecture"
)

func main() {

	lookup_uri := flag.String("lookup-uri", "", "Valid options are: gates://, galleries://, terminals://")

	flag.Parse()

	ctx := context.Background()
	lookup, err := architecture.NewLookup(ctx, *lookup_uri)

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range flag.Args() {

		results, err := lookup.Find(ctx, code)

		if err != nil {
			log.Fatal(err)
		}

		for _, r := range results {
			fmt.Println(r)
		}
	}
}
