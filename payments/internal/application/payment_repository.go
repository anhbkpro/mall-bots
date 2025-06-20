package application

import (
	"context"
	"eda-in-golang/payments/internal/models"
)

type PaymentRepository interface {
	Find(ctx context.Context, id string) (*models.Payment, error)
	Save(ctx context.Context, payment *models.Payment) error
}
