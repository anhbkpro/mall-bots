package grpc

import (
	"context"
	"eda-in-golang/customers/customerspb"
	"eda-in-golang/customers/internal/application"
	"eda-in-golang/customers/internal/domain"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	app application.App
	customerspb.UnimplementedCustomersServiceServer
}

var _ customerspb.CustomersServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	customerspb.RegisterCustomersServiceServer(registrar, &server{app: app})
	return nil
}

func (s server) RegisterCustomer(ctx context.Context, req *customerspb.RegisterCustomerRequest) (*customerspb.RegisterCustomerResponse, error) {
	id := uuid.New().String()
	err := s.app.RegisterCustomer(ctx, application.RegisterCustomer{
		ID:        id,
		Name:      req.Name,
		SmsNumber: req.SmsNumber,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register customer: %v", err)
	}

	return &customerspb.RegisterCustomerResponse{Id: id}, nil
}

func (s server) AuthorizeCustomer(ctx context.Context, req *customerspb.AuthorizeCustomerRequest) (*customerspb.AuthorizeCustomerResponse, error) {
	err := s.app.AuthorizeCustomer(ctx, application.AuthorizeCustomer{
		ID: req.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "authorize customer: %v", err)
	}
	return &customerspb.AuthorizeCustomerResponse{}, nil
}

func (s server) GetCustomer(ctx context.Context, req *customerspb.GetCustomerRequest) (*customerspb.GetCustomerResponse, error) {
	customer, err := s.app.GetCustomer(ctx, application.GetCustomer{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get customer: %v", err)
	}

	return &customerspb.GetCustomerResponse{
		Customer: s.customerFromDomain(customer),
	}, nil
}

func (s server) EnableCustomer(ctx context.Context, req *customerspb.EnableCustomerRequest) (*customerspb.EnableCustomerResponse, error) {
	err := s.app.EnableCustomer(ctx, application.EnableCustomer{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "enable customer: %v", err)
	}
	return &customerspb.EnableCustomerResponse{}, nil
}

func (s server) DisableCustomer(ctx context.Context, req *customerspb.DisableCustomerRequest) (*customerspb.DisableCustomerResponse, error) {
	err := s.app.DisableCustomer(ctx, application.DisableCustomer{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "disable customer: %v", err)
	}
	return &customerspb.DisableCustomerResponse{}, nil
}

func (s server) customerFromDomain(customer *domain.Customer) *customerspb.Customer {
	return &customerspb.Customer{
		Id:        customer.ID(),
		Name:      customer.Name,
		SmsNumber: customer.SmsNumber,
		Enabled:   customer.Enabled,
	}
}
