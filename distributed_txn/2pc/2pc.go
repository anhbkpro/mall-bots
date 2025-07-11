package pc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Transaction Coordinator
type TransactionCoordinator struct {
	orderService     *OrderService
	inventoryService *InventoryService
}

type Order struct {
	ID       int    `json:"id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type PrepareRequest struct {
	TransactionID string `json:"transaction_id"`
	Order         Order  `json:"order"`
}

type PrepareResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Order Service
type OrderService struct {
	db *sql.DB
}

func (os *OrderService) PrepareOrder(ctx context.Context, txID string, order Order) error {
	// Start transaction
	tx, err := os.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if transaction is prepared

	// Execute business logic
	_, err = tx.ExecContext(ctx,
		"INSERT INTO orders (id, product, quantity, status) VALUES ($1, $2, $3, 'pending')",
		order.ID, order.Product, order.Quantity)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Prepare the transaction
	_, err = tx.ExecContext(ctx, "PREPARE TRANSACTION $1", txID)
	if err != nil {
		return fmt.Errorf("failed to prepare transaction: %w", err)
	}

	log.Printf("Order service: Transaction %s prepared successfully", txID)
	return nil
}

func (os *OrderService) CommitPrepared(ctx context.Context, txID string) error {
	_, err := os.db.ExecContext(ctx, "COMMIT PREPARED $1", txID)
	if err != nil {
		return fmt.Errorf("failed to commit prepared transaction: %w", err)
	}
	log.Printf("Order service: Transaction %s committed", txID)
	return nil
}

func (os *OrderService) RollbackPrepared(ctx context.Context, txID string) error {
	_, err := os.db.ExecContext(ctx, "ROLLBACK PREPARED $1", txID)
	if err != nil {
		return fmt.Errorf("failed to rollback prepared transaction: %w", err)
	}
	log.Printf("Order service: Transaction %s rolled back", txID)
	return nil
}

// HTTP handlers for Order Service
func (os *OrderService) prepareHandler(w http.ResponseWriter, r *http.Request) {
	var req PrepareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := os.PrepareOrder(r.Context(), req.TransactionID, req.Order)

	response := PrepareResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Order prepared successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (os *OrderService) commitHandler(w http.ResponseWriter, r *http.Request) {
	txID := mux.Vars(r)["txid"]
	err := os.CommitPrepared(r.Context(), txID)

	response := PrepareResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Transaction committed"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (os *OrderService) rollbackHandler(w http.ResponseWriter, r *http.Request) {
	txID := mux.Vars(r)["txid"]
	err := os.RollbackPrepared(r.Context(), txID)

	response := PrepareResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Transaction rolled back"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Inventory Service
type InventoryService struct {
	db *sql.DB
}

func (is *InventoryService) PrepareInventoryUpdate(ctx context.Context, txID string, product string, quantity int) error {
	tx, err := is.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if we have enough inventory
	var currentQuantity int
	err = tx.QueryRowContext(ctx, "SELECT quantity FROM inventory WHERE product = $1 FOR UPDATE", product).Scan(&currentQuantity)
	if err != nil {
		return fmt.Errorf("failed to check inventory: %w", err)
	}

	if currentQuantity < quantity {
		return fmt.Errorf("insufficient inventory: have %d, need %d", currentQuantity, quantity)
	}

	// Update inventory
	_, err = tx.ExecContext(ctx,
		"UPDATE inventory SET quantity = quantity - $1 WHERE product = $2",
		quantity, product)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Prepare the transaction
	_, err = tx.ExecContext(ctx, "PREPARE TRANSACTION $1", txID)
	if err != nil {
		return fmt.Errorf("failed to prepare transaction: %w", err)
	}

	log.Printf("Inventory service: Transaction %s prepared successfully", txID)
	return nil
}

func (is *InventoryService) CommitPrepared(ctx context.Context, txID string) error {
	_, err := is.db.ExecContext(ctx, "COMMIT PREPARED $1", txID)
	if err != nil {
		return fmt.Errorf("failed to commit prepared transaction: %w", err)
	}
	log.Printf("Inventory service: Transaction %s committed", txID)
	return nil
}

func (is *InventoryService) RollbackPrepared(ctx context.Context, txID string) error {
	_, err := is.db.ExecContext(ctx, "ROLLBACK PREPARED $1", txID)
	if err != nil {
		return fmt.Errorf("failed to rollback prepared transaction: %w", err)
	}
	log.Printf("Inventory service: Transaction %s rolled back", txID)
	return nil
}

// HTTP handlers for Inventory Service
func (is *InventoryService) prepareHandler(w http.ResponseWriter, r *http.Request) {
	var req PrepareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := is.PrepareInventoryUpdate(r.Context(), req.TransactionID, req.Order.Product, req.Order.Quantity)

	response := PrepareResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Inventory prepared successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (is *InventoryService) commitHandler(w http.ResponseWriter, r *http.Request) {
	txID := mux.Vars(r)["txid"]
	err := is.CommitPrepared(r.Context(), txID)

	response := PrepareResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Transaction committed"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (is *InventoryService) rollbackHandler(w http.ResponseWriter, r *http.Request) {
	txID := mux.Vars(r)["txid"]
	err := is.RollbackPrepared(r.Context(), txID)

	response := PrepareResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Transaction rolled back"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Transaction Coordinator Implementation
func (tc *TransactionCoordinator) ProcessOrder(ctx context.Context, order Order) error {
	// Generate unique transaction ID
	txID := uuid.New().String()

	log.Printf("Starting distributed transaction %s for order %d", txID, order.ID)

	// Phase 1: Prepare both services
	orderPrepared := false
	inventoryPrepared := false

	// Prepare Order Service
	if err := tc.orderService.PrepareOrder(ctx, "order-"+txID, order); err != nil {
		log.Printf("Order service preparation failed: %v", err)
		return fmt.Errorf("order preparation failed: %w", err)
	}
	orderPrepared = true

	// Prepare Inventory Service
	if err := tc.inventoryService.PrepareInventoryUpdate(ctx, "inventory-"+txID, order.Product, order.Quantity); err != nil {
		log.Printf("Inventory service preparation failed: %v", err)
		// Rollback order service
		if orderPrepared {
			tc.orderService.RollbackPrepared(ctx, "order-"+txID)
		}
		return fmt.Errorf("inventory preparation failed: %w", err)
	}
	inventoryPrepared = true

	// Phase 2: Commit both services
	log.Printf("Both services prepared, committing transaction %s", txID)

	// Commit Order Service
	if err := tc.orderService.CommitPrepared(ctx, "order-"+txID); err != nil {
		log.Printf("Order service commit failed: %v", err)
		// Try to rollback inventory
		if inventoryPrepared {
			tc.inventoryService.RollbackPrepared(ctx, "inventory-"+txID)
		}
		return fmt.Errorf("order commit failed: %w", err)
	}

	// Commit Inventory Service
	if err := tc.inventoryService.CommitPrepared(ctx, "inventory-"+txID); err != nil {
		log.Printf("Inventory service commit failed: %v", err)
		// At this point, order is already committed, which is problematic
		// In production, you'd need more sophisticated recovery mechanisms
		return fmt.Errorf("inventory commit failed: %w", err)
	}

	log.Printf("Distributed transaction %s completed successfully", txID)
	return nil
}

// HTTP handler for the coordinator
func (tc *TransactionCoordinator) processOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set timeout for the entire distributed transaction
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	err := tc.ProcessOrder(ctx, order)

	response := PrepareResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Order processed successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Initialize database connections (you'd use separate DBs in practice)
	orderDB, err := sql.Open("postgres", "postgres://user:password@localhost/orderdb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to order database:", err)
	}
	defer orderDB.Close()

	inventoryDB, err := sql.Open("postgres", "postgres://user:password@localhost/inventorydb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to inventory database:", err)
	}
	defer inventoryDB.Close()

	// Initialize services
	orderService := &OrderService{db: orderDB}
	inventoryService := &InventoryService{db: inventoryDB}
	coordinator := &TransactionCoordinator{
		orderService:     orderService,
		inventoryService: inventoryService,
	}

	// Setup routes for Order Service (port 8081)
	orderRouter := mux.NewRouter()
	orderRouter.HandleFunc("/prepare", orderService.prepareHandler).Methods("POST")
	orderRouter.HandleFunc("/commit/{txid}", orderService.commitHandler).Methods("POST")
	orderRouter.HandleFunc("/rollback/{txid}", orderService.rollbackHandler).Methods("POST")

	go func() {
		log.Println("Order service starting on :8081")
		log.Fatal(http.ListenAndServe(":8081", orderRouter))
	}()

	// Setup routes for Inventory Service (port 8082)
	inventoryRouter := mux.NewRouter()
	inventoryRouter.HandleFunc("/prepare", inventoryService.prepareHandler).Methods("POST")
	inventoryRouter.HandleFunc("/commit/{txid}", inventoryService.commitHandler).Methods("POST")
	inventoryRouter.HandleFunc("/rollback/{txid}", inventoryService.rollbackHandler).Methods("POST")

	go func() {
		log.Println("Inventory service starting on :8082")
		log.Fatal(http.ListenAndServe(":8082", inventoryRouter))
	}()

	// Setup routes for Coordinator (port 8080)
	coordinatorRouter := mux.NewRouter()
	coordinatorRouter.HandleFunc("/process-order", coordinator.processOrderHandler).Methods("POST")

	log.Println("Transaction coordinator starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", coordinatorRouter))
}
