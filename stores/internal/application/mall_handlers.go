package application

import (
	"context"
	"log"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

// implements ddd.EventHandler
type MallHandlers[T ddd.AggregateEvent] struct {
	mall domain.MallRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*MallHandlers[ddd.AggregateEvent])(nil)

func NewMallHandlers(mall domain.MallRepository) *MallHandlers[ddd.AggregateEvent] {
	return &MallHandlers[ddd.AggregateEvent]{
		mall: mall,
	}
}

// The HandleEvent() method simply proxies the event into different methods based on the event name.
func (h MallHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case domain.StoreParticipationEnabledEvent:
		return h.onStoreParticipationEnabled(ctx, event)
	case domain.StoreParticipationDisabledEvent:
		return h.onStoreParticipationDisabled(ctx, event)
	case domain.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}
	return nil
}

func (h MallHandlers[T]) onStoreCreated(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreCreated)
	log.Printf("MallHandlers.onStoreCreated: Starting to add store with ID=%s, Name=%s, Location=%s", event.AggregateID(), payload.Name, payload.Location)
	return h.mall.AddStore(ctx, event.AggregateID(), payload.Name, payload.Location)
}

func (h MallHandlers[T]) onStoreParticipationEnabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.mall.SetStoreParticipation(ctx, event.AggregateID(), true)
}

func (h MallHandlers[T]) onStoreParticipationDisabled(ctx context.Context, event ddd.AggregateEvent) error {
	return h.mall.SetStoreParticipation(ctx, event.AggregateID(), false)
}

func (h MallHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.StoreRebranded)
	return h.mall.RenameStore(ctx, event.AggregateID(), payload.Name)
}
