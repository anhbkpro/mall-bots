package es

import (
	"context"
	"fmt"
	"log"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

// Why have both AggregateRepository and AggregateStore?

// This structs implements the domain.StoreRepository and ProductRepository interfaces by using generics
type AggregateRepository[T EventSourcedAggregate] struct {
	aggregateName string
	registry      registry.Registry
	store         AggregateStore
}

func NewAggregateRepository[T EventSourcedAggregate](aggregateName string, registry registry.Registry, store AggregateStore) AggregateRepository[T] {
	return AggregateRepository[T]{
		aggregateName: aggregateName,
		registry:      registry,
		store:         store,
	}
}

func (r AggregateRepository[T]) Load(ctx context.Context, aggregateID string) (agg T, err error) {
	log.Printf("AggregateRepository[%s].Load: Starting to load aggregate with ID=%s", r.aggregateName, aggregateID)

	var v any
	v, err = r.registry.Build(
		r.aggregateName,
		ddd.SetID(aggregateID),
		ddd.SetName(r.aggregateName),
	)
	if err != nil {
		log.Printf("AggregateRepository[%s].Load: Failed to build aggregate from registry: %v", r.aggregateName, err)
		return agg, err
	}

	var ok bool
	if agg, ok = v.(T); !ok {
		log.Printf("AggregateRepository[%s].Load: Type assertion failed, expected %T but got %T", r.aggregateName, agg, v)
		return agg, fmt.Errorf("%T is not the expected type %T", v, agg)
	}

	log.Printf("AggregateRepository[%s].Load: Loading aggregate data from store for ID=%s", r.aggregateName, aggregateID)
	if err = r.store.Load(ctx, agg); err != nil {
		log.Printf("AggregateRepository[%s].Load: Failed to load aggregate from store: %v", r.aggregateName, err)
		return agg, err
	}

	log.Printf("AggregateRepository[%s].Load: Successfully loaded aggregate with ID=%s, Version=%d", r.aggregateName, aggregateID, agg.Version())
	return agg, nil
}

func (r AggregateRepository[T]) Save(ctx context.Context, aggregate T) error {
	aggregateID := aggregate.ID()
	log.Printf("AggregateRepository[%s].Save: Starting to save aggregate with ID=%s, Version=%d, PendingVersion=%d",
		r.aggregateName, aggregateID, aggregate.Version(), aggregate.PendingVersion())

	if aggregate.Version() == aggregate.PendingVersion() {
		log.Printf("AggregateRepository[%s].Save: No changes detected for aggregate ID=%s, skipping save", r.aggregateName, aggregateID)
		return nil
	}

	events := aggregate.Events()
	log.Printf("AggregateRepository[%s].Save: Applying %d pending events for aggregate ID=%s", r.aggregateName, len(events), aggregateID)

	for i, event := range events {
		if err := aggregate.ApplyEvent(event); err != nil {
			log.Printf("AggregateRepository[%s].Save: Failed to apply event %d for aggregate ID=%s: %v", r.aggregateName, i, aggregateID, err)
			return err
		}
		log.Printf("AggregateRepository[%s].Save: Successfully applied event %d for aggregate ID=%s", r.aggregateName, i, aggregateID)
	}

	log.Printf("AggregateRepository[%s].Save: Saving aggregate to store for ID=%s", r.aggregateName, aggregateID)
	// => es.AggregateStore.Save(ctx, aggregate)
	// => postgres.EventStore.Save(ctx, aggregate)
	// !This will save the aggregate to the event store (events table)
	// !AND trigger the event publisher (EventPublisher.Save)
	err := r.store.Save(ctx, aggregate)
	if err != nil {
		log.Printf("AggregateRepository[%s].Save: Failed to save aggregate to store: %v", r.aggregateName, err)
		return err
	}

	log.Printf("AggregateRepository[%s].Save: Successfully saved aggregate to store for ID=%s, committing events", r.aggregateName, aggregateID)
	aggregate.CommitEvents()
	log.Printf("AggregateRepository[%s].Save: Successfully committed events for aggregate ID=%s", r.aggregateName, aggregateID)

	return nil
}
