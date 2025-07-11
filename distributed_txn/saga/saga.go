package saga

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

// Saga Orchestrator
type SagaOrchestrator struct {
	orderService     *OrderService
	inventoryService *InventoryService
	sagaStore        *SagaStore
}

type SagaStep struct {
	StepName      string     `json:"step_name"`
	Status        string     `json:"status"` // PENDING, COMPLETED, COMPENSATED, FAILED
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CompensatedAt *time.Time `json:"compensated_at,omitempty"`
}

type SagaTransaction struct {
	ID          string     `json:"id"`
	OrderID     int        `json:"order_id"`
	Status      string     `json:"status"` // STARTED, COMPLETED, COMPENSATING, COMPENSATED, FAILED
	Steps       []SagaStep `json:"steps"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Order struct {
	ID       int    `json:"id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type StepRequest struct {
	SagaID string `json:"saga_id"`
	Order  Order  `json:"order"`
}

type StepResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Saga Store - In production, use a proper database
type SagaStore struct {
	sagas map[string]*SagaTransaction
}

func NewSagaStore() *SagaStore {
	return &SagaStore{
		sagas: make(map[string]*SagaTransaction),
	}
}

func (ss *SagaStore) Save(saga *SagaTransaction) {
	ss.sagas[saga.ID] = saga
}

func (ss *SagaStore) Get(sagaID string) (*SagaTransaction, bool) {
	saga, exists := ss.sagas[sagaID]
	return saga, exists
}

func (ss *SagaStore) UpdateStepStatus(sagaID, stepName, status string) error {
	saga, exists := ss.sagas[sagaID]
	if !exists {
		return fmt.Errorf("saga %s not found", sagaID)
	}

	for i := range saga.Steps {
		if saga.Steps[i].StepName == stepName {
			saga.Steps[i].Status = status
			now := time.Now()
			if status == "COMPLETED" {
				saga.Steps[i].CompletedAt = &now
			} else if status == "COMPENSATED" {
				saga.Steps[i].CompensatedAt = &now
			}
			break
		}
	}
	return nil
}

// Order Service
type OrderService struct {
	db *sql.DB
}

func (os *OrderService) CreateOrder(ctx context.Context, sagaID string, order Order) error {
	// Create order in PENDING state
	_, err := os.db.ExecContext(ctx,
		"INSERT INTO orders (id, product, quantity, status, saga_id) VALUES ($1, $2, $3, 'PENDING', $4)",
		order.ID, order.Product, order.Quantity, sagaID)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	log.Printf("Order service: Order %d created for saga %s", order.ID, sagaID)
	return nil
}

func (os *OrderService) ConfirmOrder(ctx context.Context, sagaID string, orderID int) error {
	// Confirm the order (change status to CONFIRMED)
	result, err := os.db.ExecContext(ctx,
		"UPDATE orders SET status = 'CONFIRMED', updated_at = NOW() WHERE id = $1 AND saga_id = $2",
		orderID, sagaID)
	if err != nil {
		return fmt.Errorf("failed to confirm order: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order %d not found for saga %s", orderID, sagaID)
	}

	log.Printf("Order service: Order %d confirmed for saga %s", orderID, sagaID)
	return nil
}

func (os *OrderService) CancelOrder(ctx context.Context, sagaID string, orderID int) error {
	// Compensating action: cancel/delete the order
	result, err := os.db.ExecContext(ctx,
		"UPDATE orders SET status = 'CANCELLED', updated_at = NOW() WHERE id = $1 AND saga_id = $2",
		orderID, sagaID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order %d not found for saga %s", orderID, sagaID)
	}

	log.Printf("Order service: Order %d cancelled for saga %s", orderID, sagaID)
	return nil
}

// HTTP handlers for Order Service
func (os *OrderService) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := os.CreateOrder(r.Context(), req.SagaID, req.Order)

	response := StepResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Order created successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func (os *OrderService) confirmOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := os.ConfirmOrder(r.Context(), req.SagaID, req.Order.ID)

	response := StepResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Order confirmed successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func (os *OrderService) cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := os.CancelOrder(r.Context(), req.SagaID, req.Order.ID)

	response := StepResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Order cancelled successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

