package models

// we defined the saga data model here
type CreateOrderData struct {
	OrderID    string
	CustomerID string
	PaymentID  string
	ShoppingID string
	Items      []Item
	Total      float64
}

type Item struct {
	ProductID string
	StoreID   string
	Price     float64
	Quantity  int
}
