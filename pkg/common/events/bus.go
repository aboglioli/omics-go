package events

import "context"

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Suscribe(ctx context.Context, eventName string) (<-chan Message, error)
}
