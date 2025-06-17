package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/logger"
	"github.com/harshgupta9473/fi/repository"
	"go.uber.org/zap"
)

type ProductService struct {
	ProductRepo repository.ProductsRepoIntf
	logger      *logger.Logger
}

type ProductServiceIntf interface {
	GetALLProducts(ctx context.Context, page int, limit int) ([]*dto.Product, error)
	AddProduct(ctx context.Context, product *dto.Product) (int64, error)
	UpdateProduct(ctx context.Context, id int64, quantity int64) (*dto.Product, error, string)
}

func NewProductService(productRepo repository.ProductsRepoIntf, logger *logger.Logger) ProductServiceIntf {
	return &ProductService{
		ProductRepo: productRepo,
		logger:      logger,
	}
}

func (p *ProductService) GetALLProducts(ctx context.Context, page int, limit int) ([]*dto.Product, error) {
	if limit < 1 {
		p.logger.Error(errors.New("invalid limit"), "limit must be >= 1", zap.Int("limit", limit))
		return nil, errors.New("limit cannot be less than 1")
	}
	if page < 1 {
		p.logger.Error(errors.New("invalid page"), "page must be >= 1", zap.Int("page", page))
		return nil, errors.New("page cannot be less than 1")
	}
	offset := (page - 1) * limit

	products, err := p.ProductRepo.GetProducts(ctx, limit, offset)

	if err != nil {
		p.logger.Error(err, "failed to get products", zap.Int("limit", limit), zap.Int("offset", offset))
		return nil, err
	}
	p.logger.Info("retrieved products", zap.Int("count", len(products)))
	if products == nil {
		return []*dto.Product{}, nil
	}
	return products, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, id int64, quantity int64) (*dto.Product, error, string) {
	err := p.ProductRepo.UpdateProductQuantityByID(ctx, id, quantity)
	if err != nil {
		p.logger.Error(err, "failed to update product quantity", zap.Int64("product_id", id), zap.Int64("quantity", quantity))
		return nil, fmt.Errorf("error updating product quantity: %v", err), ""
	}

	p.logger.Info("product quantity updated", zap.Int64("product_id", id), zap.Int64("quantity", quantity))

	product, err := p.ProductRepo.GetProductByID(ctx, id)
	if err != nil {
		p.logger.Error(err, "failed to fetch product after update", zap.Int64("product_id", id))
		return nil, nil, "updated"
	}
	p.logger.Info("product fetched after update", zap.Int64("product_id", product.ID))
	return product, nil, "updated"
}

func (p *ProductService) AddProduct(ctx context.Context, product *dto.Product) (int64, error) {
	return p.ProductRepo.AddProduct(ctx, product)
}
