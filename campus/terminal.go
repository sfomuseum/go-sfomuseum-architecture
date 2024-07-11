package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (t *Terminal) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	t_id := t.WhosOnFirstId
	fmt.Fprintf(wr, "%s (terminal) %d %s\n", strings.Repeat("\t", indent), t_id, name(ctx, r, t_id))

	for _, b := range t.BoardingAreas {

		err := b.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode boarding area as tree, %w", err)
		}
	}

	for _, c := range t.CommonAreas {

		err := c.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode common area as tree, %w", err)
		}
	}

	return nil

}
