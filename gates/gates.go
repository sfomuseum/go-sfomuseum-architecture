// package gates provides methods for working with boarding gates at SFO.
package gates

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-architecture"
)

// type Gate is a struct representing a passenger gate at SFO.
type Gate struct {
	// The Who's On First ID associated with this gate.
	WhosOnFirstId int64 `json:"wof:id"`
	// The name of this gate.
	Name string `json:"wof:name"`
	// A Who's On First "existential" (`KnownUnknownFlag`) flag signaling the gate's status
	IsCurrent int64 `json:"mz:is_current"`
}

// String() will return the name of the gate.
func (g *Gate) String() string {
	return fmt.Sprintf("%d %s (%d)", g.WhosOnFirstId, g.Name, g.IsCurrent)
}

// Return the current Gate matching 'code'. Multiple matches throw an error.
func FindCurrentGate(ctx context.Context, code string) (*Gate, error) {

	lookup, err := NewLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindCurrentGateWithLookup(ctx, lookup, code)
}

// Return the current Gate matching 'code' with a custom architecture.Lookup instance. Multiple matches throw an error.
func FindCurrentGateWithLookup(ctx context.Context, lookup architecture.Lookup, code string) (*Gate, error) {

	current, err := FindGatesCurrentWithLookup(ctx, lookup, code)

	if err != nil {
		return nil, err
	}

	switch len(current) {
	case 0:
		return nil, NotFound{code}
	case 1:
		return current[0], nil
	default:
		return nil, MultipleCandidates{code}
	}

}

// Returns all Gate instances matching 'code' that are marked as current.
func FindGatesCurrent(ctx context.Context, code string) ([]*Gate, error) {

	lookup, err := NewLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindGatesCurrentWithLookup(ctx, lookup, code)
}

// Returns all Gate instances matching 'code' that are marked as current with a custom architecture.Lookup instance.
func FindGatesCurrentWithLookup(ctx context.Context, lookup architecture.Lookup, code string) ([]*Gate, error) {

	rsp, err := lookup.Find(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("Failed to find gate '%s', %w", code, err)
	}

	current := make([]*Gate, 0)

	for _, r := range rsp {

		g := r.(*Gate)

		// if g.IsCurrent == 0 {
		if g.IsCurrent != 1 {
			continue
		}

		current = append(current, g)
	}

	return current, nil
}
