package logging

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/application"

	"github.com/rs/zerolog"
)

type DomainEventHandlers struct {
	application.DomainEventHandlers
	logger zerolog.Logger
}

var _ application.DomainEventHandlers = (*DomainEventHandlers)(nil)

func LogDomainEventHandlerAccess(handlers application.DomainEventHandlers, logger zerolog.Logger) DomainEventHandlers {
	return DomainEventHandlers{
		DomainEventHandlers: handlers,
		logger:              logger,
	}
}

func (h DomainEventHandlers) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> ordering.OnOrderCreated")
	defer func() {
		h.logger.Info().Msg("<-- ordering.OnOrderCreated")
	}()

	return h.DomainEventHandlers.OnOrderCreated(ctx, event)
}

func (h DomainEventHandlers) OnOrderCanceled(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> ordering.OnOrderCanceled")
	defer func() {
		h.logger.Info().Msg("<-- ordering.OnOrderCanceled")
	}()

	return h.DomainEventHandlers.OnOrderCanceled(ctx, event)
}

func (h DomainEventHandlers) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> ordering.OnOrderReadied")
	defer func() {
		h.logger.Info().Msg("<-- ordering.OnOrderReadied")
	}()

	return h.DomainEventHandlers.OnOrderReadied(ctx, event)
}
