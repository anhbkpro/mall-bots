package application

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

type InvoiceHandlers struct {
	invoiceRepo domain.InvoiceRepository
	ignoreUnimplementedDomainEvents
}

func NewInvoiceHandlers(invoiceRepo domain.InvoiceRepository) *InvoiceHandlers {
	return &InvoiceHandlers{invoiceRepo: invoiceRepo}
}

func (h *InvoiceHandlers) OnOrderReadied(ctx context.Context, event ddd.Event) error {
	orderReadied := event.(domain.OrderReadied)
	return h.invoiceRepo.Save(ctx, orderReadied.Order.ID, orderReadied.Order.PaymentID, orderReadied.Order.GetTotal())
}
