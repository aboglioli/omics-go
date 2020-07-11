package events

type Message interface {
	Event() interface{}
	Ack() bool
}

type Subscription interface {
	Message() <-chan Message
	Unsubscribe() error
}
