package galleries

import (
	"context"
	"testing"
)

func TestFindCurrentGallery(t *testing.T) {

	tests := map[string]int64{
		"2E": 1763594985, // Kadish Gallery
	}

	ctx := context.Background()

	for code, id := range tests {

		g, err := FindCurrentGallery(ctx, code)

		if err != nil {
			t.Fatalf("Failed to find current gallery for %s, %v", code, err)
		}

		if g.WhosOnFirstId != id {
			t.Fatalf("Unexpected ID for gallery %s. Got %d but expected %d", code, g.WhosOnFirstId, id)
		}
	}
}
