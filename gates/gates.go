// package gates provides methods for working with boarding gates at SFO.
package gates

import (
	"fmt"
)

// type Gate is a struct representing a passenger gate at SFO.
type Gate struct {
	// The Who's On First ID associated with this gate.
	WhosOnFirstId int64 `json:"wof:id"`
	// The name of this gate.
	Name string `json:"wof:name"`
	// A Who's On First "existential" (`KnownUnknownFlag`) flag signaling the gate's status
	IsCurrent string `json:"mz:is_current"`
}

// String() will return the name of the gate.
func (g *Gate) String() string {
	return fmt.Sprintf("%d %s (%s)", g.WhosOnFirstId, g.Name, g.IsCurrent)
}
