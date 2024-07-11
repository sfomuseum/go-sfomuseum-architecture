package campus

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
)

func (g *Gate) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	g_id := g.WhosOnFirstId
	fmt.Fprintf(wr, "%s (gate) %d %s\n", strings.Repeat("\t", indent), g_id, name(ctx, r, g_id))

	return nil

}

func FindGates(ctx context.Context, db *sql.DB, parent_id int64) ([]*Gate, error) {

	slog.Debug("Find gates", "parent", parent_id)

	gate_ids, err := findChildIDs(ctx, db, parent_id, "gate")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (gates) for %d, %w", parent_id, err)
	}

	gates := make([]*Gate, 0)

	for _, g_id := range gate_ids {

		g_body, err := loadFeatureWithDBAndChecks(ctx, db, g_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for date %d, %w", g_id, err)
		}

		if g_body == nil {
			continue
		}

		var sfoid string

		rsp := gjson.GetBytes(g_body, "properties.sfo:id")

		if rsp.Exists() {

			sfoid = rsp.String()

		} else {

			rsp := gjson.GetBytes(g_body, "properties.wof:name")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing wof:name for %d", g_id)
			}

			sfoid = rsp.String()
		}

		name_rsp := gjson.GetBytes(g_body, "properties.wof:name")
		inception_rsp := gjson.GetBytes(g_body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(g_body, "properties.edtf:cessation")

		slog.Debug("Add gate", "sfo id", sfoid, "parent_id", parent_id, "id", g_id, "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())

		g := &Gate{
			WhosOnFirstId: g_id,
			SFOId:         sfoid,
		}

		gates = append(gates, g)
	}

	return gates, nil
}
