package testsdk

import (
	"context"
)

type Client struct{}

func (c *Client) GetSimple(ctx context.Context, id string) (*Simple, error) {
	panic("not implemented")
}

func (c *Client) CreateSimple(ctx context.Context, s *Simple) (*Simple, error) {
	panic("not implemented")
}

func (c *Client) UpdateSimple(ctx context.Context, s *Simple) (*Simple, error) {
	panic("not implemented")
}

func (c *Client) DeleteSimple(ctx context.Context, id string) error {
	panic("not implemented")
}

type AliasedInt64 int64

type Simple struct {
	String  string
	Int     int
	Int64   AliasedInt64
	Bool    bool
	Uint32  uint32
	Float32 float32
}
