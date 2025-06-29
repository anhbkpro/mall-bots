package domain

type BasketV1 struct {
	CustomerID string
	PaymentID  string
	Status     BasketStatus
	Items      map[string]Item
}

// BasketV1 implements the es.Snapshot interface
func (BasketV1) SnapshotName() string { return "baskets.BasketV1" }
