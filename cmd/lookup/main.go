package main

import (
	_ "github.com/sfomuseum/go-sfomuseum-architecture/galleries"
	_ "github.com/sfomuseum/go-sfomuseum-architecture/gates"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-architecture"
	"log"
)

func main() {

	lookup_uri := flag.String("lookup-uri", "", "...")

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
