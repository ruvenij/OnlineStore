package model

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// ProductDetails used when a new product is added by admin
type ProductDetails struct {
	ID            string  `json:"id"` // used when adding new stock, or deleting stock
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Category      string  `json:"category"`
	AddedQuantity int     `json:"addedQuantity"`
}
