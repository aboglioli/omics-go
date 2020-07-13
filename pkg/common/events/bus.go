//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package events

import "context"

type Message interface {
	Event() interface{}
	Ack() bool
}

type Subscription interface {
	Message() <-chan Message
	Unsubscribe() error
}

type EventPublisher interface {
	Publish(ctx context.Context, code string, event interface{}) error
}

type EventSubscriber interface {
	Subscribe(ctx context.Context, code string) (Subscription, error)
}
