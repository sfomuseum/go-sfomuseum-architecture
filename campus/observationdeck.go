package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (ob *ObservationDeck) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	ob_id := ob.WhosOnFirstId
	fmt.Fprintf(wr, "%s (observation deck) %d %s\n", strings.Repeat("\t", indent), ob_id, name(ctx, r, ob_id))

	for _, g := range ob.Galleries {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gallery as tree, %w", err)
		}
	}

	for _, p := range ob.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	return nil

}
