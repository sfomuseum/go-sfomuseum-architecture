package gates

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-geojson/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"sync"
)

// CompileGatesData will generate a list of `Gate` struct to be used as the source data for an `SFOMuseumLookup` instance.
// The list of gate are compiled by iterating over one or more source. `iterator_uri` is a valid `whosonfirst/go-whosonfirst-iterate` URI
// and `iterator_sources` are one more (iterator) URIs to process.
func CompileGatesData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]Gate, error) {

	lookup := make([]Gate, 0)
	mu := new(sync.RWMutex)

	iter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		path, err := emitter.PathForContext(ctx)

		if err != nil {
			return fmt.Errorf("Failed to derive path from context, %w", err)
		}

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return fmt.Errorf("Failed to parse %s, %w", path, err)
		}

		if uri_args.IsAlternate {
			return nil
		}

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil {
			return fmt.Errorf("Failed load feature from %s, %w", path, err)
		}

		wof_id := whosonfirst.Id(f)
		name := whosonfirst.Name(f)

		a := Gate{
			WOFID: wof_id,
			Name:  name,
		}

		mu.Lock()
		lookup = append(lookup, a)
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
