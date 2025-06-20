package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"eda-in-golang/ordering/internal/domain"
	"eda-in-golang/payments/paymentspb"
)

type PaymentRepository struct {
	client paymentspb.PaymentsServiceClient
}

var _ domain.PaymentRepository = (*PaymentRepository)(nil)

func NewPaymentRepository(conn *grpc.ClientConn) PaymentRepository {
	return PaymentRepository{
		client: paymentspb.NewPaymentsServiceClient(conn),
	}
}

func (r PaymentRepository) Confirm(ctx context.Context, paymentID string) error {
	log.Printf("PaymentRepository.Confirm: Confirming payment: %s", paymentID)

	_, err := r.client.ConfirmPayment(ctx, &paymentspb.ConfirmPaymentRequest{Id: paymentID})
	if err != nil {
		log.Printf("PaymentRepository.Confirm: Payment confirmation failed for payment %s: %v", paymentID, err)
		return err
	}

	log.Printf("PaymentRepository.Confirm: Payment %s confirmed successfully", paymentID)
	return err
}
