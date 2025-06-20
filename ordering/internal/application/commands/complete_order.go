package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"

	"github.com/stackus/errors"
)

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type CompleteOrderHandler struct {
	orderRepo       domain.OrderRepository
	domainPublisher ddd.EventPublisher
}

func NewCompleteOrderHandler(orderRepo domain.OrderRepository, domainPublisher ddd.EventPublisher) CompleteOrderHandler {
	return CompleteOrderHandler{orderRepo: orderRepo, domainPublisher: domainPublisher}
}

func (h CompleteOrderHandler) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	order, err := h.orderRepo.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "find order")
	}

	if err := order.Complete(cmd.InvoiceID); err != nil {
		return errors.Wrap(err, "complete order")
	}

	if err = h.orderRepo.Update(ctx, order); err != nil {
		return errors.Wrap(err, "update order")
	}

	// publish order completed event
	if err = h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return errors.Wrap(err, "publish order completed event")
	}

	return nil
}
