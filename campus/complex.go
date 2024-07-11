package campus

import (
	"context"
	"database/sql"
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

func FindComplex(ctx context.Context, db *sql.DB, complex_id int64) (*Complex, error) {

	terminals, err := FindTerminals(ctx, db, complex_id)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive terminals for complex %d, %w", complex_id, err)
	}

	c := &Complex{
		WhosOnFirstId: complex_id,
		SFOId:         "SFO",
		Terminals:     terminals,
	}

	return c, nil
}
