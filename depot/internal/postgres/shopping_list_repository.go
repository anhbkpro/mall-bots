package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/ddd"
)

type ShoppingListRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.ShoppingListRepository = (*ShoppingListRepository)(nil)

func NewShoppingListRepository(tableName string, db *sql.DB) ShoppingListRepository {
	return ShoppingListRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r ShoppingListRepository) Find(ctx context.Context, id string) (*domain.ShoppingList, error) {
	const query = "SELECT order_id, items, assigned_bot_id, status FROM %s WHERE id = $1 LIMIT 1"

	shoppingList := &domain.ShoppingList{
		AggregateBase: ddd.AggregateBase{
			ID: id,
		},
	}
	var items []byte
	var status string
	var assignedBotID sql.NullString

	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(&shoppingList.OrderID, &items, &assignedBotID, &status)
	if err != nil {
		return nil, errors.ErrInternalServerError.Err(err)
	}

	// Handle nullable assigned_bot_id
	if assignedBotID.Valid {
		shoppingList.AssignedBotID = assignedBotID.String
	} else {
		shoppingList.AssignedBotID = ""
	}

	shoppingList.Status = domain.ToShoppingListStatus(status)

	err = json.Unmarshal(items, &shoppingList.Items)
	if err != nil {
		return nil, errors.ErrInternalServerError.Err(err)
	}

	return shoppingList, nil
}

func (r ShoppingListRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.ShoppingList, error) {
	const query = "SELECT id, items, assigned_bot_id, status FROM %s WHERE order_id = $1 LIMIT 1"

	shoppingList := &domain.ShoppingList{
		AggregateBase: ddd.AggregateBase{},
	}
	var items []byte
	var status string
	var assignedBotID sql.NullString

	err := r.db.QueryRowContext(ctx, r.table(query), orderID).Scan(&shoppingList.ID, &items, &assignedBotID, &status)
	if err != nil {
		return nil, errors.ErrInternalServerError.Err(err)
	}

	// Handle nullable assigned_bot_id
	if assignedBotID.Valid {
		shoppingList.AssignedBotID = assignedBotID.String
	} else {
		shoppingList.AssignedBotID = ""
	}

	shoppingList.Status = domain.ToShoppingListStatus(status)

	err = json.Unmarshal(items, &shoppingList.Items)
	if err != nil {
		return nil, errors.ErrInternalServerError.Err(err)
	}

	return shoppingList, nil
}

func (r ShoppingListRepository) Save(ctx context.Context, list *domain.ShoppingList) error {
	const query = "INSERT INTO %s (id, order_id, items, assigned_bot_id, status) VALUES ($1, $2, $3, $4, $5)"

	items, err := json.Marshal(list.Items)
	if err != nil {
		return errors.ErrInternalServerError.Err(err)
	}

	// Handle empty assigned_bot_id as NULL for database
	var assignedBotID sql.NullString
	if list.AssignedBotID != "" {
		assignedBotID.String = list.AssignedBotID
		assignedBotID.Valid = true
	}

	_, err = r.db.ExecContext(ctx, r.table(query), list.ID, list.OrderID, items, assignedBotID, list.Status.String())

	return errors.ErrInternalServerError.Err(err)
}

func (r ShoppingListRepository) Update(ctx context.Context, list *domain.ShoppingList) error {
	const query = "UPDATE %s SET items = $2, assigned_bot_id = $3, status = $4 WHERE id = $1"

	items, err := json.Marshal(list.Items)
	if err != nil {
		return errors.ErrInternalServerError.Err(err)
	}

	// Handle empty assigned_bot_id as NULL for database
	var assignedBotID sql.NullString
	if list.AssignedBotID != "" {
		assignedBotID.String = list.AssignedBotID
		assignedBotID.Valid = true
	}

	_, err = r.db.ExecContext(ctx, r.table(query), list.ID, items, assignedBotID, list.Status.String())

	return errors.ErrInternalServerError.Err(err)
}

func (r ShoppingListRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
