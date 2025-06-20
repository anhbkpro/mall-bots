package postgres

import (
	"context"
	"database/sql"
	"eda-in-golang/payments/internal/application"
	"eda-in-golang/payments/internal/models"
	"fmt"

	"github.com/pkg/errors"
)

type PaymentRepository struct {
	tableName string
	db        *sql.DB
}

var _ application.PaymentRepository = (*PaymentRepository)(nil)

func NewPaymentRepository(tableName string, db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PaymentRepository) Find(ctx context.Context, id string) (*models.Payment, error) {
	query := "SELECT customer_id, amount FROM %s WHERE id = $1 LIMIT 1"

	payment := &models.Payment{
		ID: id,
	}
	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(&payment.CustomerID, &payment.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan payment")
	}

	return payment, nil
}

func (r PaymentRepository) Save(ctx context.Context, payment *models.Payment) error {
	query := "INSERT INTO %s (id, customer_id, amount) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET customer_id = $2, amount = $3"

	_, err := r.db.ExecContext(ctx, r.table(query), payment.ID, payment.CustomerID, payment.Amount)
	if err != nil {
		return errors.Wrap(err, "failed to save payment")
	}

	return nil
}

func (r PaymentRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
