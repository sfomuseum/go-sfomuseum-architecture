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

// type Museum is a lightweight data structure to represent dedicated Museum-related areas, distinct from galleries, at SFO  with pointers to its descendants.
type Museum struct {
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	Galleries     []*Gallery   `json:"galleries,omitempty"`
	PublicArt     []*PublicArt `json:"publicart,omitempty"`
}

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

func DeriveMuseums(ctx context.Context, db *sql.DB, parent_id int64) ([]*Museum, error) {

	slog.Debug("Derive museums", "parent id", parent_id)

	museum_ids, err := findChildIDs(ctx, db, parent_id, "museum")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (museums) for %d, %v", parent_id, err)
	}

	museums := make([]*Museum, 0)

	for _, m_id := range museum_ids {

		galleries, err := DeriveGalleries(ctx, db, m_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive galleries for museum %d, %w", m_id, err)
		}

		publicart, err := DerivePublicArt(ctx, db, m_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive public art for museum %d, %w", m_id, err)
		}

		m_body, err := loadFeatureWithDBAndChecks(ctx, db, m_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for museum %d, %w", m_id, err)
		}

		if m_body == nil {
			continue
		}

		var sfoid string

		rsp := gjson.GetBytes(m_body, "properties.sfo:id")

		if !rsp.Exists() {
			return nil, fmt.Errorf("Unable to find sfo:id for WOF record, %d", m_id)
		}

		sfoid = rsp.String()

		name_rsp := gjson.GetBytes(m_body, "properties.wof:name")
		inception_rsp := gjson.GetBytes(m_body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(m_body, "properties.edtf:cessation")

		slog.Debug("Add museum", "sfo id", sfoid, "parent id", parent_id, "id", m_id, "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())

		museum := &Museum{
			WhosOnFirstId: m_id,
			SFOId:         sfoid,
		}

		if len(galleries) > 0 {
			museum.Galleries = galleries
		}

		if len(publicart) > 0 {
			museum.PublicArt = publicart
		}

		museums = append(museums, museum)
	}

	return museums, nil
}
