package cache

import (
	"context"
)

type inmemCache struct {
	data map[string]interface{}
}

func (c *inmemCache) Set(ctx context.Context, k string, v interface{}) error {
	c.data[k] = v
	return nil
}

func (c *inmemCache) Get(ctx context.Context, k string) (interface{}, error) {
	if v, ok := c.data[k]; ok {
		return v, nil
	}

	return nil, ErrCache.Code("not_found")
}

func (c *inmemCache) Delete(ctx context.Context, k string) error {
	delete(c.data, k)
	return nil
}
