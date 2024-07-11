package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (ca *CommonArea) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	ca_id := ca.WhosOnFirstId
	fmt.Fprintf(wr, "%s (commonarea) %d %s\n", strings.Repeat("\t", indent), ca_id, name(ctx, r, ca_id))

	for _, g := range ca.Gates {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gate as tree, %w", err)
		}
	}

	for _, c := range ca.Checkpoints {

		err := c.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode checkpoint as tree, %w", err)
		}
	}

	for _, g := range ca.Galleries {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gallery as tree, %w", err)
		}
	}

	for _, p := range ca.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	for _, o := range ca.ObservationDecks {

		err := o.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode observation deck as tree, %w", err)
		}
	}

	for _, m := range ca.Museums {

		err := m.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode museum as tree, %w", err)
		}
	}

	return nil

}
