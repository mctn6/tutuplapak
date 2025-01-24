package dto

import "time"

type CreateProductRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=32"` // Required, minLength: 4, maxLength: 32
	Category string `json:"category" validate:"required"`          // Required, should be an enum of product category types
	Qty      int    `json:"qty" validate:"required,min=1"`         // Required, min: 1
	Price    int    `json:"price" validate:"required,min=100"`     // Required, min: 100
	SKU      string `json:"sku" validate:"required,max=32"`        // Required, maxLength: 32
	FileID   string `json:"fileId" validate:"required"`            // Required, should be a valid fileId
}

type ProductResponse struct {
	ProductID        string    `json:"productId"`        // string | Use any id you want
	Name             string    `json:"name"`             // string
	Category         string    `json:"category"`         // string
	Qty              int       `json:"qty"`              // number
	Price            float64   `json:"price"`            // number
	SKU              string    `json:"sku"`              // string
	FileID           string    `json:"fileId"`           // string
	FileUri          string    `json:"fileUri"`          // related file URI
	FileThumbnailUri string    `json:"fileThumbnailUri"` // related file thumbnail URI
	CreatedAt        time.Time `json:"createdAt"`        // timestamp
	UpdatedAt        time.Time `json:"updatedAt"`        // timestamp
}

type FilterProductRequest struct {
	Limit     int    `form:"limit" binding:"omitempty"`
	Offset    int    `form:"offset" binding:"omitempty"`
	Category  string `form:"category" binding:"omitempty,oneof=Food Beverage Clothes Furniture Tools"`
	ProductId string `form:"productId " binding:"omitempty"`
	SKU       string `form:"sku" binding:"omitempty"`
	SortBy    string `form:"sortBy" binding:"omitempty,oneof=createdAt updatedAt"`
}

type UpdateProductRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=32"` // Required, minLength: 4, maxLength: 32
	Category string `json:"category" validate:"required"`          // Required, should be an enum of product category types
	Qty      int    `json:"qty" validate:"required,min=1"`         // Required, min: 1
	Price    int    `json:"price" validate:"required,min=100"`     // Required, min: 100
	SKU      string `json:"sku" validate:"required,max=32"`        // Required, maxLength: 32
	FileID   string `json:"fileId" validate:"required"`            // Required, should be a valid fileId
}
