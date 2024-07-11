package campus

// type Hotel is a lightweight data structure to represent garages at SFO with pointers its descendants.
type Hotel struct {
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	PublicArt     []*PublicArt `json:"publicart,omitempty"`
}
