package logging

import (
	"context"
	"eda-in-golang/baskets/internal/application"
	"eda-in-golang/internal/ddd"

	"github.com/rs/zerolog"
)

type DomainEventHandlers struct {
	application.DomainEventHandlers
	logger zerolog.Logger
}

var _ application.DomainEventHandlers = (*DomainEventHandlers)(nil)

func LogDomainEventHandlersAccess(handlers application.DomainEventHandlers, logger zerolog.Logger) DomainEventHandlers {
	return DomainEventHandlers{
		DomainEventHandlers: handlers,
		logger:              logger,
	}
}

func (h DomainEventHandlers) OnBasketStarted(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> Baskets.OnBasketStarted")
	defer func() {
		h.logger.Info().Msg("<-- Baskets.OnBasketStarted")
	}()
	return h.DomainEventHandlers.OnBasketStarted(ctx, event)
}

func (h DomainEventHandlers) OnBasketItemAdded(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> Baskets.OnBasketItemAdded")
	defer func() {
		h.logger.Info().Msg("<-- Baskets.OnBasketItemAdded")
	}()
	return h.DomainEventHandlers.OnBasketItemAdded(ctx, event)
}

func (h DomainEventHandlers) OnBasketItemRemoved(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> Baskets.OnBasketItemRemoved")
	defer func() {
		h.logger.Info().Msg("<-- Baskets.OnBasketItemRemoved")
	}()
	return h.DomainEventHandlers.OnBasketItemRemoved(ctx, event)
}

func (h DomainEventHandlers) OnBasketCanceled(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> Baskets.OnBasketCanceled")
	defer func() {
		h.logger.Info().Msg("<-- Baskets.OnBasketCanceled")
	}()
	return h.DomainEventHandlers.OnBasketCanceled(ctx, event)
}

func (h DomainEventHandlers) OnBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	h.logger.Info().Msg("--> Baskets.OnBasketCheckedOut")
	defer func() {
		h.logger.Info().Msg("<-- Baskets.OnBasketCheckedOut")
	}()
	return h.DomainEventHandlers.OnBasketCheckedOut(ctx, event)
}
