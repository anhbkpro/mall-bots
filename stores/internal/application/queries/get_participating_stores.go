package queries

import (
	"context"
	"eda-in-golang/stores/internal/domain"
)

type (
	// Query type (DTO)
	GetParticipatingStores struct {
	}

	// Query handler
	GetParticipatingStoresHandler struct {
		participatingStores domain.ParticipatingStoreRepository
	}
)

func NewGetParticipatingStoresHandler(participatingStores domain.ParticipatingStoreRepository) GetParticipatingStoresHandler {
	return GetParticipatingStoresHandler{participatingStores: participatingStores}
}

func (h GetParticipatingStoresHandler) GetParticipatingStores(ctx context.Context, query GetParticipatingStores) ([]*domain.Store, error) {
	return h.participatingStores.FindAll(ctx)
}
