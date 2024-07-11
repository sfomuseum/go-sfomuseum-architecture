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

func (ca *CommonArea) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	ca_id := ca.WhosOnFirstId
	fmt.Fprintf(wr, "%s (commonarea) %d %s\n", strings.Repeat("\t", indent), ca_id, name(ctx, r, ca_id))

	for _, g := range ca.Gates {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gate as tree, %w", err)
		}
	}

	for _, c := range ca.Checkpoints {

		err := c.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode checkpoint as tree, %w", err)
		}
	}

	for _, g := range ca.Galleries {

		err := g.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode gallery as tree, %w", err)
		}
	}

	for _, p := range ca.PublicArt {

		err := p.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode public art as tree, %w", err)
		}
	}

	for _, o := range ca.ObservationDecks {

		err := o.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode observation deck as tree, %w", err)
		}
	}

	for _, m := range ca.Museums {

		err := m.AsTree(ctx, r, wr, indent+1)

		if err != nil {
			return fmt.Errorf("Failed to encode museum as tree, %w", err)
		}
	}

	return nil

}

func DeriveCommonAreas(ctx context.Context, db *sql.DB, parent_id int64) ([]*CommonArea, error) {

	slog.Debug("Derive common areas", "parent", parent_id)

	commonarea_ids, err := findChildIDs(ctx, db, parent_id, "commonarea")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (common areas) for %d, %v", parent_id, err)
	}

	commonareas := make([]*CommonArea, 0)

	for _, c_id := range commonarea_ids {

		gates, err := DeriveGates(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for common area %d, %w", c_id, err)
		}

		checkpoints, err := DeriveCheckpoints(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for check points %d, %w", c_id, err)
		}

		galleries, err := DeriveGalleries(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for galleries %d, %w", c_id, err)
		}

		observation_decks, err := DeriveObservationDecks(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive observation decks for galleries %d, %w", c_id, err)
		}

		museums, err := DeriveMuseums(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive museums for common area %d, %w", c_id, err)
		}

		publicart, err := DerivePublicArt(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive public art for common area %d, %w", c_id, err)
		}

		c_body, err := loadFeatureWithDBAndChecks(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature %d, %w", c_id, err)
		}

		if c_body == nil {
			continue
		}

		var sfoid string

		rsp := gjson.GetBytes(c_body, "properties.sfo:id")

		if rsp.Exists() {

			sfoid = rsp.String()

		} else {

			rsp := gjson.GetBytes(c_body, "properties.sfo:building_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Unable to find sfo:building_id for %d", c_id)
			}

			switch rsp.String() {
			case "ITB", "100":
				sfoid = "100CAD" // gis.COMMONAREA_ITB_DEPARTURES
			case "T1", "200":
				sfoid = "200CAD" // gis.COMMONAREA_T1_DEPARTURES
			case "T2", "300":
				sfoid = "300CAD" // gis.COMMONAREA_T2_DEPARTURES
			case "T3", "400":
				sfoid = "400CAD" // gis.COMMONAREA_T3_DEPARTURES
			default:
				return nil, fmt.Errorf("Unrecognized sfo:id '%s' for %d", rsp.String(), c_id)
			}
		}

		name_rsp := gjson.GetBytes(c_body, "properties.wof:name")
		inception_rsp := gjson.GetBytes(c_body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(c_body, "properties.edtf:cessation")

		slog.Debug("Add common area", "sfo id", sfoid, "id", c_id, "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())

		area := &CommonArea{
			WhosOnFirstId: c_id,
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

		commonareas = append(commonareas, area)
	}

	return commonareas, nil
}
