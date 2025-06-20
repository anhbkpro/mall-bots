package domain

import (
	"eda-in-golang/internal/ddd"
	"sort"

	"github.com/stackus/errors"
)

var (
	ErrBasketHasNoItems         = errors.Wrap(errors.ErrBadRequest, "basket has no items")
	ErrBasketCannotBeModified   = errors.Wrap(errors.ErrBadRequest, "basket cannot be modified")
	ErrBasketCannotBeCanceled   = errors.Wrap(errors.ErrBadRequest, "basket cannot be canceled")
	ErrQuantityCannotBeNegative = errors.Wrap(errors.ErrBadRequest, "quantity cannot be negative")
	ErrBasketIDCannotBeBlank    = errors.Wrap(errors.ErrBadRequest, "basket id cannot be empty")
	ErrPaymentIDCannotBeBlank   = errors.Wrap(errors.ErrBadRequest, "payment id cannot be empty")
	ErrCustomerIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "customer id cannot be empty")
)

type Basket struct {
	ddd.AggregateBase
	CustomerID string
	PaymentID  string
	Items      []Item
	Status     BasketStatus
}

func StartBasket(id, customerID string) (*Basket, error) {
	if id == "" {
		return nil, ErrBasketIDCannotBeBlank
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	basket := &Basket{
		AggregateBase: ddd.AggregateBase{ID: id},
		CustomerID:    customerID,
		Status:        BasketStatusOpen,
		Items:         make([]Item, 0),
	}

	// StartBasket triggers BasketStarted event
	basket.AddEvent(&BasketStarted{
		Basket: basket,
	})

	return basket, nil

}

func (b Basket) IsCancelable() bool {
	return b.Status == BasketStatusOpen
}

func (b Basket) IsOpen() bool {
	return b.Status == BasketStatusOpen
}

func (b *Basket) Cancel() error {
	if !b.IsCancelable() {
		return ErrBasketCannotBeCanceled
	}

	b.Status = BasketStatusCanceled
	b.Items = make([]Item, 0)

	// Cancel triggers BasketCanceled event
	b.AddEvent(&BasketCanceled{Basket: b})

	return nil
}

func (b *Basket) Checkout(paymentID string) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if len(b.Items) == 0 {
		return ErrBasketHasNoItems
	}

	if paymentID == "" {
		return ErrPaymentIDCannotBeBlank
	}

	b.PaymentID = paymentID
	b.Status = BasketStatusCheckedOut

	// Checkout triggers BasketCheckedOut event
	b.AddEvent(&BasketCheckedOut{Basket: b})

	return nil
}

func (b *Basket) hasProduct(product *Product) (int, bool) {
	for i, item := range b.Items {
		if item.ProductID == product.ID && item.StoreID == product.StoreID {
			return i, true
		}
	}

	return -1, false
}

func (b *Basket) AddItem(store *Store, product *Product, quantity int) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if quantity < 0 {
		return ErrQuantityCannotBeNegative
	}

	item := Item{
		StoreID:      store.ID,
		StoreName:    store.Name,
		ProductID:    product.ID,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		Quantity:     quantity,
	}

	if i, exists := b.hasProduct(product); exists {
		b.Items[i].Quantity += quantity
	} else {
		b.Items = append(b.Items, item)

		// sort items by store name and product name
		sort.Slice(b.Items, func(i, j int) bool {
			if b.Items[i].StoreName == b.Items[j].StoreName {
				return b.Items[i].ProductName < b.Items[j].ProductName
			}
			return b.Items[i].StoreName < b.Items[j].StoreName
		})
	}

	// AddItem triggers BasketItemAdded event
	b.AddEvent(&BasketItemAdded{Basket: b, Item: item})

	return nil
}

func (b *Basket) RemoveItem(product *Product, quantity int) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if quantity < 0 {
		return ErrQuantityCannotBeNegative
	}

	if i, exists := b.hasProduct(product); exists {
		item := b.Items[i]
		item.Quantity -= quantity

		if item.Quantity <= 0 {
			b.Items = append(b.Items[:i], b.Items[i+1:]...)
		} else {
			b.Items[i] = item
		}

		// RemoveItem triggers BasketItemRemoved event
		b.AddEvent(&BasketItemRemoved{Basket: b, Item: item})
	}

	return nil
}
