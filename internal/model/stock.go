package model

type Stock struct {
	ID              string   `json:"id"`
	Product         *Product `json:"product"`
	InitialQuantity int      `json:"initial_quantity"`
	CurrentQuantity int      `json:"current_quantity"`
}