// Inventory Service
type InventoryService struct {
	db *sql.DB
}

func (is *InventoryService) ReserveInventory(ctx context.Context, sagaID string, product string, quantity int) error {
	tx, err := is.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check current inventory
	var currentQuantity int
	err = tx.QueryRowContext(ctx, "SELECT quantity FROM inventory WHERE product = $1 FOR UPDATE", product).Scan(&currentQuantity)
	if err != nil {
		return fmt.Errorf("failed to check inventory: %w", err)
	}

	if currentQuantity < quantity {
		return fmt.Errorf("insufficient inventory: have %d, need %d", currentQuantity, quantity)
	}

	// Create reservation record
	_, err = tx.ExecContext(ctx,
		"INSERT INTO inventory_reservations (saga_id, product, quantity, status) VALUES ($1, $2, $3, 'RESERVED')",
		sagaID, product, quantity)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	// Reduce available inventory
	_, err = tx.ExecContext(ctx,
		"UPDATE inventory SET quantity = quantity - $1 WHERE product = $2",
		quantity, product)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Inventory service: Reserved %d units of %s for saga %s", quantity, product, sagaID)
	return nil
}

func (is *InventoryService) ConfirmReservation(ctx context.Context, sagaID string, product string) error {
	// Mark reservation as confirmed
	result, err := is.db.ExecContext(ctx,
		"UPDATE inventory_reservations SET status = 'CONFIRMED', updated_at = NOW() WHERE saga_id = $1 AND product = $2",
		sagaID, product)
	if err != nil {
		return fmt.Errorf("failed to confirm reservation: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("reservation not found for saga %s and product %s", sagaID, product)
	}

	log.Printf("Inventory service: Confirmed reservation for %s in saga %s", product, sagaID)
	return nil
}

func (is *InventoryService) CancelReservation(ctx context.Context, sagaID string, product string) error {
	tx, err := is.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get reservation details
	var quantity int
	err = tx.QueryRowContext(ctx,
		"SELECT quantity FROM inventory_reservations WHERE saga_id = $1 AND product = $2",
		sagaID, product).Scan(&quantity)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// Return inventory to available stock
	_, err = tx.ExecContext(ctx,
		"UPDATE inventory SET quantity = quantity + $1 WHERE product = $2",
		quantity, product)
	if err != nil {
		return fmt.Errorf("failed to restore inventory: %w", err)
	}

	// Mark reservation as cancelled
	_, err = tx.ExecContext(ctx,
		"UPDATE inventory_reservations SET status = 'CANCELLED', updated_at = NOW() WHERE saga_id = $1 AND product = $2",
		sagaID, product)
	if err != nil {
		return fmt.Errorf("failed to cancel reservation: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Inventory service: Cancelled reservation for %s in saga %s", product, sagaID)
	return nil
}

// HTTP handlers for Inventory Service
func (is *InventoryService) reserveHandler(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := is.ReserveInventory(r.Context(), req.SagaID, req.Order.Product, req.Order.Quantity)

	response := StepResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Inventory reserved successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func (is *InventoryService) confirmHandler(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := is.ConfirmReservation(r.Context(), req.SagaID, req.Order.Product)

	response := StepResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Reservation confirmed successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func (is *InventoryService) cancelHandler(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := is.CancelReservation(r.Context(), req.SagaID, req.Order.Product)

	response := StepResponse{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Reservation cancelled successfully"
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

// Saga Orchestrator Implementation
func (so *SagaOrchestrator) ProcessOrder(ctx context.Context, order Order) (*SagaTransaction, error) {
	// Create saga transaction
	sagaID := uuid.New().String()
	saga := &SagaTransaction{
		ID:        sagaID,
		OrderID:   order.ID,
		Status:    "STARTED",
		CreatedAt: time.Now(),
		Steps: []SagaStep{
			{StepName: "CREATE_ORDER", Status: "PENDING"},
			{StepName: "RESERVE_INVENTORY", Status: "PENDING"},
			{StepName: "CONFIRM_ORDER", Status: "PENDING"},
			{StepName: "CONFIRM_INVENTORY", Status: "PENDING"},
		},
	}

	so.sagaStore.Save(saga)
	log.Printf("Started saga %s for order %d", sagaID, order.ID)

	// Step 1: Create Order
	if err := so.executeStep(ctx, sagaID, "CREATE_ORDER", func() error {
		return so.orderService.CreateOrder(ctx, sagaID, order)
	}); err != nil {
		return so.compensateSaga(ctx, saga, "CREATE_ORDER")
	}

	// Step 2: Reserve Inventory
	if err := so.executeStep(ctx, sagaID, "RESERVE_INVENTORY", func() error {
		return so.inventoryService.ReserveInventory(ctx, sagaID, order.Product, order.Quantity)
	}); err != nil {
		return so.compensateSaga(ctx, saga, "RESERVE_INVENTORY")
	}

	// Step 3: Confirm Order
	if err := so.executeStep(ctx, sagaID, "CONFIRM_ORDER", func() error {
		return so.orderService.ConfirmOrder(ctx, sagaID, order.ID)
	}); err != nil {
		return so.compensateSaga(ctx, saga, "CONFIRM_ORDER")
	}

	// Step 4: Confirm Inventory
	if err := so.executeStep(ctx, sagaID, "CONFIRM_INVENTORY", func() error {
		return so.inventoryService.ConfirmReservation(ctx, sagaID, order.Product)
	}); err != nil {
		return so.compensateSaga(ctx, saga, "CONFIRM_INVENTORY")
	}

	// All steps completed successfully
	saga.Status = "COMPLETED"
	now := time.Now()
	saga.CompletedAt = &now
	so.sagaStore.Save(saga)

	log.Printf("Saga %s completed successfully", sagaID)
	return saga, nil
}

func (so *SagaOrchestrator) executeStep(ctx context.Context, sagaID, stepName string, stepFunc func() error) error {
	log.Printf("Executing step %s for saga %s", stepName, sagaID)

	if err := stepFunc(); err != nil {
		so.sagaStore.UpdateStepStatus(sagaID, stepName, "FAILED")
		log.Printf("Step %s failed for saga %s: %v", stepName, sagaID, err)
		return err
	}

	so.sagaStore.UpdateStepStatus(sagaID, stepName, "COMPLETED")
	log.Printf("Step %s completed for saga %s", stepName, sagaID)
	return nil
}

func (so *SagaOrchestrator) compensateSaga(ctx context.Context, saga *SagaTransaction, failedStep string) (*SagaTransaction, error) {
	log.Printf("Starting compensation for saga %s, failed at step %s", saga.ID, failedStep)
	saga.Status = "COMPENSATING"
	so.sagaStore.Save(saga)

	// Compensate in reverse order of execution
	compensationSteps := []struct {
		stepName       string
		shouldRun      func(string) bool
		compensateFunc func() error
	}{
		{
			stepName: "CONFIRM_INVENTORY",
			shouldRun: func(failed string) bool {
				return failed == "CONFIRM_INVENTORY" || so.isStepCompleted(saga, "CONFIRM_INVENTORY")
			},
			compensateFunc: func() error {
				// No compensation needed for confirmation - it's idempotent
				return nil
			},
		},
		{
			stepName: "CONFIRM_ORDER",
			shouldRun: func(failed string) bool {
				return failed == "CONFIRM_ORDER" || so.isStepCompleted(saga, "CONFIRM_ORDER")
			},
			compensateFunc: func() error {
				// No compensation needed for confirmation - it's idempotent
				return nil
			},
		},
		{
			stepName: "RESERVE_INVENTORY",
			shouldRun: func(failed string) bool {
				return failed == "RESERVE_INVENTORY" || so.isStepCompleted(saga, "RESERVE_INVENTORY")
			},
			compensateFunc: func() error {
				return so.inventoryService.CancelReservation(ctx, saga.ID, getProductFromSaga(saga))
			},
		},
		{
			stepName: "CREATE_ORDER",
			shouldRun: func(failed string) bool {
				return failed == "CREATE_ORDER" || so.isStepCompleted(saga, "CREATE_ORDER")
			},
			compensateFunc: func() error {
				return so.orderService.CancelOrder(ctx, saga.ID, saga.OrderID)
			},
		},
	}

	for _, comp := range compensationSteps {
		if comp.shouldRun(failedStep) {
			if err := comp.compensateFunc(); err != nil {
				log.Printf("Compensation failed for step %s in saga %s: %v", comp.stepName, saga.ID, err)
				saga.Status = "COMPENSATION_FAILED"
				so.sagaStore.Save(saga)
				return saga, fmt.Errorf("compensation failed: %w", err)
			}
			so.sagaStore.UpdateStepStatus(saga.ID, comp.stepName, "COMPENSATED")
		}
	}

	saga.Status = "COMPENSATED"
	so.sagaStore.Save(saga)
	log.Printf("Saga %s compensated successfully", saga.ID)
	return saga, fmt.Errorf("saga compensated due to failure at step %s", failedStep)
}

func (so *SagaOrchestrator) isStepCompleted(saga *SagaTransaction, stepName string) bool {
	for _, step := range saga.Steps {
		if step.StepName == stepName {
			return step.Status == "COMPLETED"
		}
	}
	return false
}

func getProductFromSaga(saga *SagaTransaction) string {
	// In a real implementation, you'd store this data in the saga
	// For now, we'll need to pass it through the context or store it
	return "laptop" // This is a simplified approach
}

// HTTP handler for the orchestrator
func (so *SagaOrchestrator) processOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	saga, err := so.ProcessOrder(ctx, order)

	response := struct {
		Success bool             `json:"success"`
		Message string           `json:"message"`
		Saga    *SagaTransaction `json:"saga,omitempty"`
	}{
		Success: err == nil,
		Message: func() string {
			if err != nil {
				return err.Error()
			}
			return "Order processed successfully"
		}(),
		Saga: saga,
	}

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func (so *SagaOrchestrator) getSagaHandler(w http.ResponseWriter, r *http.Request) {
	sagaID := mux.Vars(r)["sagaId"]

	saga, exists := so.sagaStore.Get(sagaID)
	if !exists {
		http.Error(w, "Saga not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saga)
}

func main_saga() {
	// Initialize database connections
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
	sagaStore := NewSagaStore()
	orchestrator := &SagaOrchestrator{
		orderService:     orderService,
		inventoryService: inventoryService,
		sagaStore:        sagaStore,
	}

	// Setup routes for Order Service (port 8081)
	orderRouter := mux.NewRouter()
	orderRouter.HandleFunc("/create", orderService.createOrderHandler).Methods("POST")
	orderRouter.HandleFunc("/confirm", orderService.confirmOrderHandler).Methods("POST")
	orderRouter.HandleFunc("/cancel", orderService.cancelOrderHandler).Methods("POST")

	go func() {
		log.Println("Order service starting on :8081")
		log.Fatal(http.ListenAndServe(":8081", orderRouter))
	}()

	// Setup routes for Inventory Service (port 8082)
	inventoryRouter := mux.NewRouter()
	inventoryRouter.HandleFunc("/reserve", inventoryService.reserveHandler).Methods("POST")
	inventoryRouter.HandleFunc("/confirm", inventoryService.confirmHandler).Methods("POST")
	inventoryRouter.HandleFunc("/cancel", inventoryService.cancelHandler).Methods("POST")

	go func() {
		log.Println("Inventory service starting on :8082")
		log.Fatal(http.ListenAndServe(":8082", inventoryRouter))
	}()

	// Setup routes for Saga Orchestrator (port 8080)
	orchestratorRouter := mux.NewRouter()
	orchestratorRouter.HandleFunc("/process-order", orchestrator.processOrderHandler).Methods("POST")
	orchestratorRouter.HandleFunc("/saga/{sagaId}", orchestrator.getSagaHandler).Methods("GET")

	log.Println("Saga orchestrator starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", orchestratorRouter))
}
