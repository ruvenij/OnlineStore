package request

import "OnlieStore/internal/util"

type Order struct {
	UserID    string `json:"-"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
	Price     string `json:"price" validate:"required"`
	ProductID string `json:"product_id" validate:"required"`
}

type OrderDetail struct {
	OrderID string
	Status  util.OrderStatus
}
