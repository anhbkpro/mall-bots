package domain

type ProductV1 struct {
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

func (ProductV1) SnapshotName() string { return "stores.ProductV1" }

type ProductV2 struct {
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
	Weight      float64
}

func (ProductV2) SnapshotName() string { return "stores.ProductV2" }
