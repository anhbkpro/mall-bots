package application

import (
	"context"
	"eda-in-golang/internal/ddd"
)

type DomainEventHandler interface {
	OnStoreCreated(ctx context.Context, event ddd.Event) error
	OnStoreParticipationEnabled(ctx context.Context, event ddd.Event) error
	OnStoreParticipationDisabled(ctx context.Context, event ddd.Event) error
	OnProductAdded(ctx context.Context, event ddd.Event) error
	OnProductRemoved(ctx context.Context, event ddd.Event) error
}

type ignoreUnimplementedDomainEvent struct{}

var _ DomainEventHandler = (*ignoreUnimplementedDomainEvent)(nil)

func (ignoreUnimplementedDomainEvent) OnStoreCreated(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvent) OnStoreParticipationEnabled(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvent) OnStoreParticipationDisabled(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvent) OnProductAdded(ctx context.Context, event ddd.Event) error {
	return nil
}

func (ignoreUnimplementedDomainEvent) OnProductRemoved(ctx context.Context, event ddd.Event) error {
	return nil
}
