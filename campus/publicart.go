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

// type PublicArt is a lightweight data structure to represent public art works at SFO.
type PublicArt struct {
	Element
	WhosOnFirstId int64  `json:"id"`
	SFOId         string `json:"sfomuseum:id"`
}

func (c *PublicArt) Id() int64 {
	return c.WhosOnFirstId
}

func (c *PublicArt) Placetype() string {
	return "publicart"
}

func (c *PublicArt) Walk(ctx context.Context, cb ElementCallbackFunc) error {
	return nil
}

func (p *PublicArt) AsTree(ctx context.Context, r reader.Reader, wr io.Writer, indent int) error {

	p_id := p.WhosOnFirstId
	fmt.Fprintf(wr, "%s (public art) %d %s\n", strings.Repeat("\t", indent), p_id, name(ctx, r, p_id))

	return nil

}

func DerivePublicArt(ctx context.Context, db *sql.DB, parent_id int64) ([]*PublicArt, error) {

	slog.Debug("Derive public art", "parent id", parent_id)

	publicart_ids, err := findChildIDs(ctx, db, parent_id, "publicart")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (public art) for %d, %w", parent_id, err)
	}

	publicarts := make([]*PublicArt, 0)

	for _, p_id := range publicart_ids {

		p_body, err := loadFeatureWithDBAndChecks(ctx, db, p_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for public art %d, %w", p_id, err)
		}

		if p_body == nil {
			continue
		}

		var sfomid string

		rsp := gjson.GetBytes(p_body, "properties.sfomuseum:map_id")

		if rsp.Exists() {

			sfomid = rsp.String()

		} else {

			rsp := gjson.GetBytes(p_body, "properties.sfomuseum:object_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing sfomuseum:object_id property for public art %d, %w", p_id, err)
			}

			sfomid = rsp.String()
		}

		name_rsp := gjson.GetBytes(p_body, "properties.wof:name")
		inception_rsp := gjson.GetBytes(p_body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(p_body, "properties.edtf:cessation")

		slog.Debug("Add public art", "sfo id", sfomid, "parent id", parent_id, "id", p_id, "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())

		pa := &PublicArt{
			WhosOnFirstId: p_id,
			SFOId:         sfomid,
		}

		publicarts = append(publicarts, pa)
	}

	return publicarts, nil
}
