package ddd

import (
	"context"
	"log"
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

	eventName := event.EventName()
	d.handlers[eventName] = append(d.handlers[eventName], handler)
	log.Printf("EventDispatcher.Subscribe: Handler registered for event: %s", eventName)
}

// Handlers are executed synchronously in sequence
func (d *EventDispatcher) Publish(ctx context.Context, events ...Event) error {
	log.Printf("EventDispatcher.Publish: Publishing %d events", len(events))

	// iterate over all events and publish them to all registered handlers
	for i, event := range events {
		eventName := event.EventName()
		log.Printf("EventDispatcher.Publish: Processing event %d: %s", i+1, eventName)

		handlers, ok := d.handlers[eventName]
		if !ok {
			log.Printf("EventDispatcher.Publish: No handlers found for event: %s", eventName)
			continue
		}

		log.Printf("EventDispatcher.Publish: Found %d handlers for event: %s", len(handlers), eventName)

		for j, handler := range handlers {
			log.Printf("EventDispatcher.Publish: Executing handler %d for event: %s", j+1, eventName)
			// executes all registered handlers for event
			if err := handler(ctx, event); err != nil {
				log.Printf("EventDispatcher.Publish: Handler %d failed for event %s: %v", j+1, eventName, err)
				return err
			}
			log.Printf("EventDispatcher.Publish: Handler %d completed successfully for event: %s", j+1, eventName)
		}
	}

	log.Printf("EventDispatcher.Publish: All events published successfully")
	return nil
}
