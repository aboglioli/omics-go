package cache

import "context"

type Cache interface {
	Set(ctx context.Context, k string, v interface{}) error
	Get(ctx context.Context, k string) (interface{}, error)
	Delete(ctx context.Context, k string) error
}
