package campus

import (
	"context"
	"fmt"

	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-reader"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
)

func ComplexToGeoJSONLayers(ctx context.Context, c *Complex, r reader.Reader) (map[string]*geojson.FeatureCollection, error) {

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

	feature_ids := map[string][]int64{
		"terminals":        terminal_ids,
		"boardingareas":    boardingarea_ids,
		"commonareas":      commonarea_ids,
		"observationdecks": observationdeck_ids,
		"gates":            gate_ids,
		"checkpoints":      checkpoint_ids,
		"galleries":        gallery_ids,
		"publicart":        publicart_ids,
	}

	features := make(map[string]*geojson.FeatureCollection)

	for k, ids := range feature_ids {

		fc := geojson.NewFeatureCollection()

		for _, id := range ids {

			body, err := wof_reader.LoadBytes(ctx, r, id)

			if err != nil {
				return nil, fmt.Errorf("Failed to read %d, %w", id, err)
			}

			f, err := geojson.UnmarshalFeature(body)

			if err != nil {
				return nil, fmt.Errorf("Failed to unmarshal %d, %w", id, err)
			}

			fc.Append(f)
		}

		features[k] = fc
	}

	return features, nil
}
