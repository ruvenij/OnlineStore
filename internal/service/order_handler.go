package service

import (
	"OnlieStore/internal/model"
	"OnlieStore/internal/util"
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type OrderService struct {
	mu             sync.RWMutex
	orders         map[string]*model.Order
	ordersByUserID map[string]*list.List
	latestOrderId  int
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders:         make(map[string]*model.Order),
		ordersByUserID: make(map[string]*list.List),
		latestOrderId:  1,
	}
}

func (os *OrderService) AddOrder(order *model.Order) {
	os.mu.Lock()
	defer os.mu.Unlock()

	order.ID = fmt.Sprintf("%05d", os.latestOrderId)
	order.Status = string(util.OrderStatusPlaced)

	os.orders[order.ID] = order

	if os.ordersByUserID[order.UserID] == nil {
		os.ordersByUserID[order.UserID] = list.New()
	}

	// append the latest order to the front, so that can retrieve the latest order first
	os.ordersByUserID[order.UserID].PushFront(order)
	os.latestOrderId++
}

func (os *OrderService) GetOrder(id string) (*model.Order, error) {
	os.mu.RLock()
	defer os.mu.RUnlock()

	o, ok := os.orders[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Order not found, id: %s", id))
	}

	return o, nil
}

func (os *OrderService) GetOrdersByUserID(userID string, params *model.PaginationParams) ([]*model.Order, error) {
	os.mu.RLock()
	defer os.mu.RUnlock()

	orderList, ok := os.ordersByUserID[userID]
	if !ok {
		return []*model.Order{}, nil
	}

	startIndex := (params.Page - 1) * params.Limit
	endIndex := startIndex + params.Limit
	if endIndex >= orderList.Len() {
		endIndex = orderList.Len() - 1
	}
	if startIndex < 0 {
		return []*model.Order{}, errors.New(fmt.Sprintf("Invalid page number for get orders request, page : %d",
			params.Page))
	}

	result := make([]*model.Order, 0)
	i := 0
	for e := orderList.Front(); e != nil && i <= endIndex; e = e.Next() {
		if i >= startIndex {
			result = append(result, e.Value.(*model.Order))
		}
		i++
	}

	return result, nil
}

func (os *OrderService) UpdateOrderStatus(id string, status util.OrderStatus) error {
	os.mu.Lock()
	defer os.mu.Unlock()

	o, ok := os.orders[id]
	if !ok {
		return errors.New(fmt.Sprintf("Order not found, id: %s", id))
	}

	return o.UpdateOrderStatus(status)
}
