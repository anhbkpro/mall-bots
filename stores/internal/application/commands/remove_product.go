package commands

import (
	"context"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"

	"github.com/pkg/errors"
)

type (
	RemoveProduct struct {
		ID string
	}

	RemoveProductHandler struct {
		products        domain.ProductRepository
		domainPublisher ddd.EventPublisher
	}
)

func NewRemoveProductHandler(products domain.ProductRepository, domainPublisher ddd.EventPublisher) RemoveProductHandler {
	return RemoveProductHandler{products: products, domainPublisher: domainPublisher}
}

func (h RemoveProductHandler) RemoveProduct(ctx context.Context, cmd RemoveProduct) error {
	product, err := h.products.Find(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "product not found")
	}

	if err := product.Remove(); err != nil {
		return errors.Wrap(err, "failed to remove product")
	}

	if err := h.products.Delete(ctx, cmd.ID); err != nil {
		return errors.Wrap(err, "failed to delete product")
	}

	// publish domain event
	if err := h.domainPublisher.Publish(ctx, product.GetEvents()...); err != nil {
		return errors.Wrap(err, "failed to publish domain event")
	}

	return nil
}
