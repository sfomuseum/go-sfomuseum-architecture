package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/sfomuseum/go-sfomuseum-architecture/campus"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
)

func main() {

	var iterator_uri string
	var output_mode string
	var reader_uri string

	flag.StringVar(&iterator_uri, "iterator-uri", "repo://", "...")
	flag.StringVar(&output_mode, "output-mode", "json", "...")
	flag.StringVar(&reader_uri, "reader-uri", "repo:///usr/local/data/sfomuseum-data-architecture", "...")

	flag.Parse()

	ctx := context.Background()

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

		enc := json.NewEncoder(wr)
		err = enc.Encode(c)

		if err != nil {
			log.Fatalf("Failed to encode complex, %v", err)
		}

	case "tree":

		r, err := reader.NewReader(ctx, reader_uri)

		if err != nil {
			log.Fatalf("Failed to create new reader, %v", err)
		}

		name := func(id int64) string {

			body, err := wof_reader.LoadBytes(ctx, r, id)

			if err != nil {
				slog.Warn("Failed to read bytes for ID", "id", id, "error", err)
				return ""
			}

			name, err := properties.Name(body)

			if err != nil {
				slog.Warn("Failed to read name", "id", id, "error", err)
				return ""
			}

			return name
		}

		for _, t := range c.Terminals {

			t_id := t.WhosOnFirstId
			fmt.Fprintf(wr, "%d %s\n", t_id, name(t_id))

			for _, b := range t.BoardingAreas {

				b_id := b.WhosOnFirstId
				fmt.Fprintf(wr, "\t%d %s\n", b_id, name(b_id))

				for _, g := range b.Galleries {
					g_id := g.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t%d %s\n", g_id, name(g_id))
				}

				for _, p := range b.PublicArt {
					p_id := p.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t%d %s\n", p_id, name(p_id))
				}

				for _, o := range b.ObservationDecks {
					o_id := o.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", o_id, name(o_id))

					for _, g := range o.Galleries {
						g_id := g.WhosOnFirstId
						fmt.Fprintf(wr, "\t\t\t\t%d %s\n", g_id, name(g_id))
					}

					for _, p := range o.PublicArt {
						p_id := p.WhosOnFirstId
						fmt.Fprintf(wr, "\t\t\t\t%d %s\n", p_id, name(p_id))
					}
				}
			}

			for _, c := range t.CommonAreas {

				c_id := c.WhosOnFirstId
				fmt.Fprintf(wr, "\t%d %s\n", c_id, name(c_id))

				for _, g := range c.Galleries {
					g_id := g.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t%d %s\n", g_id, name(g_id))
				}

				for _, p := range c.PublicArt {
					p_id := p.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t%d %s\n", p_id, name(p_id))
				}

				for _, o := range c.ObservationDecks {
					o_id := o.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", o_id, name(o_id))

					for _, g := range o.Galleries {
						g_id := g.WhosOnFirstId
						fmt.Fprintf(wr, "\t\t\t\t%d %s\n", g_id, name(g_id))
					}

					for _, p := range o.PublicArt {
						p_id := p.WhosOnFirstId
						fmt.Fprintf(wr, "\t\t\t\t%d %s\n", p_id, name(p_id))
					}
				}

			}

		}

	default:
		log.Fatalf("Invalid or unsupported output mode, %v", err)
	}
}
