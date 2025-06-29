package commands

import (
	"context"
	"log"

	"eda-in-golang/stores/internal/domain"
)

type (
	CreateStoreCmd struct {
		ID       string
		Name     string
		Location string
	}

	CreateStoreHandler struct {
		stores domain.StoreRepository // concrete implementation: es.AggregateRepository[*domain.Store] => see module.go
	}
)

func NewCreateStoreHandler(stores domain.StoreRepository) CreateStoreHandler {
	return CreateStoreHandler{
		stores: stores,
	}
}

// CreateStore is the command handler for creating a store
func (h CreateStoreHandler) CreateStore(ctx context.Context, cmd CreateStoreCmd) error {
	log.Printf("CreateStore: Starting to create store with ID=%s, Name=%s, Location=%s", cmd.ID, cmd.Name, cmd.Location)

	// store is an implementation of es.EventSourcedAggregate
	store, err := domain.CreateStore(cmd.ID, cmd.Name, cmd.Location) // create the store aggregate
	if err != nil {
		log.Printf("CreateStore: Failed to create store domain object: %v", err)
		return err
	}

	log.Printf("CreateStore: Successfully created store domain object with ID=%s", store.ID())

	// domain.StoreRepository interface (es.AggregateRepository[*domain.Store] implements the domain.StoreRepository interface)
	// h.stores.Save(ctx, store) = domain.StoreRepository.Save(ctx, store)
	// => es.AggregateRepository[*domain.Store].Save(ctx, store)
	// => es.AggregateStore.Save(ctx, aggregate)
	// => postgres.EventStore.Save(ctx, aggregate)

	err = h.stores.Save(ctx, store) // save the store aggregate to the event store
	if err != nil {
		log.Printf("CreateStore: Failed to save store to repository: %v", err)
		return err
	}

	log.Printf("CreateStore: Successfully saved store with ID=%s to repository", store.ID())
	return nil
}
