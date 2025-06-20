package postgres

import (
	"context"
	"database/sql"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type OrderRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(tableName string, db *sql.DB) *OrderRepository {
	return &OrderRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r OrderRepository) Find(ctx context.Context, id string) (*domain.Order, error) {
	query := r.table("SELECT customer_id, payment_id, shopping_id, invoice_id, items, status FROM orders WHERE id = $1 LIMIT 1")

	order := &domain.Order{
		AggregateBase: ddd.AggregateBase{
			ID: id,
		},
	}

	var items []byte
	var status string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.CustomerID,
		&order.PaymentID,
		&order.ShoppingID,
		&order.InvoiceID,
		&items,
		&status,
	)
	if err != nil {
		return nil, errors.Wrap(err, "query row context")
	}

	order.Status = domain.OrderStatus(status)

	if err := json.Unmarshal(items, &order.Items); err != nil {
		return nil, errors.Wrap(err, "unmarshal items")
	}

	return order, nil
}

func (r OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	query := "INSERT INTO orders (id, customer_id, payment_id, shopping_id, invoice_id, items, status) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	items, err := json.Marshal(order.Items)
	if err != nil {
		return errors.Wrap(err, "marshal items")
	}

	_, err = r.db.ExecContext(ctx, r.table(query), order.ID, order.CustomerID, order.PaymentID, order.ShoppingID, order.InvoiceID, items, order.Status)
	if err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}

func (r OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := "UPDATE orders SET customer_id = $1, payment_id = $2, shopping_id = $3, invoice_id = $4, items = $5, status = $6 WHERE id = $7"

	items, err := json.Marshal(order.Items)
	if err != nil {
		return errors.Wrap(err, "marshal items")
	}

	_, err = r.db.ExecContext(ctx, r.table(query), order.CustomerID, order.PaymentID, order.ShoppingID, order.InvoiceID, items, order.Status, order.ID)
	if err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}

func (r OrderRepository) table(query string) string {
	return fmt.Sprintf("%s.%s", r.tableName, query)
}
