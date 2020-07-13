package main

import (
	"context"
	"fmt"
	"omics/pkg/shared/events"
	"sync"
	"time"
)

type ProductCreated struct {
	ID   string
	Name string
}

type ProductUpdated struct {
	ID      string
	OldName string
	NewName string
}

type Ready struct{}

func main() {
	eventBus := events.InMemEventBus()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		sub, err := eventBus.Subscribe(context.Background(), "product.created")
		if err != nil {
			panic(err)
		}

		c := 0
		for msg := range sub.Message() {
			event, ok := msg.Event().(ProductCreated)
			if !ok {
				panic(fmt.Sprintf("Invalid event: %#v", msg.Event()))
			}

			fmt.Printf("[ProductCreated]: %#v\n", event)

			if c >= 3 {
				sub.Unsubscribe()
			}
			c++
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		sub, err := eventBus.Subscribe(context.Background(), "product.updated")
		if err != nil {
			panic(err)
		}

		c := 0
		for msg := range sub.Message() {
			event, ok := msg.Event().(ProductUpdated)
			if !ok {
				panic(fmt.Sprintf("Invalid event: %#v", msg.Event()))
			}

			fmt.Printf("[ProductUpdated]: %#v\n", event)

			if c >= 2 {
				sub.Unsubscribe()
			}
			c++
		}

		wg.Done()
	}()

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(100 * time.Millisecond)
			eventBus.Publish(
				context.Background(),
				"product.created",
				ProductCreated{
					ID:   fmt.Sprintf("P0%d", i),
					Name: fmt.Sprintf("Product #%d", i),
				},
			)

			eventBus.Publish(
				context.Background(),
				"product.updated",
				ProductUpdated{
					ID:      fmt.Sprintf("P0%d", i),
					OldName: fmt.Sprintf("Product #%d", i),
					NewName: fmt.Sprintf("Product #%d", i+255),
				},
			)
		}

		eventBus.Publish(context.Background(), "ready", Ready{})
	}()

	sub, err := eventBus.Subscribe(context.Background(), "ready")
	if err != nil {
		panic("Error subscribing to 'ready'")
	}

	msg := <-sub.Message()
	fmt.Println(msg.Event())

	wg.Wait()
}
