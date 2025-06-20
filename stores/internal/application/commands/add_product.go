package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"

	"github.com/pkg/errors"
)

type (
	// Command type (DTO)
	AddProduct struct {
		ID          string
		StoreID     string
		Name        string
		Description string
		SKU         string
		Price       float64
	}

	// Command handler
	AddProductHandler struct {
		stores          domain.StoreRepository
		products        domain.ProductRepository
		domainPublisher ddd.EventPublisher
	}
)

func NewAddProductHandler(stores domain.StoreRepository, products domain.ProductRepository, domainPublisher ddd.EventPublisher) AddProductHandler {
	return AddProductHandler{
		stores:          stores,
		products:        products,
		domainPublisher: domainPublisher,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProduct) error {
	if _, err := h.stores.Find(ctx, cmd.StoreID); err != nil {
		return errors.Wrap(err, "store not found")
	}

	product, err := domain.CreateProduct(cmd.ID, cmd.StoreID, cmd.Name, cmd.Description, cmd.SKU, cmd.Price)
	if err != nil {
		return errors.Wrap(err, "failed to create product")
	}

	if err := h.products.Save(ctx, product); err != nil {
		return errors.Wrap(err, "failed to save product")
	}

	// publish domain event
	if err := h.domainPublisher.Publish(ctx, product.GetEvents()...); err != nil {
		return errors.Wrap(err, "failed to publish domain event")
	}

	return nil
}
