package events

import "context"

type EventBus interface {
	Publish(ctx context.Context, code string, event interface{}) error
	Subscribe(ctx context.Context, code string) (Subscription, error)
}
