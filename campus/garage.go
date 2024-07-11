package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

// type Garage is a lightweight data structure to represent garages at SFO with pointers its descendants.
type Garage struct {
	Element
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	PublicArt     []*PublicArt `json:"publicart,omitempty"`
}

func (g *Garage) Id() int64 {
	return g.WhosOnFirstId
}

func (g *Garage) Placetype() string {
	return "garage"
}

func (g *Garage) Walk(ctx context.Context, cb ElementCallbackFunc) error {

	for _, pa := range g.PublicArt {

		err := walkElement(ctx, pa, cb)

		if err != nil {
			return nil
		}
	}

	return nil
}

func (g *Garage) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	g_id := g.WhosOnFirstId
	fmt.Fprintf(wr, "%s (garage) %d %s\n", strings.Repeat("\t", indent), g_id, name(ctx, r, g_id))

	for _, p := range g.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	return nil

}
