package postgres

import (
	"context"
	"database/sql"
	"eda-in-golang/payments/internal/application"
	"eda-in-golang/payments/internal/models"
	"fmt"

	"github.com/pkg/errors"
)

type InvoiceRepository struct {
	tableName string
	db        *sql.DB
}

var _ application.InvoiceRepository = (*InvoiceRepository)(nil)

func NewInvoiceRepository(tableName string, db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r InvoiceRepository) Find(ctx context.Context, id string) (*models.Invoice, error) {
	query := "SELECT order_id, amount, status FROM %s WHERE id = $1 LIMIT 1"

	invoice := &models.Invoice{
		ID: id,
	}
	var status string
	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(&invoice.OrderID, &invoice.Amount, &status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan invoice")
	}

	invoice.Status, err = r.statusToModel(status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert status to model")
	}

	return invoice, nil
}

func (r InvoiceRepository) Save(ctx context.Context, invoice *models.Invoice) error {
	query := "INSERT INTO %s (id, order_id, amount, status) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET order_id = $2, amount = $3, status = $4"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.ID, invoice.OrderID, invoice.Amount, invoice.Status.String())
	if err != nil {
		return errors.Wrap(err, "failed to save invoice")
	}

	return nil
}

func (r InvoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
	query := "UPDATE %s SET order_id = $1, amount = $2, status = $3 WHERE id = $4"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.OrderID, invoice.Amount, invoice.Status.String(), invoice.ID)
	if err != nil {
		return errors.Wrap(err, "failed to update invoice")
	}

	return nil
}

func (r InvoiceRepository) statusToModel(status string) (models.InvoiceStatus, error) {
	switch status {
	case models.InvoiceStatusPending.String():
		return models.InvoiceStatusPending, nil
	case models.InvoiceStatusPaid.String():
		return models.InvoiceStatusPaid, nil
	case models.InvoiceStatusFailed.String():
		return models.InvoiceStatusFailed, nil
	default:
		return models.InvoiceStatusUnknown, errors.Errorf("invalid status: %s", status)
	}
}

func (r InvoiceRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
