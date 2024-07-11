package campus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (c *Complex) AsJSON(ctx context.Context, wr io.Writer) error {

	enc := json.NewEncoder(wr)
	return enc.Encode(c)
}

func (c *Complex) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	c_id := c.WhosOnFirstId
	fmt.Fprintf(wr, "%s (complex) %d %s\n", strings.Repeat("\t", indent), c_id, name(ctx, r, c_id))

	for _, t := range c.Terminals {

		err := t.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode Terminal as tree, %w", err)
		}
	}

	return nil
}
