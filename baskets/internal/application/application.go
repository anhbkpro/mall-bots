package application

import (
	"context"
	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/ddd"

	"github.com/pkg/errors"
)

type (
	// Command types (DTOs)
	StartBasket struct {
		ID         string
		CustomerID string
	}

	CancelBasket struct {
		ID string
	}

	CheckoutBasket struct {
		ID        string
		PaymentID string
	}

	AddItem struct {
		ID        string
		ProductID string
		Quantity  int
	}

	RemoveItem struct {
		ID        string
		ProductID string
		Quantity  int
	}

	GetBasket struct {
		ID string
	}

	App interface {
		StartBasket(ctx context.Context, cmd StartBasket) error
		CancelBasket(ctx context.Context, cmd CancelBasket) error
		CheckoutBasket(ctx context.Context, cmd CheckoutBasket) error
		AddItem(ctx context.Context, cmd AddItem) error
		RemoveItem(ctx context.Context, cmd RemoveItem) error
		GetBasket(ctx context.Context, cmd GetBasket) (*domain.Basket, error)
	}

	Application struct {
		basketRepo           domain.BasketRepository
		orderRepo            domain.OrderRepository
		productRepo          domain.ProductRepository
		storeRepo            domain.StoreRepository
		domainEventPublisher ddd.EventPublisher
	}
)

var _ App = (*Application)(nil)

func NewApplication(
	basketRepo domain.BasketRepository,
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
	storeRepo domain.StoreRepository,
	domainEventPublisher ddd.EventPublisher,
) *Application {
	return &Application{
		basketRepo:           basketRepo,
		orderRepo:            orderRepo,
		productRepo:          productRepo,
		storeRepo:            storeRepo,
		domainEventPublisher: domainEventPublisher,
	}
}

func (a Application) StartBasket(ctx context.Context, cmd StartBasket) error {
	basket, err := domain.StartBasket(cmd.ID, cmd.CustomerID)
	if err != nil {
		return err
	}

	if err := a.basketRepo.Save(ctx, basket); err != nil {
		return err
	}

	// publish domain event
	if err := a.domainEventPublisher.Publish(ctx, basket.GetEvents()...); err != nil {
		return err
	}

	return nil
}

func (a Application) CancelBasket(ctx context.Context, cmd CancelBasket) error {
	basket, err := a.basketRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = basket.Cancel()
	if err != nil {
		return err
	}

	if err := a.basketRepo.Save(ctx, basket); err != nil {
		return err
	}

	// publish domain event
	if err := a.domainEventPublisher.Publish(ctx, basket.GetEvents()...); err != nil {
		return err
	}

	return nil

}

func (a Application) CheckoutBasket(ctx context.Context, cmd CheckoutBasket) error {
	basket, err := a.basketRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = basket.Checkout(cmd.PaymentID)
	if err != nil {
		return errors.Wrap(err, "checkout basket")
	}

	if err := a.basketRepo.Update(ctx, basket); err != nil {
		return errors.Wrap(err, "update basket")
	}

	// publish domain event
	if err := a.domainEventPublisher.Publish(ctx, basket.GetEvents()...); err != nil {
		return err
	}

	return nil
}

func (a Application) AddItem(ctx context.Context, cmd AddItem) error {
	basket, err := a.basketRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	product, err := a.productRepo.Find(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	store, err := a.storeRepo.Find(ctx, product.StoreID)
	if err != nil {
		return err
	}

	err = basket.AddItem(store, product, cmd.Quantity)
	if err != nil {
		return err
	}

	if err := a.basketRepo.Update(ctx, basket); err != nil {
		return errors.Wrap(err, "update basket")
	}

	// publish domain event
	if err := a.domainEventPublisher.Publish(ctx, basket.GetEvents()...); err != nil {
		return err
	}

	return nil
}

func (a Application) RemoveItem(ctx context.Context, cmd RemoveItem) error {
	product, err := a.productRepo.Find(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	basket, err := a.basketRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = basket.RemoveItem(product, cmd.Quantity)
	if err != nil {
		return err
	}

	if err := a.basketRepo.Update(ctx, basket); err != nil {
		return errors.Wrap(err, "update basket")
	}

	// publish domain event
	if err := a.domainEventPublisher.Publish(ctx, basket.GetEvents()...); err != nil {
		return err
	}

	return nil
}

func (a Application) GetBasket(ctx context.Context, cmd GetBasket) (*domain.Basket, error) {
	return a.basketRepo.Find(ctx, cmd.ID)
}
