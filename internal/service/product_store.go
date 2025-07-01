package service

import (
	"OnlieStore/internal/model"
	"OnlieStore/internal/util"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type ProductStore struct {
	mu              sync.RWMutex
	stock           map[string]*model.Stock // key - product id, value - product stock
	latestProdIndex int                     // next available index to be used as the product id when adding new product
	stockList       []*model.Stock          // sorted list of products
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		stock:           make(map[string]*model.Stock),
		stockList:       make([]*model.Stock, 0),
		latestProdIndex: 1,
	}
}

func (ps *ProductStore) AddProduct(input *model.ProductDetails) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	productStock := &model.Stock{
		ID: fmt.Sprintf("%s%05d", util.ProductSuffix, ps.latestProdIndex),
		Product: &model.Product{
			ID:       fmt.Sprintf("%s%05d", util.ProductSuffix, ps.latestProdIndex),
			Name:     input.Name,
			Price:    input.Price,
			Category: input.Category,
		},
		InitialQuantity: input.AddedQuantity,
		CurrentQuantity: input.AddedQuantity,
	}

	ps.stock[productStock.ID] = productStock
	ps.stockList = append(ps.stockList, productStock)

	// sort the list when a new product is added
	sort.Slice(
		ps.stockList,
		func(i, j int) bool {
			return ps.stockList[i].ID < ps.stockList[j].ID
		})

	ps.latestProdIndex++
}

func (ps *ProductStore) GetProduct(id string) (*model.Stock, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	p, ok := ps.stock[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Product not found, id: %s", id))
	}

	return p, nil
}

func (ps *ProductStore) IsProductAvailableToBuy(id string, quantity int) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if p, ok := ps.stock[id]; ok && p.CurrentQuantity >= quantity {
		return true
	}

	return false
}

func (ps *ProductStore) GetProducts(params *model.PaginationParams) ([]*model.Stock, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	// product list is in sorted order already
	startIndex := (params.Page - 1) * params.Limit
	endIndex := startIndex + params.Limit

	if startIndex < 0 {
		return nil, errors.New(fmt.Sprintf("Invalid page number received for the request, Page : %d", params.Page))
	}

	if endIndex > len(ps.stockList) {
		endIndex = len(ps.stockList)
	}

	return ps.stockList[startIndex:endIndex], nil
}

func (ps *ProductStore) UpdateProductQuantity(id string, action int, quantity int) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	p, ok := ps.stock[id]
	if !ok {
		return errors.New(fmt.Sprintf("Product %s is not available", id))
	}

	if action == util.ActionProductDecrease {
		p.CurrentQuantity -= quantity // reduce qty because of a user buy action
	} else if action == util.ActionProductIncrease {
		p.InitialQuantity += quantity // increase qty after adding new stocks
		p.CurrentQuantity += quantity
	} else {
		return errors.New(fmt.Sprintf("Invalid action : %d", action))
	}

	return nil
}
