package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"

	"github.com/stackus/errors"
)

type ReadyOrder struct {
	ID string
}

type ReadyOrderHandler struct {
	orderRepo       domain.OrderRepository
	domainPublisher ddd.EventPublisher
}

func NewReadyOrderHandler(orderRepo domain.OrderRepository, domainPublisher ddd.EventPublisher) ReadyOrderHandler {
	return ReadyOrderHandler{orderRepo: orderRepo, domainPublisher: domainPublisher}
}

func (h ReadyOrderHandler) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	order, err := h.orderRepo.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "find order")
	}

	if err := order.Ready(); err != nil {
		return errors.Wrap(err, "ready order")
	}

	if err := h.orderRepo.Update(ctx, order); err != nil {
		return errors.Wrap(err, "update order")
	}

	// publish order readied event
	if err := h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return errors.Wrap(err, "publish order readied event")
	}

	return nil
}
