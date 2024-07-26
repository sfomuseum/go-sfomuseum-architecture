package terminals

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

// CompileTerminalsData will generate a list of `Terminal` struct to be used as the source data for an `SFOMuseumLookup` instance.
// The list of terminal are compiled by iterating over one or more source. `iterator_uri` is a valid `whosonfirst/go-whosonfirst-iterate` URI
// and `iterator_sources` are one more (iterator) URIs to process.
func CompileTerminalsData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]*Terminal, error) {

	lookup := make([]*Terminal, 0)
	mu := new(sync.RWMutex)

	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		if strings.HasSuffix(path, "~") {
			return nil
		}

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return fmt.Errorf("Failed to parse %s, %w", path, err)
		}

		if uri_args.IsAlternate {
			return nil
		}

		body, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		wof_id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID for %s, %w", path, err)
		}

		wof_name, err := properties.Name(body)

		if err != nil {
			return fmt.Errorf("Failed to derive name for %s, %w", path, err)
		}

		fl, err := properties.IsCurrent(body)

		if err != nil {
			return fmt.Errorf("Failed to determine is current for %s, %v", path, err)
		}

		preferred_names := make([]string, 0)
		variant_names := make([]string, 0)

		names := properties.Names(body)

		for k, k_names := range names {

			if strings.HasSuffix(k, "_preferred") {

				for _, n := range k_names {
					preferred_names = append(preferred_names, n)
				}
			} else if strings.HasSuffix(k, "_variant") {

				for _, n := range k_names {
					variant_names = append(variant_names, n)
				}
			} else {
			}

		}

		inception := properties.Inception(body)
		cessation := properties.Cessation(body)

		g := &Terminal{
			WhosOnFirstId:  wof_id,
			Name:           wof_name,
			IsCurrent:      fl.Flag(),
			PreferredNames: preferred_names,
			VariantNames:   variant_names,
			Inception:      inception,
			Cessation:      cessation,
		}

		sfom_rsp := gjson.GetBytes(body, "properties.sfomuseum:terminal_id")

		if sfom_rsp.Exists() {
			g.SFOMuseumId = sfom_rsp.String()
		}

		mu.Lock()
		lookup = append(lookup, g)
		mu.Unlock()

		return nil
	}

	iter, err := iterator.NewIterator(ctx, iterator_uri, iter_cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to create iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to iterate sources, %w", err)
	}

	return lookup, nil
}
