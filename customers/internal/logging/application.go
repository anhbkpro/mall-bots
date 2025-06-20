package logging

import (
	"context"
	"eda-in-golang/customers/internal/application"
	"eda-in-golang/customers/internal/domain"

	"github.com/rs/zerolog"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(app application.App, logger zerolog.Logger) Application {
	return Application{
		App:    app,
		logger: logger,
	}
}

func (a Application) RegisterCustomer(ctx context.Context, cmd application.RegisterCustomer) error {
	a.logger.Info().Msg("--> Customers.RegisterCustomer")
	defer func() {
		a.logger.Info().Msg("<-- Customers.RegisterCustomer")
	}()

	return a.App.RegisterCustomer(ctx, cmd)
}

func (a Application) AuthorizeCustomer(ctx context.Context, cmd application.AuthorizeCustomer) error {
	a.logger.Info().Msg("--> Customers.AuthorizeCustomer")
	defer func() {
		a.logger.Info().Msg("<-- Customers.AuthorizeCustomer")
	}()

	return a.App.AuthorizeCustomer(ctx, cmd)
}

func (a Application) GetCustomer(ctx context.Context, cmd application.GetCustomer) (*domain.Customer, error) {
	a.logger.Info().Msg("--> Customers.GetCustomer")
	defer func() {
		a.logger.Info().Msg("<-- Customers.GetCustomer")
	}()

	return a.App.GetCustomer(ctx, cmd)
}

func (a Application) EnableCustomer(ctx context.Context, cmd application.EnableCustomer) error {
	a.logger.Info().Msg("--> Customers.EnableCustomer")
	defer func() {
		a.logger.Info().Msg("<-- Customers.EnableCustomer")
	}()

	return a.App.EnableCustomer(ctx, cmd)
}

func (a Application) DisableCustomer(ctx context.Context, cmd application.DisableCustomer) error {
	a.logger.Info().Msg("--> Customers.DisableCustomer")
	defer func() {
		a.logger.Info().Msg("<-- Customers.DisableCustomer")
	}()

	return a.App.DisableCustomer(ctx, cmd)
}
