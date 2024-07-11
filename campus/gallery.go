package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (g *Gallery) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	g_id := g.WhosOnFirstId
	fmt.Fprintf(wr, "%s (gallery) %d %s\n", strings.Repeat("\t", indent), g_id, name(ctx, r, g_id))

	return nil

}
