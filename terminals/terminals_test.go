package terminals

import (
	"context"
	"testing"
)

func TestFindCurrentTerminal(t *testing.T) {

	tests := map[string]int64{
		"T2": 1763588123,
	}

	ctx := context.Background()

	for code, id := range tests {

		g, err := FindCurrentTerminal(ctx, code)

		if err != nil {
			t.Fatalf("Failed to find current terminal for %s, %v", code, err)
		}

		if g.WhosOnFirstId != id {
			t.Fatalf("Unexpected ID for terminal %s. Got %d but expected %d", code, g.WhosOnFirstId, id)
		}
	}
}
