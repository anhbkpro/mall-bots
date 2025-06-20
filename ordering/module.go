package ordering

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/monolith"
	"eda-in-golang/ordering/internal/application"
	"eda-in-golang/ordering/internal/grpc"
	"eda-in-golang/ordering/internal/handlers"
	"eda-in-golang/ordering/internal/logging"
	"eda-in-golang/ordering/internal/postgres"
	"eda-in-golang/ordering/internal/rest"
	"fmt"

	"github.com/pkg/errors"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	ordersRepo := postgres.NewOrderRepository("ordering.orders", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return errors.Wrap(err, "dial ordering service")
	}
	customersRepo := grpc.NewCustomerRepository(conn)
	paymentsRepo := grpc.NewPaymentRepository(conn)
	invoiceRepo := grpc.NewInvoiceRepository(conn)
	shoppingRepo := grpc.NewShoppingListRepository(conn)
	notificationsRepo := grpc.NewNotificationRepository(conn)

	// setup application
	var app application.App
	app = application.New(
		ordersRepo,
		customersRepo,
		paymentsRepo,
		shoppingRepo,
		domainDispatcher,
	)
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup application handlers
	notificationHandlers := logging.LogDomainEventHandlerAccess(
		application.NewNotificationHandlers(notificationsRepo),
		mono.Logger(),
	)

	invoiceHandlers := logging.LogDomainEventHandlerAccess(
		application.NewInvoiceHandlers(invoiceRepo),
		mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}
	handlers.RegisterNotificationHandlers(notificationHandlers, domainDispatcher)
	handlers.RegisterInvoiceHandlers(invoiceHandlers, domainDispatcher)

	fmt.Println("ordering module started")
	return nil
}
