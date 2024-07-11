package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (p *PublicArt) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	p_id := p.WhosOnFirstId
	fmt.Fprintf(wr, "%s (public art) %d %s\n", strings.Repeat("\t", indent), p_id, name(ctx, r, p_id))

	return nil

}
