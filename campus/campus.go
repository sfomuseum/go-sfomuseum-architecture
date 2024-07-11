// package campus provides methods for working with the SFO airport campus.
package campus

import (
	"context"
)

// type Campus is a lightweight data structure to represent the SFO campus with pointers its descendants.
type Campus struct {
	Element
	WhosOnFirstId int64        `json:"id"`
	SFOId         string       `json:"sfo:id"`
	Complex       *Complex     `json:"complex"`
	Garages       []*Garage    `json:"garages"`
	Hotels        []*Hotel     `json:"hotels"`
	PublicArt     []*PublicArt `json:"buildings,omitempty"`
}

func (c *Campus) Id() int64 {
	return c.WhosOnFirstId
}

func (c *Campus) Placetype() string {
	return "campus"
}

func (c *Campus) Walk(ctx context.Context, cb ElementCallbackFunc) error {

	err := walkElement(ctx, c.Complex, cb)

	if err != nil {
		return err
	}

	for _, g := range c.Garages {

		err := walkElement(ctx, g, cb)

		if err != nil {
			return err
		}
	}

	for _, h := range c.Hotels {

		err := walkElement(ctx, h, cb)

		if err != nil {
			return err
		}
	}

	for _, pa := range c.PublicArt {

		err := walkElement(ctx, pa, cb)

		if err != nil {
			return err
		}
	}

	return nil
}
