package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (cp *Checkpoint) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	cp_id := cp.WhosOnFirstId
	fmt.Fprintf(wr, "%s (checkpoint) %d %s\n", strings.Repeat("\t", indent), cp_id, name(ctx, r, cp_id))

	return nil

}
