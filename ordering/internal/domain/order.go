package domain

import (
	"eda-in-golang/internal/ddd"
	"log"

	"github.com/stackus/errors"
)

var (
	ErrOrderHasNoItems         = errors.Wrap(errors.ErrBadRequest, "order has no items")
	ErrOrderCannotBeCancelled  = errors.Wrap(errors.ErrBadRequest, "order cannot be cancelled")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "customer id cannot be blank")
	ErrPaymentIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "payment id cannot be blank")
)

type Order struct {
	ddd.AggregateBase
	CustomerID string
	PaymentID  string
	InvoiceID  string
	ShoppingID string
	Items      []*Item
	Status     OrderStatus
}

func CreateOrder(id, customerID, paymentID string, items []*Item) (*Order, error) {
	log.Printf("Domain.CreateOrder: Creating order with ID: %s", id)

	if len(items) == 0 {
		log.Printf("Domain.CreateOrder: Validation failed - no items provided")
		return nil, ErrOrderHasNoItems
	}

	if customerID == "" {
		log.Printf("Domain.CreateOrder: Validation failed - customer ID is blank")
		return nil, ErrCustomerIDCannotBeBlank
	}

	if paymentID == "" {
		log.Printf("Domain.CreateOrder: Validation failed - payment ID is blank")
		return nil, ErrPaymentIDCannotBeBlank
	}

	log.Printf("Domain.CreateOrder: All validations passed, creating order object")

	order := &Order{
		AggregateBase: ddd.AggregateBase{ID: id},
		CustomerID:    customerID,
		PaymentID:     paymentID,
		InvoiceID:     "",
		ShoppingID:    "",
		Items:         items,
		Status:        OrderStatusPending,
	}

	order.AddEvent(OrderCreated{Order: order})

	log.Printf("Domain.CreateOrder: Order created successfully with %d items, status: %s", len(items), order.Status.String())

	return order, nil
}

func (o *Order) Cancel() error {
	if o.Status != OrderStatusPending {
		return ErrOrderCannotBeCancelled
	}

	o.Status = OrderStatusCancelled
	o.AddEvent(OrderCanceled{Order: o})

	return nil
}

func (o *Order) Ready() error {
	o.Status = OrderStatusReady
	o.AddEvent(OrderReadied{Order: o})

	return nil
}

func (o *Order) Complete(invoiceID string) error {
	o.InvoiceID = invoiceID
	o.Status = OrderStatusCompleted
	o.AddEvent(OrderCompleted{Order: o})

	return nil
}

func (o Order) GetTotal() float64 {
	total := 0.0
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
