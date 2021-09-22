// package galleries provides methods for working with boarding galleries at SFO.
package galleries

import (
	"fmt"
)

// type Gallery is a struct representing a passenger gallery at SFO.
type Gallery struct {
	// The Who's On First ID associated with this gallery.
	WhosOnFirstId int64 `json:"wof:id"`
	// The SFO Museum ID associated with this gallery.
	SFOMuseumId int64 `json:"sfomuseum:id"`
	// The map label (ID) associated with this gallery.
	MapId string `json:"map_id"`
	// The name of this gallery.
	Name string `json:"wof:name"`
	// The (EDTF) inception date for the gallery
	Inception string `json:"edtf:inception"`
	// The (EDTF) cessation date for the gallery
	Cessation string `json:"edtf:cessation"`
}

// String() will return the name of the gallery.
func (g *Gallery) String() string {
	return fmt.Sprintf("%d#%d %s-%s %s", g.WhosOnFirstId, g.SFOMuseumId, g.Inception, g.Cessation, g.Name)
}
