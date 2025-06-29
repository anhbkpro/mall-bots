package domain

const (
	BasketStartedEvent     = "baskets.BasketStarted"
	BasketItemAddedEvent   = "baskets.BasketItemAdded"
	BasketItemRemovedEvent = "baskets.BasketItemRemoved"
	BasketCanceledEvent    = "baskets.BasketCanceled"
	BasketCheckedOutEvent  = "baskets.BasketCheckedOut"
)

type BasketStarted struct {
	CustomerID string
}

// BasketStarted implements the registry.Registrable interface
// => to be able to register the event with the registry.
func (BasketStarted) Key() string {
	return BasketStartedEvent
}

type BasketItemAdded struct {
	Item Item
}

// BasketItemAdded implements the registry.Registrable interface
// => to be able to register the event with the registry.
func (BasketItemAdded) Key() string {
	return BasketItemAddedEvent
}

type BasketItemRemoved struct {
	ProductID string
	Quantity  int
}

// BasketItemRemoved implements the registry.Registrable interface
// => to be able to register the event with the registry.
func (BasketItemRemoved) Key() string {
	return BasketItemRemovedEvent
}

type BasketCanceled struct {
}

// BasketCanceled implements the registry.Registrable interface
// => to be able to register the event with the registry.
func (BasketCanceled) Key() string {
	return BasketCanceledEvent
}

type BasketCheckedOut struct {
	PaymentID  string
	CustomerID string
	Items      map[string]Item
}

// BasketCheckedOut implements the registry.Registrable interface
// => to be able to register the event with the registry.
func (BasketCheckedOut) Key() string {
	return BasketCheckedOutEvent
}
