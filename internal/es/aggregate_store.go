package es

import (
	"context"
	"fmt"

	"eda-in-golang/internal/ddd"
)

type EventSourcedAggregate interface {
	ddd.IDer
	AggregateName() string
	ddd.Eventer
	Versioner
	EventApplier
	EventCommitter
}

type AggregateStoreMiddleware func(store AggregateStore) AggregateStore

// AggregateStore is the interface that wraps the basic Save and Load methods.
// EventPublisher, EventStore, SnapshotStore are all implementations of AggregateStore
// Purpose: chain-of-responsibility pattern
type AggregateStore interface {
	Load(ctx context.Context, aggregate EventSourcedAggregate) error
	Save(ctx context.Context, aggregate EventSourcedAggregate) error
}

func AggregateStoreWithMiddleware(store AggregateStore, mws ...AggregateStoreMiddleware) AggregateStore {
	//	var s AggregateStore
	s := store
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i := len(mws) - 1; i >= 0; i-- {
		fmt.Println("---AggregateStoreWithMiddleware: Applying middleware", i)
		s = mws[i](s)
	}
	return s
}
