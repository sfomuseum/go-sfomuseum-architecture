package terminals

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sfomuseum/go-sfomuseum-architecture"
	"github.com/sfomuseum/go-sfomuseum-architecture/data"
)

const DATA_JSON string = "terminals.json"

var lookup_table *sync.Map
var lookup_idx int64

var lookup_init sync.Once
var lookup_init_err error

type TerminalsLookupFunc func(context.Context)

type TerminalsLookup struct {
	architecture.Lookup
}

func init() {
	ctx := context.Background()
	architecture.RegisterLookup(ctx, "terminals", NewLookup)

	lookup_idx = int64(0)
}

// NewLookup will return an `architecture.Lookup` instance. By default the lookup table is derived from precompiled (embedded) data in `data/terminals.json`
// by passing in `sfomuseum://` as the URI. It is also possible to create a new lookup table with the following URI options:
//
//	`sfomuseum://github`
//
// This will cause the lookup table to be derived from the data stored at https://raw.githubusercontent.com/sfomuseum/go-sfomuseum-architecture/main/data/terminals.json. This might be desirable if there have been updates to the underlying data that are not reflected in the locally installed package's pre-compiled data.
//
//	`sfomuseum://iterator?uri={URI}&source={SOURCE}`
//
// This will cause the lookup table to be derived, at runtime, from data emitted by a `whosonfirst/go-whosonfirst-iterate` instance. `{URI}` should be a valid `whosonfirst/go-whosonfirst-iterate/iterator` URI and `{SOURCE}` is one or more URIs for the iterator to process.
func NewLookup(ctx context.Context, uri string) (architecture.Lookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// Reminder: u.Scheme is used by the architecture.Lookup constructor

	switch u.Host {
	case "iterator":

		q := u.Query()

		iterator_uri := q.Get("uri")
		iterator_sources := q["source"]

		return NewLookupFromIterator(ctx, iterator_uri, iterator_sources...)

	case "github":

		data_url := fmt.Sprintf("https://raw.githubusercontent.com/sfomuseum/go-sfomuseum-architecture/main/data/%s", DATA_JSON)
		rsp, err := http.Get(data_url)

		if err != nil {
			return nil, fmt.Errorf("Failed to load remote data from Github, %w", err)
		}

		lookup_func := NewLookupFuncWithReader(ctx, rsp.Body)
		return NewLookupWithLookupFunc(ctx, lookup_func)

	default:

		fs := data.FS
		fh, err := fs.Open(DATA_JSON)

		if err != nil {
			return nil, fmt.Errorf("Failed to load local precompiled data, %w", err)
		}

		lookup_func := NewLookupFuncWithReader(ctx, fh)
		return NewLookupWithLookupFunc(ctx, lookup_func)
	}
}

// NewLookupWithReader will return an `TerminalsLookupFunc` function instance that, when invoked, will populate an `architecture.Lookup` instance with data stored in `r`.
// `r` will be closed when the `TerminalsLookupFunc` function instance is invoked.
// It is assumed that the data in `r` will be formatted in the same way as the procompiled (embedded) data stored in `data/sfomuseum.json`.
func NewLookupFuncWithReader(ctx context.Context, r io.ReadCloser) TerminalsLookupFunc {

	defer r.Close()

	var terminals_list []*Terminal

	dec := json.NewDecoder(r)
	err := dec.Decode(&terminals_list)

	if err != nil {

		lookup_func := func(ctx context.Context) {
			lookup_init_err = err
		}

		return lookup_func
	}

	return NewLookupFuncWithTerminals(ctx, terminals_list)
}

// NewLookupFuncWithTerminals will return an `TerminalsLookupFunc` function instance that, when invoked, will populate an `architecture.Lookup` instance with data stored in `terminals_list`.
func NewLookupFuncWithTerminals(ctx context.Context, terminals_list []*Terminal) TerminalsLookupFunc {

	lookup_func := func(ctx context.Context) {

		table := new(sync.Map)

		for _, data := range terminals_list {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			appendData(ctx, table, data)
		}

		lookup_table = table
	}

	return lookup_func
}

// NewLookupWithLookupFunc will return an `architecture.Lookup` instance derived by data compiled using `lookup_func`.
func NewLookupWithLookupFunc(ctx context.Context, lookup_func TerminalsLookupFunc) (architecture.Lookup, error) {

	fn := func() {
		lookup_func(ctx)
	}

	lookup_init.Do(fn)

	if lookup_init_err != nil {
		return nil, lookup_init_err
	}

	l := TerminalsLookup{}
	return &l, nil
}

func NewLookupFromIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) (architecture.Lookup, error) {

	terminals_list, err := CompileTerminalsData(ctx, iterator_uri, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile terminals data, %w", err)
	}

	lookup_func := NewLookupFuncWithTerminals(ctx, terminals_list)
	return NewLookupWithLookupFunc(ctx, lookup_func)
}

func (l *TerminalsLookup) Find(ctx context.Context, code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, fmt.Errorf("Code '%s' not found", code)
	}

	terminals_list := make([]interface{}, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		terminals_list = append(terminals_list, row.(*Terminal))
	}

	return terminals_list, nil
}

func (l *TerminalsLookup) Append(ctx context.Context, data interface{}) error {
	return appendData(ctx, lookup_table, data.(*Terminal))
}

func appendData(ctx context.Context, table *sync.Map, data *Terminal) error {

	idx := atomic.AddInt64(&lookup_idx, 1)

	pointer := fmt.Sprintf("pointer:%d", idx)
	table.Store(pointer, data)

	str_wofid := strconv.FormatInt(data.WhosOnFirstId, 10)

	possible_codes := []string{
		data.Name,
		str_wofid,
	}

	for _, n := range data.PreferredNames {
		possible_codes = append(possible_codes, n)
	}

	for _, n := range data.VariantNames {
		possible_codes = append(possible_codes, n)
	}

	if data.SFOMuseumId != "" {
		possible_codes = append(possible_codes, data.SFOMuseumId)
	}

	for _, code := range possible_codes {

		if code == "" {
			continue
		}

		pointers := make([]string, 0)
		has_pointer := false

		others, ok := table.Load(code)

		if ok {

			pointers = others.([]string)
		}

		for _, dupe := range pointers {

			if dupe == pointer {
				has_pointer = true
				break
			}
		}

		if has_pointer {
			continue
		}

		pointers = append(pointers, pointer)
		table.Store(code, pointers)
	}

	return nil
}
