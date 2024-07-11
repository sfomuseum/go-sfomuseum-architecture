package campus

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/whosonfirst/go-reader"
)

func (m *Museum) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	m_id := m.WhosOnFirstId
	fmt.Fprintf(wr, "%s (museum) %d %s\n", strings.Repeat("\t", indent), m_id, name(ctx, r, m_id))

	for _, g := range m.Galleries {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gallery as tree, %w", err)
		}
	}

	for _, p := range m.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	return nil

}
