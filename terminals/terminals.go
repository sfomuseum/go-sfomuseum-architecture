// package terminals provides methods for working with boarding terminals at SFO.
package terminals

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-architecture"
)

// type Terminal is a struct representing a passenger terminal at SFO.
type Terminal struct {
	// The Who's On First ID associated with this terminal.
	WhosOnFirstId int64 `json:"wof:id"`
	// The SFO Museum name/label for this terminal
	SFOMuseumId string `json:"sfomuseum:terminal_id,omityempty"`
	// The name of this terminal.
	Name string `json:"wof:name"`
	// A Who's On First "existential" (`KnownUnknownFlag`) flag signaling the terminal's status
	IsCurrent int64 `json:"mz:is_current"`
	// The list of name:{LANG}_x_preferred names for this terminal
	PreferredNames []string `json:"name:preferred,omitempty"`
	// The list of name:{LANG}_x_variant names for this terminal
	VariantNames []string `json:"name:variant,omitempty"`
}

// String() will return the name of the terminal.
func (g *Terminal) String() string {
	return fmt.Sprintf("%d %s (%d)", g.WhosOnFirstId, g.Name, g.IsCurrent)
}

// Return the current Terminal matching 'code'. Multiple matches throw an error.
func FindCurrentTerminal(ctx context.Context, code string) (*Terminal, error) {

	lookup, err := NewLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindCurrentTerminalWithLookup(ctx, lookup, code)
}

// Return the current Terminal matching 'code' with a custom architecture.Lookup instance. Multiple matches throw an error.
func FindCurrentTerminalWithLookup(ctx context.Context, lookup architecture.Lookup, code string) (*Terminal, error) {

	current, err := FindTerminalsCurrentWithLookup(ctx, lookup, code)

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

// Returns all Terminal instances matching 'code' that are marked as current.
func FindTerminalsCurrent(ctx context.Context, code string) ([]*Terminal, error) {

	lookup, err := NewLookup(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new lookup, %w", err)
	}

	return FindTerminalsCurrentWithLookup(ctx, lookup, code)
}

// Returns all Terminal instances matching 'code' that are marked as current with a custom architecture.Lookup instance.
func FindTerminalsCurrentWithLookup(ctx context.Context, lookup architecture.Lookup, code string) ([]*Terminal, error) {

	rsp, err := lookup.Find(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("Failed to find terminal '%s', %w", code, err)
	}

	current := make([]*Terminal, 0)

	for _, r := range rsp {

		g := r.(*Terminal)

		// if g.IsCurrent == 0 {
		if g.IsCurrent != 1 {
			continue
		}

		current = append(current, g)
	}

	return current, nil
}
