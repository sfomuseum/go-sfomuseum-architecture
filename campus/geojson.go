package campus

import (
	"context"

	"github.com/whosonfirst/go-reader"
)

func ToGeoJSON(ctx context.Context, c *Complex, r reader.Reader) error {

	terminal_ids := make([]int64, 0)
	boardingarea_ids := make([]int64, 0)
	commonarea_ids := make([]int64, 0)
	observationdeck_ids := make([]int64, 0)
	gallery_ids := make([]int64, 0)
	checkpoint_ids := make([]int64, 0)
	gate_ids := make([]int64, 0)
	publicart_ids := make([]int64, 0)

	for _, t := range c.Terminals {

		terminal_ids = append(terminal_ids, t.WhosOnFirstId)

		for _, b := range t.BoardingAreas {
			boardingarea_ids = append(boardingarea_ids, b.WhosOnFirstId)

			for _, g := range b.Gates {
				gate_ids = append(gate_ids, g.WhosOnFirstId)
			}

			for _, c := range b.Checkpoints {
				checkpoint_ids = append(checkpoint_ids, c.WhosOnFirstId)
			}

			for _, g := range b.Galleries {
				gallery_ids = append(gallery_ids, g.WhosOnFirstId)
			}

			for _, p := range b.PublicArt {
				publicart_ids = append(publicart_ids, p.WhosOnFirstId)
			}

			for _, o := range b.ObservationDecks {

				observationdeck_ids = append(observationdeck_ids, o.WhosOnFirstId)

				for _, g := range o.Galleries {
					gallery_ids = append(gallery_ids, g.WhosOnFirstId)
				}

				for _, p := range o.PublicArt {
					publicart_ids = append(publicart_ids, p.WhosOnFirstId)
				}

			}
		}

		for _, c := range t.CommonAreas {

			commonarea_ids = append(commonarea_ids, c.WhosOnFirstId)

			for _, g := range c.Gates {
				gate_ids = append(gate_ids, g.WhosOnFirstId)
			}

			for _, c := range c.Checkpoints {
				checkpoint_ids = append(checkpoint_ids, c.WhosOnFirstId)
			}

			for _, g := range c.Galleries {
				gallery_ids = append(gallery_ids, g.WhosOnFirstId)
			}

			for _, p := range c.PublicArt {
				publicart_ids = append(publicart_ids, p.WhosOnFirstId)
			}

			for _, o := range c.ObservationDecks {

				observationdeck_ids = append(observationdeck_ids, o.WhosOnFirstId)

				for _, g := range o.Galleries {
					gallery_ids = append(gallery_ids, g.WhosOnFirstId)
				}

				for _, p := range o.PublicArt {
					publicart_ids = append(publicart_ids, p.WhosOnFirstId)
				}
			}
		}
	}

	return nil
}
