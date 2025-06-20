package domain

type InvoiceID string

func (i InvoiceID) String() string {
	return string(i)
}
