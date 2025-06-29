package domain

import (
	"context"
)

// There will not be any need to write any new implementations for the two new interfaces;
// AggregateRepository uses generics, and we can again have type safety and save a little on typing.
type StoreRepository interface {
	Load(ctx context.Context, storeID string) (*Store, error)
	Save(ctx context.Context, store *Store) error
}
