package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/logger"
	"go.uber.org/zap"
)

type ProductsRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

type ProductsRepoIntf interface {
	AddProduct(ctx context.Context, product *dto.Product) (int64, error)
	GetProducts(ctx context.Context, limit int, offset int) ([]*dto.Product, error)
	UpdateProductQuantityByID(ctx context.Context, id int64, quantity int64) error
	GetProductByID(ctx context.Context, id int64) (*dto.Product, error)
}

func NewProductsRepository(db *sql.DB, logger *logger.Logger) ProductsRepoIntf {
	return &ProductsRepository{
		db:     db,
		logger: logger,
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

	if err != nil {
		p.logger.Error(err, "failed to insert product", zap.String("name", product.Name))
		return -1, err
	}

	p.logger.Info("product added successfully", zap.Int64("product_id", id), zap.String("name", product.Name))

	return id, nil
}

func (p *ProductsRepository) GetProducts(ctx context.Context, limit int, offset int) ([]*dto.Product, error) {
	query := `SELECT id, name, type, sku, image_url, description, quantity, price FROM products
              ORDER BY id
              LIMIT $1 OFFSET $2`
	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		p.logger.Error(err, "failed to query products", zap.Int("limit", limit), zap.Int("offset", offset))
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
			p.logger.Error(err, "failed to scan product row")
			return nil, err
		}
		products = append(products, &product)
	}
	p.logger.Info("products retrieved", zap.Int("count", len(products)))
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
		p.logger.Error(err, "failed to update product quantity", zap.Int64("product_id", id), zap.Int64("quantity", quantity))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.logger.Error(err, "failed to get rows affected for update", zap.Int64("product_id", id))
		return err
	}

	if rowsAffected == 0 {
		msg := fmt.Sprintf("no product found with id %d", id)
		p.logger.Info(msg, zap.Int64("product_id", id))
		return fmt.Errorf(msg)
	}

	p.logger.Info("product quantity updated", zap.Int64("product_id", id), zap.Int64("quantity", quantity))

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
			p.logger.Info("product not found", zap.Int64("product_id", id))
			return nil, nil
		}
		p.logger.Error(err, "failed to get product by ID", zap.Int64("product_id", id))
		return nil, err
	}
	p.logger.Info("product retrieved", zap.Int64("product_id", product.ID))
	return &product, nil
}
