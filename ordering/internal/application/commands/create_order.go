package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
	"log"

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
	log.Printf("CreateOrder: Starting order creation for ID: %s, Customer: %s, Payment: %s", cmd.ID, cmd.CustomerID, cmd.PaymentID)
	log.Printf("CreateOrder: Items count: %d", len(cmd.Items))

	// Log item details
	for i, item := range cmd.Items {
		log.Printf("CreateOrder: Item %d - ProductID: %s, StoreID: %s, Quantity: %d, Price: %.2f",
			i+1, item.ProductID, item.StoreID, item.Quantity, item.Price)
	}

	order, err := domain.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		log.Printf("CreateOrder: Failed to create order domain object: %v", err)
		return errors.Wrap(err, "create order")
	}
	log.Printf("CreateOrder: Domain order created successfully with status: %s", order.Status.String())

	// authorize customer
	log.Printf("CreateOrder: Authorizing customer: %s", cmd.CustomerID)
	if err := h.customerRepo.Authorize(ctx, cmd.CustomerID); err != nil {
		log.Printf("CreateOrder: Customer authorization failed: %v", err)
		return errors.Wrap(err, "authorize customer")
	}
	log.Printf("CreateOrder: Customer authorized successfully")

	// validate payment
	log.Printf("CreateOrder: Confirming payment: %s", cmd.PaymentID)
	if err := h.paymentRepo.Confirm(ctx, cmd.PaymentID); err != nil {
		log.Printf("CreateOrder: Payment confirmation failed: %v", err)
		return errors.Wrap(err, "confirm payment")
	}
	log.Printf("CreateOrder: Payment confirmed successfully")

	// create order
	log.Printf("CreateOrder: Saving order to repository")
	if err := h.orderRepo.Save(ctx, order); err != nil {
		log.Printf("CreateOrder: Failed to save order: %v", err)
		return errors.Wrap(err, "save order")
	}
	log.Printf("CreateOrder: Order saved successfully to repository")

	// publish order created event
	log.Printf("CreateOrder: Publishing order created events")
	if err := h.domainPublisher.Publish(ctx, order.GetEvents()...); err != nil {
		log.Printf("CreateOrder: Failed to publish order created events: %v", err)
		return errors.Wrap(err, "publish order created event")
	}
	log.Printf("CreateOrder: Order created events published successfully")

	log.Printf("CreateOrder: Order creation completed successfully for ID: %s", cmd.ID)
	return nil
}
