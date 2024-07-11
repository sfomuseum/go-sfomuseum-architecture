package campus

// type Garage is a lightweight data structure to represent garages at SFO with pointers its descendants.
type Garage struct {
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	PublicArt     []*PublicArt `json:"publicart,omitempty"`
}
