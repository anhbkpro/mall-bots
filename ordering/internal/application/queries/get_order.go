package queries

import (
	"context"
	"eda-in-golang/ordering/internal/domain"

	"github.com/pkg/errors"
)

type GetOrder struct {
	ID string
}

type GetOrderHandler struct {
	orderRepo domain.OrderRepository
}

func NewGetOrderHandler(orderRepo domain.OrderRepository) GetOrderHandler {
	return GetOrderHandler{orderRepo: orderRepo}
}

func (h GetOrderHandler) GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error) {
	order, err := h.orderRepo.Find(ctx, query.ID)
	if err != nil {
		return nil, errors.Wrap(err, "find order")
	}

	return order, nil
}
