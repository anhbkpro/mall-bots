package grpc

import (
	"context"
	"eda-in-golang/payments/internal/application"
	"eda-in-golang/payments/paymentspb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	app application.App
	paymentspb.UnimplementedPaymentsServiceServer
}

var _ paymentspb.PaymentsServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	paymentspb.RegisterPaymentsServiceServer(registrar, &server{app: app})
	return nil
}

func (s server) AuthorizePayment(ctx context.Context, req *paymentspb.AuthorizePaymentRequest) (*paymentspb.AuthorizePaymentResponse, error) {
	id := uuid.New().String()
	err := s.app.AuthorizePayment(ctx, application.AuthorizePayment{
		ID:         id,
		CustomerID: req.GetCustomerId(),
		Amount:     req.GetAmount(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "authorize payment: %v", err)
	}

	return &paymentspb.AuthorizePaymentResponse{Id: id}, nil
}
