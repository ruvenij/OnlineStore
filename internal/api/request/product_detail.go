package request

type ProductDetails struct {
	Name          string `json:"name" validate:"required,min=5,max=15"`
	Price         string `json:"price" validate:"required"`
	Category      string `json:"category" validate:"required"`
	AddedQuantity int    `json:"addedQuantity" validate:"required,gt=0"`
}
