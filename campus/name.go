package campus

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
)

func treeName(ctx context.Context, r reader.Reader, el Element, indent int) string {

	id := el.Id()
	pt := el.Placetype()

	return fmt.Sprintf("%s (%s) %d %s\n", strings.Repeat("\t", indent), pt, id, name(ctx, r, id))
}

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
