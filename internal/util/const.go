package util

const (
	ProductSuffix = "P"
	UserSuffix    = "U"
)

type OrderStatus string

const (
	OrderStatusPlaced    OrderStatus = "placed"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusError     OrderStatus = "error"
)

const (
	ActionProductIncrease = iota
	ActionProductDecrease
)
