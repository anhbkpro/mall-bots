package logging

import (
	"context"
	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/application/commands"
	"eda-in-golang/ordering/internal/application/queries"
	"eda-in-golang/ordering/internal/domain"

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

func (a Application) CreateOrder(ctx context.Context, cmd commands.CreateOrder) error {
	a.logger.Info().Msg("--> ordering.CreateOrder")
	defer func() {
		a.logger.Info().Msg("<-- ordering.CreateOrder")
	}()

	return a.App.CreateOrder(ctx, cmd)
}

func (a Application) ReadyOrder(ctx context.Context, cmd commands.ReadyOrder) error {
	a.logger.Info().Msg("--> ordering.ReadyOrder")
	defer func() {
		a.logger.Info().Msg("<-- ordering.ReadyOrder")
	}()

	return a.App.ReadyOrder(ctx, cmd)
}

func (a Application) CancelOrder(ctx context.Context, cmd commands.CancelOrder) error {
	a.logger.Info().Msg("--> ordering.CancelOrder")
	defer func() {
		a.logger.Info().Msg("<-- ordering.CancelOrder")
	}()

	return a.App.CancelOrder(ctx, cmd)
}

func (a Application) CompleteOrder(ctx context.Context, cmd commands.CompleteOrder) error {
	a.logger.Info().Msg("--> ordering.CompleteOrder")
	defer func() {
		a.logger.Info().Msg("<-- ordering.CompleteOrder")
	}()

	return a.App.CompleteOrder(ctx, cmd)
}

func (a Application) GetOrder(ctx context.Context, query queries.GetOrder) (*domain.Order, error) {
	a.logger.Info().Msg("--> ordering.GetOrder")
	defer func() {
		a.logger.Info().Msg("<-- ordering.GetOrder")
	}()

	return a.App.GetOrder(ctx, query)
}
