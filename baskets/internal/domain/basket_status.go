package domain

type BasketStatus string

const (
	BasketStatusUnknown    BasketStatus = ""
	BasketStatusOpen       BasketStatus = "open"
	BasketStatusCanceled   BasketStatus = "canceled"
	BasketStatusCheckedOut BasketStatus = "checked_out"
)

func (s BasketStatus) String() string {
	switch s {
	case BasketStatusOpen, BasketStatusCanceled, BasketStatusCheckedOut:
		return string(s)
	default:
		return string(BasketStatusUnknown)
	}
}

func ToBasketStatus(s string) BasketStatus {
	switch s {
	case BasketStatusOpen.String():
		return BasketStatusOpen
	case BasketStatusCanceled.String():
		return BasketStatusCanceled
	case BasketStatusCheckedOut.String():
		return BasketStatusCheckedOut
	default:
		return BasketStatusUnknown
	}
}
