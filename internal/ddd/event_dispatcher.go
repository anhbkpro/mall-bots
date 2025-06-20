package ddd

import (
	"context"
	"sync"
)

type EventSubscriber interface {
	Subscribe(event Event, handler EventHandler)
}

type EventPublisher interface {
	Publish(ctx context.Context, events ...Event) error
}

type EventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.Mutex
}

var _ interface {
	EventSubscriber
	EventPublisher
} = (*EventDispatcher)(nil)

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Multiple handlers can subscribe to the same event
func (d *EventDispatcher) Subscribe(event Event, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers[event.EventName()] = append(d.handlers[event.EventName()], handler)
}

// Handlers are executed synchronously in sequence
func (d *EventDispatcher) Publish(ctx context.Context, events ...Event) error {
	// iterate over all events and publish them to all registered handlers
	for _, event := range events {
		handlers, ok := d.handlers[event.EventName()]
		if !ok {
			continue
		}

		for _, handler := range handlers {
			// executes all registered handlers for event
			if err := handler(ctx, event); err != nil {
				return err
			}
		}
	}

	return nil
}
