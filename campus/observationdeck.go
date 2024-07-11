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

func (ob *ObservationDeck) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	ob_id := ob.WhosOnFirstId
	fmt.Fprintf(wr, "%s (observation deck) %d %s\n", strings.Repeat("\t", indent), ob_id, name(ctx, r, ob_id))

	for _, g := range ob.Galleries {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gallery as tree, %w", err)
		}
	}

	for _, p := range ob.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	return nil

}

func DeriveObservationDecks(ctx context.Context, db *sql.DB, t_id int64) ([]*ObservationDeck, error) {

	slog.Debug("Derive observation decks", "parent id", t_id)

	deck_ids, err := findChildIDs(ctx, db, t_id, "observationdeck")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (observation decks) for %d, %v", t_id, err)
	}

	decks := make([]*ObservationDeck, 0)

	for _, d_id := range deck_ids {

		galleries, err := DeriveGalleries(ctx, db, d_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive galleries for observation deck %d, %w", d_id, err)
		}

		publicart, err := DerivePublicArt(ctx, db, d_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive public art for observation deck %d, %w", d_id, err)
		}

		d_body, err := loadFeatureWithDBAndChecks(ctx, db, d_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for observation deck %d, %w", d_id, err)
		}

		if d_body == nil {
			continue
		}

		var sfoid string

		rsp := gjson.GetBytes(d_body, "properties.sfo:id")

		if !rsp.Exists() {
			return nil, fmt.Errorf("Unable to find sfo:id for WOF record, %d", d_id)
		}

		sfoid = rsp.String()

		name_rsp := gjson.GetBytes(d_body, "properties.wof:name")
		inception_rsp := gjson.GetBytes(d_body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(d_body, "properties.edtf:cessation")

		slog.Debug("Add observation deck", "sfo id", sfoid, "parent id", t_id, "id", d_id, "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())

		deck := &ObservationDeck{
			WhosOnFirstId: d_id,
			SFOId:         sfoid,
		}

		if len(galleries) > 0 {
			deck.Galleries = galleries
		}

		if len(publicart) > 0 {
			deck.PublicArt = publicart
		}

		decks = append(decks, deck)
	}

	return decks, nil
}
