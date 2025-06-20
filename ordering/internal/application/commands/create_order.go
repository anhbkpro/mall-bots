package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"

	"github.com/stackus/errors"
)

type CreateOrder struct {
	ID         string
	CustomerID string
	PaymentID  string
	Items      []*domain.Item
}

type CreateOrderHandler struct {
	orderRepo       domain.OrderRepository
	customerRepo    domain.CustomerRepository
	paymentRepo     domain.PaymentRepository
	shoppingRepo    domain.ShoppingRepository
	domainPublisher ddd.EventPublisher
}

func NewCreateOrderHandler(orderRepo domain.OrderRepository, customerRepo domain.CustomerRepository, paymentRepo domain.PaymentRepository, shoppingRepo domain.ShoppingRepository, domainPublisher ddd.EventPublisher) CreateOrderHandler {
	return CreateOrderHandler{orderRepo: orderRepo, customerRepo: customerRepo, paymentRepo: paymentRepo, shoppingRepo: shoppingRepo, domainPublisher: domainPublisher}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	order, err := domain.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return errors.Wrap(err, "create order")
	}

	// authorize customer
	if err := h.customerRepo.Authorize(ctx, cmd.CustomerID); err != nil {
		return errors.Wrap(err, "authorize customer")
	}

	// validate payment
	if err := h.paymentRepo.Confirm(ctx, cmd.PaymentID); err != nil {
		return errors.Wrap(err, "confirm payment")
	}

	// create order
	if err := h.orderRepo.Save(ctx, order); err != nil {
		return errors.Wrap(err, "save order")
	}

	// publish order created event
	if err := h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		return errors.Wrap(err, "publish order created event")
	}

	return nil
}
