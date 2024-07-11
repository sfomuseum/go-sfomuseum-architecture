package campus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
)

func (c *Complex) AsJSON(ctx context.Context, wr io.Writer) error {

	enc := json.NewEncoder(wr)
	return enc.Encode(c)
}

func (c *Complex) AsTree(ctx context.Context, r reader.Reader, wr io.Writer) error {

	name := func(id int64) string {

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

	for _, t := range c.Terminals {

		t_id := t.WhosOnFirstId
		fmt.Fprintf(wr, "%d %s\n", t_id, name(t_id))

		for _, b := range t.BoardingAreas {

			b_id := b.WhosOnFirstId
			fmt.Fprintf(wr, "\t%d %s\n", b_id, name(b_id))

			for _, g := range b.Galleries {
				g_id := g.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", g_id, name(g_id))
			}

			for _, p := range b.PublicArt {
				p_id := p.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", p_id, name(p_id))
			}

			for _, o := range b.ObservationDecks {
				o_id := o.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", o_id, name(o_id))

				for _, g := range o.Galleries {
					g_id := g.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", g_id, name(g_id))
				}

				for _, p := range o.PublicArt {
					p_id := p.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", p_id, name(p_id))
				}
			}

			for _, m := range b.Museums {
				m_id := m.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", m_id, name(m_id))

				for _, g := range m.Galleries {
					g_id := g.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", g_id, name(g_id))
				}

				for _, p := range m.PublicArt {
					p_id := p.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", p_id, name(p_id))
				}
			}

		}

		for _, c := range t.CommonAreas {

			c_id := c.WhosOnFirstId
			fmt.Fprintf(wr, "\t%d %s\n", c_id, name(c_id))

			for _, g := range c.Galleries {
				g_id := g.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", g_id, name(g_id))
			}

			for _, p := range c.PublicArt {
				p_id := p.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", p_id, name(p_id))
			}

			for _, o := range c.ObservationDecks {
				o_id := o.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", o_id, name(o_id))

				for _, g := range o.Galleries {
					g_id := g.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", g_id, name(g_id))
				}

				for _, p := range o.PublicArt {
					p_id := p.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", p_id, name(p_id))
				}
			}

			for _, m := range c.Museums {
				m_id := m.WhosOnFirstId
				fmt.Fprintf(wr, "\t\t%d %s\n", m_id, name(m_id))

				for _, g := range m.Galleries {
					g_id := g.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", g_id, name(g_id))
				}

				for _, p := range m.PublicArt {
					p_id := p.WhosOnFirstId
					fmt.Fprintf(wr, "\t\t\t%d %s\n", p_id, name(p_id))
				}
			}

		}

	}

	return nil
}
