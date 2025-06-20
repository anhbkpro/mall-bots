package application

import (
	"context"
	"eda-in-golang/payments/internal/models"

	"github.com/stackus/errors"
)

type (
	AuthorizePayment struct {
		ID         string
		CustomerID string
		Amount     float64
	}

	ConfirmPayment struct {
		ID string
	}

	CreateInvoice struct {
		ID        string
		OrderID   string
		PaymentID string
		Amount    float64
	}

	AdjustInvoice struct {
		ID     string
		Amount float64
	}

	PayInvoice struct {
		ID string
	}

	CancelInvoice struct {
		ID string
	}

	App interface {
		AuthorizePayment(ctx context.Context, cmd AuthorizePayment) error
		ConfirmPayment(ctx context.Context, cmd ConfirmPayment) error
		CreateInvoice(ctx context.Context, cmd CreateInvoice) error
		AdjustInvoice(ctx context.Context, cmd AdjustInvoice) error
		PayInvoice(ctx context.Context, cmd PayInvoice) error
		CancelInvoice(ctx context.Context, cmd CancelInvoice) error
	}

	Application struct {
		invoiceRepo InvoiceRepository
		paymentRepo PaymentRepository
		ordersRepo  OrderRepository
	}
)

var _ App = (*Application)(nil)

func New(
	invoiceRepo InvoiceRepository,
	paymentRepo PaymentRepository,
	ordersRepo OrderRepository,
) *Application {
	return &Application{
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
		ordersRepo:  ordersRepo,
	}
}

func (a Application) AuthorizePayment(ctx context.Context, cmd AuthorizePayment) error {
	return a.paymentRepo.Save(ctx, &models.Payment{
		ID:         cmd.ID,
		CustomerID: cmd.CustomerID,
		Amount:     cmd.Amount,
	})
}

func (a Application) ConfirmPayment(ctx context.Context, cmd ConfirmPayment) error {
	if payment, err := a.paymentRepo.Find(ctx, cmd.ID); err != nil || payment == nil {
		return errors.Wrap(err, "payment not found")
	}

	return nil
}

func (a Application) CreateInvoice(ctx context.Context, cmd CreateInvoice) error {
	return a.invoiceRepo.Save(ctx, &models.Invoice{
		ID:      cmd.ID,
		OrderID: cmd.OrderID,
		Amount:  cmd.Amount,
		Status:  models.InvoiceStatusPending,
	})
}

func (a Application) AdjustInvoice(ctx context.Context, cmd AdjustInvoice) error {
	invoice, err := a.invoiceRepo.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "invoice not found")
	}

	invoice.Amount += cmd.Amount
	return a.invoiceRepo.Update(ctx, invoice)
}

func (a Application) PayInvoice(ctx context.Context, cmd PayInvoice) error {
	invoice, err := a.invoiceRepo.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "invoice not found")
	}

	if invoice.Status != models.InvoiceStatusPending {
		return errors.Wrap(errors.ErrBadRequest, "invoice is not pending")
	}

	invoice.Status = models.InvoiceStatusPaid
	return a.invoiceRepo.Update(ctx, invoice)
}

func (a Application) CancelInvoice(ctx context.Context, cmd CancelInvoice) error {
	invoice, err := a.invoiceRepo.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "invoice not found")
	}

	if invoice.Status != models.InvoiceStatusPending {
		return errors.Wrap(errors.ErrBadRequest, "invoice is not pending")
	}

	invoice.Status = models.InvoiceStatusFailed
	return a.invoiceRepo.Update(ctx, invoice)
}
