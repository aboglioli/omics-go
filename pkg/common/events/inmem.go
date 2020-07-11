package events

import (
	"context"
)

// Message
type inmemMessage struct {
	event interface{}
}

func NewInMemMessage(event interface{}) Message {
	return inmemMessage{
		event: event,
	}
}

func (msg inmemMessage) Event() interface{} {
	return msg.event
}

func (msg inmemMessage) Ack() bool {
	return true
}

// Subscription
type inmemSubscription struct {
	code    string
	msg     chan Message
	deleted bool
}

func (sub *inmemSubscription) Message() <-chan Message {
	return sub.msg
}

func (sub *inmemSubscription) Unsubscribe() error {
	sub.deleted = true
	return nil
}

type inmemEventBus struct {
	subscriptions []*inmemSubscription
}

func InMemEventBus() EventBus {
	return &inmemEventBus{
		subscriptions: make([]*inmemSubscription, 0),
	}
}

func (eb *inmemEventBus) Publish(ctx context.Context, code string, event interface{}) error {
	for _, sub := range eb.subscriptions {
		if !sub.deleted && sub.code == code {
			go func(sub *inmemSubscription) {
				select {
				case sub.msg <- inmemMessage{
					event: event,
				}:
				case <-ctx.Done():
				}
			}(sub)
		} else if sub.deleted {
			subscriptions := make([]*inmemSubscription, 0)
			for _, subscription := range eb.subscriptions {
				if sub == subscription {
					close(sub.msg)
					continue
				}
				subscriptions = append(subscriptions, subscription)
			}
			eb.subscriptions = subscriptions
		}
	}

	return nil
}

func (eb *inmemEventBus) Subscribe(ctx context.Context, code string) (Subscription, error) {
	sub := &inmemSubscription{
		code:    code,
		msg:     make(chan Message),
		deleted: false,
	}
	eb.subscriptions = append(eb.subscriptions, sub)
	return sub, nil
}
