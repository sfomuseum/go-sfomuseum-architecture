// package campus provides methods for working with the SFO airport campus.
package campus

// type Campus is a lightweight data structure to represent the SFO campus with pointers its descendants.
type Campus struct {
	WhosOnFirstId int64     `json:"id"`
	SFOId         string    `json:"sfo:id"`
	Complex       *Complex  `json:"complex"`
	Garages       []*Garage `json:"garages"`
	// Buildings []*Building `json:"buildings,omitempty"`
	PublicArt []*PublicArt `json:"buildings,omitempty"`
}

// type Garage is a lightweight data structure to represent garages at SFO with pointers its descendants.
type Garage struct {
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	PublicArt     []*PublicArt `json:"publicart,omitempty"`
}

// type Complex is a lightweight data structure to represent the terminal complex at SFO with pointers its descendants.
type Complex struct {
	WhosOnFirstId int64       `json:"id"`
	SFOId         string      `json:"sfo:id"`
	Terminals     []*Terminal `json:"terminals"`
}

// type ObservationDeck is a lightweight data structure to represent observation decks at SFO with pointers its descendants.
type ObservationDeck struct {
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	PublicArt     []*PublicArt `json:"publicart,omitempty"`
	Galleries     []*Gallery   `json:"galleries,omitempty"`
}

// type Terminal is a lightweight data structure to represent terminals at SFO with pointers its descendants.
type Terminal struct {
	WhosOnFirstId int64           `json:"id"`
	SFOId         string          `json:"sfo:id"`
	CommonAreas   []*CommonArea   `json:"commonareas,omitempty"`
	BoardingAreas []*BoardingArea `json:"boardingareas,omitempty"`
}

// type CommonArea is a lightweight data structure to represent common areas at SFO with pointers its descendants.
type CommonArea struct {
	WhosOnFirstId    int64              `json:"id"`
	SFOId            string             `json:"sfo:id"`
	Gates            []*Gate            `json:"gates,omitempty"`
	Checkpoints      []*Checkpoint      `json:"checkpoints,omitempty"`
	Galleries        []*Gallery         `json:"galleries,omitempty"`
	PublicArt        []*PublicArt       `json:"publicart,omitempty"`
	ObservationDecks []*ObservationDeck `json:"observationdecks,omitempty"` // for example T2
}

// type BoardingArea is a lightweight data structure to represent boarding areas at SFO with pointers its descendants.
type BoardingArea struct {
	WhosOnFirstId    int64              `json:"id"`
	SFOId            string             `json:"sfo:id"`
	Gates            []*Gate            `json:"gates,omitempty"`
	Checkpoints      []*Checkpoint      `json:"checkpoints,omitempty"`
	Galleries        []*Gallery         `json:"galleries,omitempty"`
	PublicArt        []*PublicArt       `json:"publicart,omitempty"`
	ObservationDecks []*ObservationDeck `json:"observationdecks,omitempty"`
}

// type Gallery is a lightweight data structure to represent SFO Museum galleries at SFO.
type Gallery struct {
	WhosOnFirstId int64  `json:"id"`
	SFOId         string `json:"sfomuseum:id"`
}

// type Gate is a lightweight data structure to represent passenger gates at SFO.
type Gate struct {
	WhosOnFirstId int64  `json:"id"`
	SFOId         string `json:"sfo:id"`
}

// type Checkpoint is a lightweight data structure to represent security checkpoints at SFO.
type Checkpoint struct {
	WhosOnFirstId int64  `json:"id"`
	SFOId         string `json:"sfo:id"`
}

// type PublicArt is a lightweight data structure to represent public art works at SFO.
type PublicArt struct {
	WhosOnFirstId int64  `json:"id"`
	SFOId         string `json:"sfomuseum:id"`
}
