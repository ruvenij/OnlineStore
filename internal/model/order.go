package model

import (
	"OnlieStore/internal/util"
	"errors"
)

type Order struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	ProductID string  `json:"product_id"`
	Status    string  `json:"status"`
}

func (order *Order) UpdateOrderStatus(newStatus util.OrderStatus) error {
	if order.Status == string(util.OrderStatusCancelled) {
		return errors.New("Order is already cancelled, Unable to update the status ")
	}

	order.Status = string(newStatus)
	return nil
}
