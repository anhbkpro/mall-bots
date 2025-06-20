package models

type InvoiceStatus string

const (
	InvoiceStatusUnknown InvoiceStatus = ""
	InvoiceStatusPending InvoiceStatus = "pending"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusFailed  InvoiceStatus = "failed"
)

type Invoice struct {
	ID      string
	OrderID string
	Amount  float64
	Status  InvoiceStatus
}

func (s InvoiceStatus) String() string {
	switch s {
	case InvoiceStatusPending, InvoiceStatusPaid, InvoiceStatusFailed:
		return string(s)
	default:
		return string(InvoiceStatusUnknown)
	}
}
