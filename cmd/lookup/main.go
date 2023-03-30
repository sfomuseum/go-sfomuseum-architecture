package main

import (
	_ "github.com/sfomuseum/go-sfomuseum-architecture/galleries"
	_ "github.com/sfomuseum/go-sfomuseum-architecture/gates"
	_ "github.com/sfomuseum/go-sfomuseum-architecture/terminals"
)

import (
	"context"
	"flag"
	"fmt"
	"log"

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

		for _, a := range results {
			fmt.Println(a)
		}
	}
}
