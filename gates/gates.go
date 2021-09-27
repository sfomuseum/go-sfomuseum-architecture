// package gates provides methods for working with boarding gates at SFO.
package gates

// type Gate is a struct representing a passenger gate at SFO.
type Gate struct {
	// The Who's On First ID associated with this gate.
	WhosOnFirstId int64 `json:"wof:id"`
	// The name of this gate.
	Name string `json:"wof:name"`
}

// String() will return the name of the gate.
func (g *Gate) String() string {
	return g.Name
}
