package rest

import (
	"context"
	"eda-in-golang/customers/customerspb"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterGateway(ctx context.Context, mux *chi.Mux, grpcAddr string) error {
	apiRoot := "/api/customers"

	gateway := runtime.NewServeMux()

	err := customerspb.RegisterCustomersServiceHandlerFromEndpoint(ctx, gateway, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		return err
	}

	// mount the gateway under the api root
	mux.Mount(apiRoot, gateway)

	return nil
}
