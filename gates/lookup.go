package gates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-architecture"
	"github.com/sfomuseum/go-sfomuseum-architecture/data"
	"io"
	_ "log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var lookup_table *sync.Map
var lookup_idx int64

var lookup_init sync.Once
var lookup_init_err error

type GatesLookupFunc func(context.Context)

type GatesLookup struct {
	architecture.Lookup
}

func init() {
	ctx := context.Background()
	architecture.RegisterLookup(ctx, "gates", NewLookup)

	lookup_idx = int64(0)
}

// NewLookup will return an `architecture.Lookup` instance derived from precompiled (embedded) data in `data/gates.json`
func NewLookup(ctx context.Context, uri string) (architecture.Lookup, error) {

	fs := data.FS
	fh, err := fs.Open("gates.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to load data, %v", err)
	}

	lookup_func := NewLookupFuncWithReader(ctx, fh)
	return NewLookupWithLookupFunc(ctx, lookup_func)
}

// NewLookupWithReader will return an `GatesLookupFunc` function instance that, when invoked, will populate an `architecture.Lookup` instance with data stored in `r`.
// `r` will be closed when the `GatesLookupFunc` function instance is invoked.
// It is assumed that the data in `r` will be formatted in the same way as the procompiled (embedded) data stored in `data/sfomuseum.json`.
func NewLookupFuncWithReader(ctx context.Context, r io.ReadCloser) GatesLookupFunc {

	defer r.Close()

	var gates_list []*Gate

	dec := json.NewDecoder(r)
	err := dec.Decode(&gates_list)

	if err != nil {

		lookup_func := func(ctx context.Context) {
			lookup_init_err = err
		}

		return lookup_func
	}

	return NewLookupFuncWithGates(ctx, gates_list)
}

// NewLookupFuncWithGates will return an `GatesLookupFunc` function instance that, when invoked, will populate an `architecture.Lookup` instance with data stored in `gates_list`.
func NewLookupFuncWithGates(ctx context.Context, gates_list []*Gate) GatesLookupFunc {

	lookup_func := func(ctx context.Context) {

		table := new(sync.Map)

		for _, data := range gates_list {

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
func NewLookupWithLookupFunc(ctx context.Context, lookup_func GatesLookupFunc) (architecture.Lookup, error) {

	fn := func() {
		lookup_func(ctx)
	}

	lookup_init.Do(fn)

	if lookup_init_err != nil {
		return nil, lookup_init_err
	}

	l := GatesLookup{}
	return &l, nil
}

func NewLookupFromIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) (architecture.Lookup, error) {

	gates_list, err := CompileGatesData(ctx, iterator_uri, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile gates data, %w", err)
	}

	lookup_func := NewLookupFuncWithGates(ctx, gates_list)
	return NewLookupWithLookupFunc(ctx, lookup_func)
}

func (l *GatesLookup) Find(ctx context.Context, code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, fmt.Errorf("Code '%s' not found", code)
	}

	gates_list := make([]interface{}, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		gates_list = append(gates_list, row.(*Gate))
	}

	return gates_list, nil
}

func (l *GatesLookup) Append(ctx context.Context, data interface{}) error {
	return appendData(ctx, lookup_table, data.(*Gate))
}

func appendData(ctx context.Context, table *sync.Map, data *Gate) error {

	idx := atomic.AddInt64(&lookup_idx, 1)

	pointer := fmt.Sprintf("pointer:%d", idx)
	table.Store(pointer, data)

	str_wofid := strconv.FormatInt(data.WhosOnFirstId, 10)

	possible_codes := []string{
		data.Name,
		str_wofid,
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
