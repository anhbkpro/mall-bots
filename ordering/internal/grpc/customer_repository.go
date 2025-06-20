package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"eda-in-golang/customers/customerspb"
	"eda-in-golang/ordering/internal/domain"
)

type CustomerRepository struct {
	client customerspb.CustomersServiceClient
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(conn *grpc.ClientConn) CustomerRepository {
	return CustomerRepository{client: customerspb.NewCustomersServiceClient(conn)}
}

func (r CustomerRepository) Authorize(ctx context.Context, customerID string) error {
	log.Printf("CustomerRepository.Authorize: Authorizing customer: %s", customerID)

	_, err := r.client.AuthorizeCustomer(ctx, &customerspb.AuthorizeCustomerRequest{Id: customerID})
	if err != nil {
		log.Printf("CustomerRepository.Authorize: Authorization failed for customer %s: %v", customerID, err)
		return err
	}

	log.Printf("CustomerRepository.Authorize: Customer %s authorized successfully", customerID)
	return err
}
