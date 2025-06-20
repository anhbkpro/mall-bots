package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type OrderRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(tableName string, db *sql.DB) OrderRepository {
	return OrderRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r OrderRepository) Find(ctx context.Context, orderID string) (*domain.Order, error) {
	const query = "SELECT customer_id, payment_id, shopping_id, invoice_id, items, status FROM %s WHERE id = $1 LIMIT 1"

	log.Printf("OrderRepository.Find: Finding order with ID: %s", orderID)

	order := &domain.Order{
		AggregateBase: ddd.AggregateBase{
			ID: orderID,
		},
	}

	var items []byte
	var status string
	var shoppingID sql.NullString
	var invoiceID sql.NullString

	log.Printf("OrderRepository.Find: Executing database query for order %s", orderID)
	err := r.db.QueryRowContext(ctx, r.table(query), orderID).Scan(&order.CustomerID, &order.PaymentID, &shoppingID, &invoiceID, &items, &status)
	if err != nil {
		log.Printf("OrderRepository.Find: Database query failed for order %s: %v", orderID, err)
		return nil, errors.Wrap(err, "scanning order")
	}

	log.Printf("OrderRepository.Find: Order %s found - Customer: %s, Payment: %s, Status: %s",
		orderID, order.CustomerID, order.PaymentID, status)
	log.Printf("OrderRepository.Find: Raw items JSON length: %d bytes", len(items))
	if len(items) > 0 {
		log.Printf("OrderRepository.Find: Raw items JSON: %s", string(items))
	}

	// Handle nullable fields
	if shoppingID.Valid {
		order.ShoppingID = shoppingID.String
		log.Printf("OrderRepository.Find: ShoppingID set to: %s", order.ShoppingID)
	} else {
		order.ShoppingID = ""
		log.Printf("OrderRepository.Find: ShoppingID is empty")
	}

	if invoiceID.Valid {
		order.InvoiceID = invoiceID.String
		log.Printf("OrderRepository.Find: InvoiceID set to: %s", order.InvoiceID)
	} else {
		order.InvoiceID = ""
		log.Printf("OrderRepository.Find: InvoiceID is empty")
	}

	order.Status = domain.ToOrderStatus(status)

	log.Printf("OrderRepository.Find: Unmarshaling items JSON for order %s", orderID)
	err = json.Unmarshal(items, &order.Items)
	if err != nil {
		log.Printf("OrderRepository.Find: Failed to unmarshal items for order %s: %v", orderID, err)
		log.Printf("OrderRepository.Find: JSON unmarshal error details - items length: %d, items content: %s", len(items), string(items))
		return nil, errors.Wrap(err, "unmarshalling items")
	}

	log.Printf("OrderRepository.Find: Successfully loaded order %s with %d items", orderID, len(order.Items))

	// Log details of each item
	for i, item := range order.Items {
		log.Printf("OrderRepository.Find: Item %d - ProductID: %s, StoreID: %s, StoreName: %s, ProductName: %s, Quantity: %d, Price: %.2f",
			i+1, item.ProductID, item.StoreID, item.StoreName, item.ProductName, item.Quantity, item.Price)
	}

	return order, nil
}

func (r OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	const query = "INSERT INTO %s (id, customer_id, payment_id, shopping_id, invoice_id, items, status) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	log.Printf("OrderRepository.Save: Saving order ID: %s, Customer: %s, Status: %s",
		order.ID, order.CustomerID, order.Status.String())

	items, err := json.Marshal(order.Items)
	if err != nil {
		log.Printf("OrderRepository.Save: Failed to marshal items: %v", err)
		return errors.Wrap(err, "marshalling items")
	}
	log.Printf("OrderRepository.Save: Items marshaled successfully, JSON length: %d", len(items))

	// Handle empty strings as NULL for database
	var shoppingID sql.NullString
	if order.ShoppingID != "" {
		shoppingID.String = order.ShoppingID
		shoppingID.Valid = true
		log.Printf("OrderRepository.Save: ShoppingID set to: %s", order.ShoppingID)
	} else {
		log.Printf("OrderRepository.Save: ShoppingID is empty, will be NULL")
	}

	var invoiceID sql.NullString
	if order.InvoiceID != "" {
		invoiceID.String = order.InvoiceID
		invoiceID.Valid = true
		log.Printf("OrderRepository.Save: InvoiceID set to: %s", order.InvoiceID)
	} else {
		log.Printf("OrderRepository.Save: InvoiceID is empty, will be NULL")
	}

	log.Printf("OrderRepository.Save: Executing INSERT query")
	_, err = r.db.ExecContext(ctx, r.table(query), order.ID, order.CustomerID, order.PaymentID, shoppingID, invoiceID, items, order.Status.String())
	if err != nil {
		log.Printf("OrderRepository.Save: Database INSERT failed: %v", err)
		return errors.Wrap(err, "inserting order")
	}

	log.Printf("OrderRepository.Save: Order saved successfully to database")
	return nil
}

func (r OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	const query = "UPDATE %s SET customer_id = $2, payment_id = $3, shopping_id = $4, invoice_id = $5, items = $6, status = $7 WHERE id = $1"

	items, err := json.Marshal(order.Items)
	if err != nil {
		return errors.Wrap(err, "marshalling items")
	}

	// Handle empty strings as NULL for database
	var shoppingID sql.NullString
	if order.ShoppingID != "" {
		shoppingID.String = order.ShoppingID
		shoppingID.Valid = true
	}

	var invoiceID sql.NullString
	if order.InvoiceID != "" {
		invoiceID.String = order.InvoiceID
		invoiceID.Valid = true
	}

	_, err = r.db.ExecContext(ctx, r.table(query), order.ID, order.CustomerID, order.PaymentID, shoppingID, invoiceID, items, order.Status.String())
	if err != nil {
		return errors.Wrap(err, "updating order")
	}

	return nil
}

func (r OrderRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
