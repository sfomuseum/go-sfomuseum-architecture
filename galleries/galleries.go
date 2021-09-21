// package galleries provides methods for working with boarding galleries at SFO.
package galleries

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
}

// String() will return the name of the gallery.
func (g *Gallery) String() string {
	return g.Name
}
