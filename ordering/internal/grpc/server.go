package grpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/application/commands"
	"eda-in-golang/ordering/internal/application/queries"
	"eda-in-golang/ordering/internal/domain"
	"eda-in-golang/ordering/orderingpb"
)

type server struct {
	app application.App
	orderingpb.UnimplementedOrderingServiceServer
}

var _ orderingpb.OrderingServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	orderingpb.RegisterOrderingServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateOrder(ctx context.Context, request *orderingpb.CreateOrderRequest) (*orderingpb.CreateOrderResponse, error) {
	id := uuid.New().String()

	log.Printf("gRPC.CreateOrder: Received request for customer: %s, payment: %s, items count: %d",
		request.GetCustomerId(), request.GetPaymentId(), len(request.Items))
	log.Printf("gRPC.CreateOrder: Generated order ID: %s", id)

	items := make([]*domain.Item, 0, len(request.Items))
	for i, item := range request.Items {
		log.Printf("gRPC.CreateOrder: Processing item %d - ProductID: %s, StoreID: %s, Quantity: %d",
			i+1, item.GetProductId(), item.GetStoreId(), item.GetQuantity())
		items = append(items, s.itemToDomain(item))
	}

	log.Printf("gRPC.CreateOrder: Calling application CreateOrder command")
	err := s.app.CreateOrder(ctx, commands.CreateOrder{
		ID:         id,
		CustomerID: request.GetCustomerId(),
		PaymentID:  request.GetPaymentId(),
		Items:      items,
	})

	if err != nil {
		log.Printf("gRPC.CreateOrder: Application CreateOrder failed: %v", err)
		return &orderingpb.CreateOrderResponse{Id: id}, err
	}

	log.Printf("gRPC.CreateOrder: Order created successfully, returning ID: %s", id)
	return &orderingpb.CreateOrderResponse{Id: id}, err
}

func (s server) CancelOrder(ctx context.Context, request *orderingpb.CancelOrderRequest) (*orderingpb.CancelOrderResponse, error) {
	err := s.app.CancelOrder(ctx, commands.CancelOrder{ID: request.GetId()})

	return &orderingpb.CancelOrderResponse{}, err
}

func (s server) ReadyOrder(ctx context.Context, request *orderingpb.ReadyOrderRequest) (*orderingpb.ReadyOrderResponse, error) {
	err := s.app.ReadyOrder(ctx, commands.ReadyOrder{ID: request.GetId()})
	return &orderingpb.ReadyOrderResponse{}, err
}

func (s server) CompleteOrder(ctx context.Context, request *orderingpb.CompleteOrderRequest) (*orderingpb.CompleteOrderResponse, error) {
	err := s.app.CompleteOrder(ctx, commands.CompleteOrder{ID: request.GetId()})
	return &orderingpb.CompleteOrderResponse{}, err
}

func (s server) GetOrder(ctx context.Context, request *orderingpb.GetOrderRequest) (*orderingpb.GetOrderResponse, error) {
	orderID := request.GetId()
	log.Printf("gRPC.GetOrder: Received request for order ID: %s", orderID)

	log.Printf("gRPC.GetOrder: Calling application GetOrder query")
	order, err := s.app.GetOrder(ctx, queries.GetOrder{ID: orderID})
	if err != nil {
		log.Printf("gRPC.GetOrder: Application GetOrder failed: %v", err)
		return nil, err
	}

	log.Printf("gRPC.GetOrder: Order found successfully, converting to protobuf response")
	response := &orderingpb.GetOrderResponse{
		Order: s.orderFromDomain(order),
	}

	log.Printf("gRPC.GetOrder: Successfully returning order %s with %d items", orderID, len(order.Items))
	return response, nil
}

func (s server) orderFromDomain(order *domain.Order) *orderingpb.Order {
	log.Printf("gRPC.orderFromDomain: Converting domain order %s to protobuf", order.ID)

	items := make([]*orderingpb.Item, 0, len(order.Items))
	for i, item := range order.Items {
		log.Printf("gRPC.orderFromDomain: Converting item %d - ProductID: %s, StoreID: %s, Quantity: %d, Price: %.2f",
			i+1, item.ProductID, item.StoreID, item.Quantity, item.Price)
		items = append(items, s.itemFromDomain(item))
	}

	response := &orderingpb.Order{
		Id:         order.ID,
		CustomerId: order.CustomerID,
		PaymentId:  order.PaymentID,
		Items:      items,
		Status:     order.Status.String(),
	}

	log.Printf("gRPC.orderFromDomain: Successfully converted order %s with %d items, status: %s",
		order.ID, len(items), order.Status.String())
	return response
}

func (s server) itemToDomain(item *orderingpb.Item) *domain.Item {
	return &domain.Item{
		ProductID:   item.GetProductId(),
		StoreID:     item.GetStoreId(),
		StoreName:   item.GetStoreName(),
		ProductName: item.GetProductName(),
		Price:       item.GetPrice(),
		Quantity:    int(item.GetQuantity()),
	}
}

func (s server) itemFromDomain(item *domain.Item) *orderingpb.Item {
	log.Printf("gRPC.itemFromDomain: Converting domain item - ProductID: %s, StoreID: %s", item.ProductID, item.StoreID)

	response := &orderingpb.Item{
		StoreId:     item.StoreID,
		ProductId:   item.ProductID,
		StoreName:   item.StoreName,
		ProductName: item.ProductName,
		Price:       item.Price,
		Quantity:    int32(item.Quantity),
	}

	log.Printf("gRPC.itemFromDomain: Successfully converted item - ProductID: %s, StoreID: %s, Quantity: %d, Price: %.2f",
		item.ProductID, item.StoreID, item.Quantity, item.Price)
	return response
}
