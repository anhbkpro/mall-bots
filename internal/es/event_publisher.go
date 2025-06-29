package es

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"eda-in-golang/internal/ddd"
)

type EventPublisher struct {
	AggregateStore
	publisher ddd.EventPublisher[ddd.AggregateEvent]
}

var _ AggregateStore = (*EventPublisher)(nil)

func NewEventPublisher(publisher ddd.EventPublisher[ddd.AggregateEvent]) AggregateStoreMiddleware {
	eventPublisher := EventPublisher{
		publisher: publisher,
	}

	return func(store AggregateStore) AggregateStore {
		eventPublisher.AggregateStore = store
		return eventPublisher
	}
}

// Here we make the publishing action part of the saving action
// No one would need to remember which action needed to be first,
// and the errors from the Save() method can be automatically handled.
func (p EventPublisher) Save(ctx context.Context, aggregate EventSourcedAggregate) error {
	if err := p.AggregateStore.Save(ctx, aggregate); err != nil {
		return err
	}
	caller := getCurrentCaller()
	log.Printf("[Middleware] EventPublisher.Save: Child method called by: %s", caller)
	log.Printf("[Middleware] EventPublisher.Save: Publishing %d events for aggregate ID=%s", len(aggregate.Events()), aggregate.ID())
	// publish the events so the event handlers will be called
	return p.publisher.Publish(ctx, aggregate.Events()...)
}

func getCurrentCaller() string {
	// Skip 2 frames: getCurrentCaller() and the current method
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown caller"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown function"
	}

	return fmt.Sprintf("%s:%d %s", file, line, fn.Name())
}
