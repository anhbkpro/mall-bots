package logging

import (
	"context"
	"eda-in-golang/baskets/internal/application"
	"eda-in-golang/baskets/internal/domain"

	"github.com/rs/zerolog"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(application application.App, logger zerolog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) StartBasket(ctx context.Context, cmd application.StartBasket) error {
	a.logger.Info().Msg("--> Baskets.StartBasket")
	defer func() {
		a.logger.Info().Msg("<-- Baskets.StartBasket")
	}()

	return a.App.StartBasket(ctx, cmd)
}

func (a Application) CancelBasket(ctx context.Context, cmd application.CancelBasket) error {
	a.logger.Info().Msg("--> Baskets.CancelBasket")
	defer func() {
		a.logger.Info().Msg("<-- Baskets.CancelBasket")
	}()

	return a.App.CancelBasket(ctx, cmd)
}

func (a Application) CheckoutBasket(ctx context.Context, cmd application.CheckoutBasket) error {
	a.logger.Info().Msg("--> Baskets.CheckoutBasket")
	defer func() {
		a.logger.Info().Msg("<-- Baskets.CheckoutBasket")
	}()
	return a.App.CheckoutBasket(ctx, cmd)
}

func (a Application) AddItem(ctx context.Context, cmd application.AddItem) error {
	a.logger.Info().Msg("--> Baskets.AddItem")
	defer func() {
		a.logger.Info().Msg("<-- Baskets.AddItem")
	}()
	return a.App.AddItem(ctx, cmd)
}

func (a Application) RemoveItem(ctx context.Context, cmd application.RemoveItem) error {
	a.logger.Info().Msg("--> Baskets.RemoveItem")
	defer func() {
		a.logger.Info().Msg("<-- Baskets.RemoveItem")
	}()
	return a.App.RemoveItem(ctx, cmd)
}

func (a Application) GetBasket(ctx context.Context, cmd application.GetBasket) (basket *domain.Basket, err error) {
	a.logger.Info().Msg("--> Baskets.GetBasket")
	defer func() {
		a.logger.Info().Msg("<-- Baskets.GetBasket")
	}()
	return a.App.GetBasket(ctx, cmd)
}
