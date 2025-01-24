package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"tutuplapak/db"
	"tutuplapak/dto"
	"tutuplapak/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) CreateProduct(req dto.CreateProductRequest) (models.Product, error) {
	query := `
				WITH inserted_product AS (
					INSERT INTO products (name, category, qty, price, sku, fileId)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING *
				)
				SELECT 
					inserted_product.id,
					inserted_product.name,
					inserted_product.category,
					inserted_product.qty,
					inserted_product.price,
					inserted_product.sku,
					inserted_product.created_at,
					inserted_product.updated_at,
					files.fileid AS file_id,
					files.fileuri AS file_uri,
					files.filethumbnailuri AS file_thumbnail_uri
				FROM inserted_product
				JOIN files ON inserted_product.fileId = files.fileid;
			`

	var product models.Product
	err := db.DB.QueryRow(query, req.Name, req.Category, req.Qty, req.Price, req.SKU, req.FileID).Scan(
		&product.ID,
		&product.Name,
		&product.Category,
		&product.Qty,
		&product.Price,
		&product.SKU,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.File.FileID,
		&product.File.FileUri,
		&product.File.FileThumbnailUri,
	)
	if err != nil {
		return models.Product{}, fmt.Errorf("failed to create product: %v", err)
	}

	return product, nil
}

func (r *ProductRepository) FilterProducts(filters map[string]string) ([]models.Product, error) {
	query := `
		SELECT 
			products.id,
			products.name,
			products.category,
			products.qty,
			products.price,
			products.sku,
			files.fileid,
			files.fileuri,
			files.filethumbnailuri,
			products.created_at,
			products.updated_at
		FROM products
		JOIN files
		ON files.fileid = products.fileid
	`

	args := []interface{}{}
	argCount := 1

	whereClause := " WHERE 1=1"

	for key, value := range filters {
		switch key {
		case "product_id":
			whereClause += fmt.Sprintf(" AND products.id = $%d", argCount)
			args = append(args, value)
			argCount++
		case "category":
			whereClause += fmt.Sprintf(" AND products.category = $%d", argCount)
			args = append(args, value)
			argCount++
		case "sku":
			whereClause += fmt.Sprintf(" AND products.sku = $%d", argCount)
			args = append(args, value)
			argCount++
		default:
			// Ignore unknown filters
			continue
		}
	}

	// Append the WHERE clause to the query
	query += whereClause

	// Handle SORT BY
	if sortBy, ok := filters["sort_by"]; ok {
		switch sortBy {
		case "newest":
			query += " ORDER BY GREATEST(products.created_at, products.updated_at) DESC"
		case "cheapest":
			query += " ORDER BY products.price ASC"
		default:
			if strings.HasPrefix(sortBy, "sold-") {
				// Extract the number of seconds from "sold-x"
				secondsStr := strings.TrimPrefix(sortBy, "sold-")
				seconds, err := strconv.Atoi(secondsStr)
				if err == nil && seconds > 0 {
					// Assuming you have a "sales" table with a "sold_at" timestamp column
					query += fmt.Sprintf(`
						ORDER BY (
							SELECT COUNT(*) 
							FROM sales 
							WHERE sales.product_id = products.id 
							AND sales.sold_at >= NOW() - INTERVAL '%d seconds'
						) DESC`, seconds)
				}
			}
		}
	}

	limit, _ := strconv.Atoi(filters["limit"])
	offset, _ := strconv.Atoi(filters["offset"])
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
		argCount++
	}

	rows, err := r.DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Qty,
			&product.Price,
			&product.SKU,
			&product.File.FileID,
			&product.File.FileUri,
			&product.File.FileThumbnailUri,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetProductById(id int) (*models.Product, error) {
	query := `
		SELECT 
			products.id,
			products.name,
			products.category,
			products.qty,
			products.price,
			products.sku,
			files.fileid,
			files.fileuri,
			files.filethumbnailuri,
			products.created_at,
			products.updated_at
		FROM products
		JOIN files
		ON files.fileid = products.fileid
		WHERE productId = $1
	`

	var product models.Product
	err := r.DB.QueryRowContext(context.Background(), query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Category,
		&product.Qty,
		&product.Price,
		&product.SKU,
		&product.File.FileID,
		&product.File.FileUri,
		&product.File.FileThumbnailUri,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateProduct(id int, updatedProduct models.Product) error {
	query := `
		UPDATE products
		SET name = $1, category = $2, qty = $3, price = $4, sku = $5, fileId = $6
		WHERE id = $7
	`

	_, err := r.DB.Exec(
		query,
		updatedProduct.Name,
		updatedProduct.Category,
		updatedProduct.Qty,
		updatedProduct.Price,
		updatedProduct.SKU,
		updatedProduct.FileID,
		id,
	)
	return err
}

func (r *ProductRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	return nil
}

func (r *ProductRepository) FileExists(fileId string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM files WHERE fileId = $1)`
	err := db.DB.QueryRow(query, fileId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
