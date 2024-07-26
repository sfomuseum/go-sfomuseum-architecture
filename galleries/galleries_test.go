package galleries

import (
	"context"
	"log/slog"
	"testing"
)

type galleryTest struct {
	Id   int64
	Code string
	Date string
}

func TestFindCurrentGallery(t *testing.T) {

	tests := map[string]int64{
		// "2E": 1763594985, // Kadish Gallery
		"2E": 1914601395,
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

func TestFindGalleryForDate(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	tests := []*galleryTest{
		&galleryTest{Id: 1763594985, Code: "2E", Date: "2022"},
		&galleryTest{Id: 1914650743, Code: "1G", Date: "2024-07-23"},
		&galleryTest{Id: 1914600907, Code: "3F", Date: "2024-06-18"},
		&galleryTest{Id: 1360392589, Code: "F04", Date: "2002"},
		// &galleryTest{Id: 1159157059, Code: "F-04", Date: "2017-12~"},
	}

	ctx := context.Background()

	for _, gallery := range tests {

		g, err := FindGalleryForDate(ctx, gallery.Code, gallery.Date)

		if err != nil {
			t.Fatalf("Failed to find gallery %s for %s, %v", gallery.Code, gallery.Date, err)
		}

		if g.WhosOnFirstId != gallery.Id {
			t.Fatalf("Unexpected ID for gallery %s. Got %d but expected %d", gallery.Code, g.WhosOnFirstId, gallery.Id)
		}
	}

}
