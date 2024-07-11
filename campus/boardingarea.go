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

func (b *BoardingArea) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	b_id := b.WhosOnFirstId
	fmt.Fprintf(wr, "%s (boardingarea) %d %s\n", strings.Repeat("\t", indent), b_id, name(ctx, r, b_id))

	for _, g := range b.Gates {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gate as tree, %w", err)
		}
	}

	for _, c := range b.Checkpoints {

		err := c.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode checkpoint as tree, %w", err)
		}
	}

	for _, g := range b.Galleries {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gallery as tree, %w", err)
		}
	}

	for _, p := range b.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	for _, o := range b.ObservationDecks {

		err := o.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode observation deck as tree, %w", err)
		}
	}

	for _, m := range b.Museums {

		err := m.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode museum as tree, %w", err)
		}
	}

	return nil

}

func DeriveBoardingAreas(ctx context.Context, db *sql.DB, id int64) ([]*BoardingArea, error) {

	slog.Debug("Derive boarding areas", "parent", id)

	boardingarea_ids, err := findChildIDs(ctx, db, id, "boardingarea")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (boarding areas areas) for %d, %v", id, err)
	}

	boardingareas := make([]*BoardingArea, 0)

	for _, b_id := range boardingarea_ids {

		gates, err := DeriveGates(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for boarding area %d, %w", b_id, err)
		}

		checkpoints, err := DeriveCheckpoints(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive check points for boarding area %d, %w", b_id, err)
		}

		galleries, err := DeriveGalleries(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive galleries for boarding area %d, %w", b_id, err)
		}

		publicart, err := DerivePublicArt(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive public art for boarding area %d, %w", b_id, err)
		}

		observation_decks, err := DeriveObservationDecks(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive observation decks for boarding area %d, %w", b_id, err)
		}

		museums, err := DeriveMuseums(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive museums for boarding area %d, %w", b_id, err)
		}

		b_body, err := loadFeatureWithDBAndChecks(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for %d, %w", b_id, err)
		}

		if b_body == nil {
			continue
		}

		var sfoid string

		rsp := gjson.GetBytes(b_body, "properties.sfo:id")

		if rsp.Exists() {

			sfoid = rsp.String()

		} else {

			rsp := gjson.GetBytes(b_body, "properties.sfo:building_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing sfo:building_id for boarding area %d", b_id)
			}

			sfoid = rsp.String()
		}

		name_rsp := gjson.GetBytes(b_body, "properties.wof:name")
		inception_rsp := gjson.GetBytes(b_body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(b_body, "properties.edtf:cessation")

		slog.Debug("Add boardinarea", "sfo id", sfoid, "id", b_id, "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())

		area := &BoardingArea{
			WhosOnFirstId: b_id,
			SFOId:         sfoid,
		}

		if len(gates) > 0 {
			area.Gates = gates
		}

		if len(checkpoints) > 0 {
			area.Checkpoints = checkpoints
		}

		if len(galleries) > 0 {
			area.Galleries = galleries
		}

		if len(publicart) > 0 {
			area.PublicArt = publicart
		}

		if len(observation_decks) > 0 {
			area.ObservationDecks = observation_decks
		}

		if len(museums) > 0 {
			area.Museums = museums
		}

		boardingareas = append(boardingareas, area)
	}

	return boardingareas, nil

}
