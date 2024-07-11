package campus

import (
	"context"
	"io"

	"github.com/whosonfirst/go-reader"
)

type Element interface {
	Id() int64
	Placetype() string
	Walk(context.Context, ElementCallbackFunc) error
	AsTree(context.Context, reader.Reader, io.Writer, int) error
}

type ElementCallbackFunc func(context.Context, Element) error

func walkElement(ctx context.Context, el Element, cb ElementCallbackFunc) error {

	err := cb(ctx, el)

	if err != nil {
		return err
	}

	err = el.Walk(ctx, cb)

	if err != nil {
		return err
	}

	return nil
}
