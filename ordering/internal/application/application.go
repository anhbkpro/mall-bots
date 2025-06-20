package application

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/application/commands"
	"eda-in-golang/ordering/internal/application/queries"
	"eda-in-golang/ordering/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateOrder(ctx context.Context, cmd commands.CreateOrder) error
		ReadyOrder(ctx context.Context, cmd commands.ReadyOrder) error
		CancelOrder(ctx context.Context, cmd commands.CancelOrder) error
		CompleteOrder(ctx context.Context, cmd commands.CompleteOrder) error
	}
	Queries interface {
		GetOrder(ctx context.Context, query queries.GetOrder) (*domain.Order, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.CreateOrderHandler
		commands.ReadyOrderHandler
		commands.CancelOrderHandler
		commands.CompleteOrderHandler
	}
	appQueries struct {
		queries.GetOrderHandler
	}
)

var _ App = (*Application)(nil)

func New(orderRepo domain.OrderRepository, customerRepo domain.CustomerRepository, paymentRepo domain.PaymentRepository, shoppingRepo domain.ShoppingRepository, domainPublisher ddd.EventPublisher) *Application {
	return &Application{
		appCommands: appCommands{
			CreateOrderHandler:   commands.NewCreateOrderHandler(orderRepo, customerRepo, paymentRepo, shoppingRepo, domainPublisher),
			ReadyOrderHandler:    commands.NewReadyOrderHandler(orderRepo, domainPublisher),
			CancelOrderHandler:   commands.NewCancelOrderHandler(orderRepo, shoppingRepo, domainPublisher),
			CompleteOrderHandler: commands.NewCompleteOrderHandler(orderRepo, domainPublisher),
		},
		appQueries: appQueries{
			GetOrderHandler: queries.NewGetOrderHandler(orderRepo),
		},
	}
}
