package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (b *BoardingArea) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	b_id := b.WhosOnFirstId
	fmt.Fprintf(wr, "%s (boardingarea) %d %s\n", strings.Repeat("\t", indent), b_id, name(ctx, r, b_id))

	for _, g := range b.Gates {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gate as tree, %w", err)
		}
	}

	for _, c := range b.Checkpoints {

		err := c.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode checkpoint as tree, %w", err)
		}
	}

	for _, g := range b.Galleries {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gallery as tree, %w", err)
		}
	}

	for _, p := range b.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	for _, o := range b.ObservationDecks {

		err := o.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode observation deck as tree, %w", err)
		}
	}

	for _, m := range b.Museums {

		err := m.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode museum as tree, %w", err)
		}
	}

	return nil

}
