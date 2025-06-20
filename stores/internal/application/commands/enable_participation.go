package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"

	"github.com/pkg/errors"
)

type (
	// Command type (DTO)
	EnableParticipation struct {
		ID string
	}

	// Command handler
	EnableParticipationHandler struct {
		stores          domain.StoreRepository
		domainPublisher ddd.EventPublisher
	}
)

func NewEnableParticipationHandler(stores domain.StoreRepository, domainPublisher ddd.EventPublisher) EnableParticipationHandler {
	return EnableParticipationHandler{stores: stores, domainPublisher: domainPublisher}
}

func (h EnableParticipationHandler) EnableParticipation(ctx context.Context, cmd EnableParticipation) error {
	store, err := h.stores.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "store not found")
	}

	if err := store.EnableParticipation(); err != nil {
		return errors.Wrap(err, "failed to enable participation")
	}

	if err := h.stores.Update(ctx, store); err != nil {
		return errors.Wrap(err, "failed to update store")
	}

	if err := h.domainPublisher.Publish(ctx, store.GetEvents()...); err != nil {
		return errors.Wrap(err, "failed to publish domain event")
	}

	return nil
}
