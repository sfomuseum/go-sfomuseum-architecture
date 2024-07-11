package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

// type Hotel is a lightweight data structure to represent garages at SFO with pointers its descendants.
type Hotel struct {
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	PublicArt     []*PublicArt `json:"publicart,omitempty"`
}

func (h *Hotel) Id() int64 {
	return h.WhosOnFirstId
}

func (h *Hotel) Placetype() string {
	return "hotel"
}

func (h *Hotel) Walk(ctx context.Context, cb ElementCallbackFunc) error {

	for _, pa := range h.PublicArt {

		err := walkElement(ctx, pa, cb)

		if err != nil {
			return nil
		}
	}

	return nil
}

func (h *Hotel) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	h_id := h.WhosOnFirstId
	fmt.Fprintf(wr, "%s (hotel) %d %s\n", strings.Repeat("\t", indent), h_id, name(ctx, r, h_id))

	for _, p := range h.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	return nil

}
