package gates

import (
	"context"
	"testing"
)

func TestFindCurrentGate(t *testing.T) {

	tests := map[string]int64{
		"A9": 1763588417,
	}

	ctx := context.Background()

	for code, id := range tests {

		g, err := FindCurrentGate(ctx, code)

		if err != nil {
			t.Fatalf("Failed to find current gate for %s, %v", code, err)
		}

		if g.WhosOnFirstId != id {
			t.Fatalf("Unexpected ID for gate %s. Got %d but expected %d", code, g.WhosOnFirstId, id)
		}
	}
}
