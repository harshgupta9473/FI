package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/harshgupta9473/fi/dto"
)

type ProductsRepository struct {
	db *sql.DB
}

type ProductsRepoIntf interface {
	AddProduct(ctx context.Context, product *dto.Product) (int64, error)
	GetProducts(ctx context.Context, limit int, offset int) ([]*dto.Product, error)
	UpdateProductQuantityByID(ctx context.Context, id int64, quantity int64) error
	GetProductByID(ctx context.Context, id int64) (*dto.Product, error)
}

func NewProductsRepository(db *sql.DB) ProductsRepoIntf {
	return &ProductsRepository{
		db: db,
	}
}

func (p *ProductsRepository) AddProduct(ctx context.Context, product *dto.Product) (int64, error) {
	query := `
		INSERT INTO products (name, type, sku, image_url, description, quantity, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int64
	err := p.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Type,
		product.SKU,
		product.ImageURL,
		product.Description,
		product.Quantity,
		product.Price,
	).Scan(&id)

	return id, err
}

func (p *ProductsRepository) GetProducts(ctx context.Context, limit int, offset int) ([]*dto.Product, error) {
	query := `SELECT id, name, type, sku, image_url, description, quantity, price FROM products
              ORDER BY id
              LIMIT $1 OFFSET $2`
	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*dto.Product
	for rows.Next() {
		var product dto.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Type,
			&product.SKU,
			&product.ImageURL,
			&product.Description,
			&product.Quantity,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (p *ProductsRepository) UpdateProductQuantityByID(ctx context.Context, id int64, quantity int64) error {
	query := `
		UPDATE products
		SET quantity = $1
		WHERE id = $2
	`
	result, err := p.db.ExecContext(ctx, query, quantity, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no product found with id %d", id)
	}

	return nil
}

func (p *ProductsRepository) GetProductByID(ctx context.Context, id int64) (*dto.Product, error) {
	query := `SELECT id, name, type, sku, image_url, description, quantity, price FROM products WHERE id=$1 `
	var product dto.Product
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Type,
		&product.SKU,
		&product.ImageURL,
		&product.Description,
		&product.Quantity,
		&product.Price,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}
