package events

type Event interface {
	Name() string
}

type Message interface {
	Event() Event
	Ack() bool
}
