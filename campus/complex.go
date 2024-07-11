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

func DeriveComplex(ctx context.Context, db *sql.DB, complex_id int64) (*Complex, error) {

	if complex_id == 0 {

		id, err := findMostRecentComplexID(ctx, db, FIRST_SFO)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive most recent complex ID, %w", err)
		}

		complex_id = id
	}

	terminals, err := DeriveTerminals(ctx, db, complex_id)

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

func findMostRecentComplexID(ctx context.Context, db *sql.DB, id int64) (int64, error) {

	q := "SELECT DISTINCT(superseded_by_id) FROM supersedes WHERE superseded_id = ?"

	rows, err := db.QueryContext(ctx, q, id)

	if err != nil {

		if err == sql.ErrNoRows {
			return id, nil
		}

		return -1, err
	}

	defer rows.Close()
	possible := make([]int64, 0)

	for rows.Next() {

		var superseded_by int64
		err := rows.Scan(&superseded_by)

		if err != nil {
			return -1, fmt.Errorf("Failed to scan row, %w", err)
		}

		possible = append(possible, superseded_by)
	}

	err = rows.Close()

	if err != nil {
		return -1, err
	}

	err = rows.Err()

	if err != nil {
		return -1, err
	}

	switch len(possible) {
	case 0:
		return id, nil
	case 1:
		return findMostRecentComplexID(ctx, db, possible[0])
	default:
		return -1, fmt.Errorf("Multiple results for '%d', not implemented", id)
	}
}
