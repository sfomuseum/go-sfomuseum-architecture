// package campus provides methods for working with the SFO airport campus.
package campus

// type Campus is a lightweight data structure to represent the SFO campus with pointers its descendants.
type Campus struct {
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	Complex       *Complex     `json:"complex"`
	Garages       []*Garage    `json:"garages"`
	Hotels        []*Hotel     `json:"hotels"`
	PublicArt     []*PublicArt `json:"buildings,omitempty"`
}
