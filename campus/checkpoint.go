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

func (cp *Checkpoint) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	cp_id := cp.WhosOnFirstId
	fmt.Fprintf(wr, "%s (checkpoint) %d %s\n", strings.Repeat("\t", indent), cp_id, name(ctx, r, cp_id))

	return nil

}

func FindCheckpoints(ctx context.Context, db *sql.DB, parent_id int64) ([]*Checkpoint, error) {

	slog.Debug("Find check points", "parent id", parent_id)

	checkpoint_ids, err := findChildIDs(ctx, db, parent_id, "checkpoint")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (checkpoints) for %d, %w", parent_id, err)
	}

	checkpoints := make([]*Checkpoint, 0)

	for _, cp_id := range checkpoint_ids {

		cp_body, err := loadFeatureWithDBAndChecks(ctx, db, cp_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for %d, %w", cp_id, err)
		}

		if cp_body == nil {
			continue
		}

		var sfoid string

		rsp := gjson.GetBytes(cp_body, "properties.sfo:id")

		if !rsp.Exists() {
			return nil, fmt.Errorf("Missing sfo:id for %d", cp_id)
		}

		sfoid = rsp.String()

		name_rsp := gjson.GetBytes(cp_body, "properties.wof:name")
		inception_rsp := gjson.GetBytes(cp_body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(cp_body, "properties.edtf:cessation")

		slog.Debug("Add checkpoint", "sfo id", sfoid, "parent id", parent_id, "id", cp_id, "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())

		cp := &Checkpoint{
			WhosOnFirstId: cp_id,
			SFOId:         sfoid,
		}

		checkpoints = append(checkpoints, cp)
	}

	return checkpoints, nil
}
