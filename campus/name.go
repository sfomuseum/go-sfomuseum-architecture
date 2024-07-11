package campus

import (
	"context"
	"log/slog"

	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
)

func name(ctx context.Context, r reader.Reader, id int64) string {

	body, err := wof_reader.LoadBytes(ctx, r, id)

	if err != nil {
		slog.Warn("Failed to read bytes for ID", "id", id, "error", err)
		return ""
	}

	name, err := properties.Name(body)

	if err != nil {
		slog.Warn("Failed to read name", "id", id, "error", err)
		return ""
	}

	return name
}
