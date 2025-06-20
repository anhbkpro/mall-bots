package queries

import (
	"context"
	"eda-in-golang/ordering/internal/domain"
	"log"

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
	log.Printf("GetOrderHandler.GetOrder: Processing query for order ID: %s", query.ID)

	log.Printf("GetOrderHandler.GetOrder: Calling repository to find order")
	order, err := h.orderRepo.Find(ctx, query.ID)
	if err != nil {
		log.Printf("GetOrderHandler.GetOrder: Repository Find failed for order %s: %v", query.ID, err)
		return nil, errors.Wrap(err, "find order")
	}

	log.Printf("GetOrderHandler.GetOrder: Successfully retrieved order %s with %d items, status: %s",
		query.ID, len(order.Items), order.Status.String())
	return order, nil
}
