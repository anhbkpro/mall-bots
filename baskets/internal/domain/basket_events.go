package domain

type BasketStarted struct {
	Basket *Basket
}

// BasketStarted implements the Event interface.
func (BasketStarted) EventName() string {
	return "baskets.BasketStarted"
}

type BasketItemAdded struct {
	Basket *Basket
	Item   Item
}

// BasketItemAdded implements the Event interface.
func (BasketItemAdded) EventName() string {
	return "baskets.BasketItemAdded"
}

type BasketItemRemoved struct {
	Basket *Basket
	Item   Item
}

// BasketItemRemoved implements the Event interface.
func (BasketItemRemoved) EventName() string {
	return "baskets.BasketItemRemoved"
}

type BasketCanceled struct {
	Basket *Basket
}

// BasketCanceled implements the Event interface.
func (BasketCanceled) EventName() string {
	return "baskets.BasketCanceled"
}

type BasketCheckedOut struct {
	Basket *Basket
}

// BasketCheckedOut implements the Event interface, so it is a valid event key in the event dispatcher
func (BasketCheckedOut) EventName() string {
	return "baskets.BasketCheckedOut"
}
