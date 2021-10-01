package galleries

import (
	"context"
	"testing"
)

func TestFindCurrentGallery(t *testing.T) {

	tests := map[string]int64{
		"D16": 1745882461,
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
