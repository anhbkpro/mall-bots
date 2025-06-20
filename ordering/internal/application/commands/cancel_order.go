package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"

	"github.com/stackus/errors"
)

type CancelOrder struct {
	ID string
}

type CancelOrderHandler struct {
	orderRepo       domain.OrderRepository
	shoppingRepo    domain.ShoppingRepository
	domainPublisher ddd.EventPublisher
}

func NewCancelOrderHandler(orderRepo domain.OrderRepository, shoppingRepo domain.ShoppingRepository, domainPublisher ddd.EventPublisher) CancelOrderHandler {
	return CancelOrderHandler{
		orderRepo:       orderRepo,
		shoppingRepo:    shoppingRepo,
		domainPublisher: domainPublisher,
	}
}

func (h CancelOrderHandler) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := h.orderRepo.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "find order")
	}

	if err = order.Cancel(); err != nil {
		return errors.Wrap(err, "cancel order")
	}

	if err = h.shoppingRepo.Cancel(ctx, order.ShoppingID); err != nil {
		return errors.Wrap(err, "cancel shopping")
	}

	if err = h.orderRepo.Update(ctx, order); err != nil {
		return errors.Wrap(err, "update order")
	}

	// publish order canceled event
	if err = h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return errors.Wrap(err, "publish order canceled event")
	}

	return nil
}
