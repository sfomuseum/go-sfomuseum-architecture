package galleries

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"sync"
)

// CompileGalleriesData will generate a list of `Gallery` struct to be used as the source data for an `SFOMuseumLookup` instance.
// The list of gate are compiled by iterating over one or more source. `iterator_uri` is a valid `whosonfirst/go-whosonfirst-iterate` URI
// and `iterator_sources` are one more (iterator) URIs to process.
func CompileGalleriesData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]*Gallery, error) {

	lookup := make([]*Gallery, 0)
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

		body, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Failed load feature from %s, %w", path, err)
		}

		name_rsp := gjson.GetBytes(body, "properties.wof:name")
		wofid_rsp := gjson.GetBytes(body, "properties.wof:id")
		sfomid_rsp := gjson.GetBytes(body, "properties.sfomuseum:gallery_id")

		if !name_rsp.Exists() {
			return fmt.Errorf("Missing wof:name property (%s)", path)
		}

		if !wofid_rsp.Exists() {
			return fmt.Errorf("Missing wof:id property (%s)", path)
		}

		if !sfomid_rsp.Exists() {
			return fmt.Errorf("Missing sfomuseum:gallery_id property (%s)", path)
		}

		mapid_rsp := gjson.GetBytes(body, "properties.sfomuseum:map_id")
		inception_rsp := gjson.GetBytes(body, "properties.edtf:inception")
		cessation_rsp := gjson.GetBytes(body, "properties.edtf:cessation")

		g := &Gallery{
			WhosOnFirstId: wofid_rsp.Int(),
			SFOMuseumId:   sfomid_rsp.Int(),
			MapId:         mapid_rsp.String(),
			Name:          name_rsp.String(),
			Inception:     inception_rsp.String(),
			Cessation:     cessation_rsp.String(),
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
