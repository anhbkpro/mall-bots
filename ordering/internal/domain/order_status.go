package domain

type OrderStatus string

const (
	OrderStatusUnknown   OrderStatus = ""
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusInProcess OrderStatus = "in-process"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending, OrderStatusInProcess, OrderStatusReady, OrderStatusCompleted, OrderStatusCancelled:
		return string(s)
	default:
		return ""
	}
}

func ToOrderStatus(s string) OrderStatus {
	switch s {
	case OrderStatusPending.String():
		return OrderStatusPending
	case OrderStatusInProcess.String():
		return OrderStatusInProcess
	case OrderStatusReady.String():
		return OrderStatusReady
	case OrderStatusCompleted.String():
		return OrderStatusCompleted
	case OrderStatusCancelled.String():
		return OrderStatusCancelled
	default:
		return OrderStatusUnknown
	}
}
