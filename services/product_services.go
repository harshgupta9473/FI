package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/repository"
)

type ProductService struct {
	ProductRepo repository.ProductsRepoIntf
}

type ProductServiceIntf interface {
	GetALLProducts(ctx context.Context, page int, limit int) ([]*dto.Product, error)
	AddProduct(ctx context.Context, product *dto.Product) (int64, error)
	UpdateProduct(ctx context.Context, id int64, quantity int64) (*dto.Product, error, string)
}

func NewProductService(productRepo repository.ProductsRepoIntf) ProductServiceIntf {
	return &ProductService{
		ProductRepo: productRepo,
	}
}

func (p *ProductService) GetALLProducts(ctx context.Context, page int, limit int) ([]*dto.Product, error) {
	if limit < 1 {
		return nil, errors.New("limit cannot be less than 1")
	}
	if page < 1 {
		return nil, errors.New("page cannot be less than 1")
	}
	offset := (page - 1) * limit

	products, err := p.ProductRepo.GetProducts(ctx, limit, offset)

	if err != nil {
		return nil, err
	}

	if products == nil {
		return []*dto.Product{}, nil
	}
	return products, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, id int64, quantity int64) (*dto.Product, error, string) {
	err := p.ProductRepo.UpdateProductQuantityByID(ctx, id, quantity)
	if err != nil {
		return nil, fmt.Errorf("error updating product quantity: %v", err), ""
	}
	product, err := p.ProductRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, nil, "updated"
	}
	return product, nil, "updated"
}

func (p *ProductService) AddProduct(ctx context.Context, product *dto.Product) (int64, error) {
	return p.ProductRepo.AddProduct(ctx, product)
}
