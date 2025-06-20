package grpc

import (
	"context"
	"eda-in-golang/ordering/orderingpb"
	"eda-in-golang/payments/internal/application"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type OrderRepository struct {
	client orderingpb.OrderingServiceClient
}

var _ application.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(conn *grpc.ClientConn) OrderRepository {
	return OrderRepository{client: orderingpb.NewOrderingServiceClient(conn)}
}

func (r OrderRepository) Complete(ctx context.Context, invoiceID, orderID string) error {
	_, err := r.client.CompleteOrder(ctx, &orderingpb.CompleteOrderRequest{
		Id:        orderID,
		InvoiceId: invoiceID,
	})
	if err != nil {
		return errors.Wrap(err, "complete order")
	}
	return nil
}
