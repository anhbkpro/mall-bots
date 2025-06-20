package application

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type NotificationHandlers struct {
	notificationRepo domain.NotificationRepository
	ignoreUnimplementedDomainEvents
}

func NewNotificationHandlers(notificationRepo domain.NotificationRepository) *NotificationHandlers {
	return &NotificationHandlers{notificationRepo: notificationRepo}
}

func (h *NotificationHandlers) OnOrderCreated(ctx context.Context, event ddd.Event) error {
	orderCreated := event.(domain.OrderCreated)
	return h.notificationRepo.NotifyOrderCreated(ctx, orderCreated.Order.ID, orderCreated.Order.CustomerID)
}

func (h *NotificationHandlers) OnOrderCanceled(ctx context.Context, event ddd.Event) error {
	orderCanceled := event.(domain.OrderCanceled)
	return h.notificationRepo.NotifyOrderCanceled(ctx, orderCanceled.Order.ID, orderCanceled.Order.CustomerID)
}

func (h *NotificationHandlers) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	orderReadied := event.(domain.OrderReadied)
	return h.notificationRepo.NotifyOrderReady(ctx, orderReadied.Order.ID, orderReadied.Order.CustomerID)
}
