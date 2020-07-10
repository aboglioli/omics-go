package security

import (
	"context"
	"omics/pkg/errors"
)

type fakeCache struct {
	data map[string]interface{}
}

func FakeCache() *fakeCache {
	return &fakeCache{
		data: make(map[string]interface{}),
	}
}

func (c *fakeCache) Set(ctx context.Context, k string, v interface{}) error {
	c.data[k] = v
	return nil
}

func (c *fakeCache) Get(ctx context.Context, k string) (interface{}, error) {
	if v, ok := c.data[k]; ok {
		return v, nil
	}

	return nil, errors.ErrTODO
}

func (c *fakeCache) Delete(ctx context.Context, k string) error {
	delete(c.data, k)
	return nil
}
