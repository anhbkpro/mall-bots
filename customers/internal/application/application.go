package application

import (
	"context"
	"eda-in-golang/customers/internal/domain"
	"eda-in-golang/internal/ddd"
)

type (
	// [DTO] RegisterCustomer is the command to register a new customer
	RegisterCustomer struct {
		ID        string
		Name      string
		SmsNumber string
	}

	AuthorizeCustomer struct {
		ID string
	}

	GetCustomer struct {
		ID string
	}

	EnableCustomer struct {
		ID string
	}

	DisableCustomer struct {
		ID string
	}

	App interface {
		RegisterCustomer(ctx context.Context, cmd RegisterCustomer) error
		AuthorizeCustomer(ctx context.Context, cmd AuthorizeCustomer) error
		GetCustomer(ctx context.Context, cmd GetCustomer) (*domain.Customer, error)
		EnableCustomer(ctx context.Context, cmd EnableCustomer) error
		DisableCustomer(ctx context.Context, cmd DisableCustomer) error
	}

	Application struct {
		customers       domain.CustomerRepository
		domainPublisher ddd.EventPublisher[ddd.AggregateEvent]
	}
)

var _ App = (*Application)(nil)

func New(customers domain.CustomerRepository, domainPublisher ddd.EventPublisher[ddd.AggregateEvent]) Application {
	return Application{
		customers:       customers,
		domainPublisher: domainPublisher,
	}
}

func (a Application) RegisterCustomer(ctx context.Context, cmd RegisterCustomer) error {
	customer, err := domain.RegisterCustomer(cmd.ID, cmd.Name, cmd.SmsNumber)
	if err != nil {
		return err
	}

	if err := a.customers.Save(ctx, customer); err != nil {
		return err
	}

	// publish events to the domain event bus
	if err := a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}

func (a Application) AuthorizeCustomer(ctx context.Context, cmd AuthorizeCustomer) error {
	customer, err := a.customers.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := customer.Authorize(); err != nil {
		return err
	}

	// publish events to the domain event bus
	if err := a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}

func (a Application) GetCustomer(ctx context.Context, cmd GetCustomer) (*domain.Customer, error) {
	return a.customers.Find(ctx, cmd.ID)
}

func (a Application) EnableCustomer(ctx context.Context, cmd EnableCustomer) error {
	customer, err := a.customers.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := customer.Enable(); err != nil {
		return err
	}

	if err := a.customers.Update(ctx, customer); err != nil {
		return err
	}

	if err := a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}

func (a Application) DisableCustomer(ctx context.Context, cmd DisableCustomer) error {
	customer, err := a.customers.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := customer.Disable(); err != nil {
		return err
	}

	if err := a.customers.Update(ctx, customer); err != nil {
		return err
	}

	if err := a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}
