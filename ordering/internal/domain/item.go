package domain

type Item struct {
	ProductID   string  `json:"product_id"`
	StoreID     string  `json:"store_id"`
	StoreName   string  `json:"store_name"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}
