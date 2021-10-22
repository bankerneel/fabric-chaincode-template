package schema

import (
	"main/core/utils"
)

//UserTerm structure for contract schema
type Product struct {
	ProductId   int32  `json:"product_id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	utils.MetaData
}

//Product request structure
type ProductRequest struct {
	UUID string `json:"uuid"`
	Product
}

type GetProductRequest struct {
	UserID    int32 `json:"user_id"`
	ProductId int32 `json:"product_id,omitempty" metadata:",optional"`
}

type GetProductResponse struct {
	Message string     `json:"message,omitempty" metadata:",optional"`
	Product []*Product `json:"product"`
}
