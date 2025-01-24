package v1

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"tutuplapak/dto"
	"tutuplapak/repositories"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	Repo *repositories.ProductRepository
}

type UpdateProductRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=4,max=32"` // Optional, minLength: 4, maxLength: 32
	Category *string  `json:"category,omitempty" validate:"omitempty"`          // Optional, should be an enum of product category types
	Qty      *int     `json:"qty,omitempty" validate:"omitempty,min=1"`         // Optional, min: 1
	Price    *float64 `json:"price,omitempty" validate:"omitempty,min=100"`     // Optional, min: 100
	SKU      *string  `json:"sku,omitempty" validate:"omitempty,max=32"`        // Optional, maxLength: 32
	FileID   *string  `json:"fileId,omitempty" validate:"omitempty"`            // Optional, should be a valid fileId
}

var validCategories = map[string]bool{
	"Food":      true,
	"Beverage":  true,
	"Clothes":   true,
	"Furniture": true,
	"Tools":     true,
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{
		Repo: repositories.NewProductRepository(db),
	}
}

func validateCategory(category string) bool {
	_, exists := validCategories[category]
	return exists
}

func isValidSortBy(sortBy string) bool {
	// Valid fixed values
	if sortBy == "newest" || sortBy == "cheapest" {
		return true
	}

	// Check if it's a "sold-x" value
	if strings.HasPrefix(sortBy, "sold-") {
		// Extract the number part after "sold-"
		xStr := strings.TrimPrefix(sortBy, "sold-")
		x, err := strconv.Atoi(xStr)
		if err != nil || x < 0 {
			return false // Not a valid number or negative
		}
		return true
	}

	// If none of the above, it's invalid
	return false
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the request struct
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validateCategory(req.Category) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}

	// Validate fileId exists in the database
	exists, err := h.Repo.FileExists(req.FileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate fileId"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileId does not exist"})
		return
	}

	product, err := h.Repo.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ProductResponse{
		ProductID:        strconv.Itoa(product.ID),
		Name:             product.Name,
		Category:         product.Category,
		Qty:              product.Qty,
		Price:            product.Price,
		SKU:              product.SKU,
		FileID:           product.File.FileID,
		FileUri:          product.File.FileUri,
		FileThumbnailUri: product.File.FileThumbnailUri,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter dto.FilterProductRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]string)

	if filter.ProductId != "" {
		filters["product_id"] = filter.ProductId
	}

	if filter.Category != "" {
		filters["category"] = filter.Category
	}

	if filter.SKU != "" {
		filters["sku"] = filter.SKU
	}

	if filter.SortBy != "" {
		filters["sort_by"] = filter.SortBy
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters["limit"] = strconv.Itoa(limit)
	filters["offset"] = strconv.Itoa(offset)

	products, err := h.Repo.FilterProducts(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ProductResponse, 0)
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ProductID:        strconv.Itoa(product.ID),
			Name:             product.Name,
			Category:         product.Category,
			Qty:              product.Qty,
			Price:            product.Price,
			SKU:              product.SKU,
			FileID:           product.File.FileID,
			FileUri:          product.File.FileUri,
			FileThumbnailUri: product.File.FileThumbnailUri,
			CreatedAt:        product.CreatedAt,
			UpdatedAt:        product.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productId := c.Param("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	parsedProductId, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to parse product id"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the request struct
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validateCategory(req.Category) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}

	if err := h.Repo.UpdateProduct(parsedProductId, *product); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	updatedProduct, err := h.Repo.GetProductById(parsedProductId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found after update"})
		return
	}

	response := dto.ProductResponse{
		ProductID:        strconv.Itoa(updatedProduct.ID),
		Name:             updatedProduct.Name,
		Category:         updatedProduct.Category,
		Qty:              updatedProduct.Qty,
		Price:            updatedProduct.Price,
		SKU:              updatedProduct.SKU,
		FileID:           updatedProduct.File.FileID,
		FileUri:          updatedProduct.File.FileUri,
		FileThumbnailUri: updatedProduct.File.FileThumbnailUri,
		CreatedAt:        updatedProduct.CreatedAt,
		UpdatedAt:        updatedProduct.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productId := c.Param("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	parsedProductId, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.Repo.DeleteProduct(parsedProductId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Product deleted")
}
