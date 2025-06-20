package queries

import (
	"context"
	"eda-in-golang/stores/internal/domain"
)

type (
	// Query type (DTO)
	GetStores struct {
	}

	// Query handler
	GetStoresHandler struct {
		stores domain.StoreRepository
	}
)

func NewGetStoresHandler(stores domain.StoreRepository) GetStoresHandler {
	return GetStoresHandler{stores: stores}
}

func (h GetStoresHandler) GetStores(ctx context.Context, query GetStores) ([]*domain.Store, error) {
	return h.stores.FindAll(ctx)
}
