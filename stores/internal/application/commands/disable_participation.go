package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"

	"github.com/pkg/errors"
)

type (
	// Command type (DTO)
	DisableParticipation struct {
		ID string
	}

	// Command handler
	DisableParticipationHandler struct {
		stores          domain.StoreRepository
		domainPublisher ddd.EventPublisher
	}
)

func NewDisableParticipationHandler(stores domain.StoreRepository, domainPublisher ddd.EventPublisher) DisableParticipationHandler {
	return DisableParticipationHandler{stores: stores, domainPublisher: domainPublisher}
}

func (h DisableParticipationHandler) DisableParticipation(ctx context.Context, cmd DisableParticipation) error {
	store, err := h.stores.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "store not found")
	}

	if err := store.DisableParticipation(); err != nil {
		return errors.Wrap(err, "failed to disable participation")
	}

	if err := h.stores.Update(ctx, store); err != nil {
		return errors.Wrap(err, "failed to update store")
	}

	if err := h.domainPublisher.Publish(ctx, store.GetEvents()...); err != nil {
		return errors.Wrap(err, "failed to publish domain event")
	}

	return nil
}
