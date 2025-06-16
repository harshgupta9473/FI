package services

import (
	"errors"
	"fmt"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/repository"
)

type ProductService struct {
	ProductRepo repository.ProductsRepoIntf
}

type ProductServiceIntf interface {
	GetALLProducts(page int, limit int) ([]*dto.Product, error)
	AddProduct(product *dto.Product) (int64, error)
	UpdateProduct(id int64, quantity int64) (*dto.Product, error, string)
}

func NewProductService(productRepo repository.ProductsRepoIntf) ProductServiceIntf {
	return &ProductService{
		ProductRepo: productRepo,
	}
}

func (p *ProductService) GetALLProducts(page int, limit int) ([]*dto.Product, error) {
	if limit < 1 {
		return nil, errors.New("limit cannot be less than 1")
	}
	if page < 1 {
		return nil, errors.New("page cannot be less than 1")
	}
	offset := (page - 1) * limit

	products, err := p.ProductRepo.GetProducts(limit, offset)

	if err != nil {
		return nil, err
	}

	if products == nil {
		return []*dto.Product{}, nil
	}
	return products, nil
}

func (p *ProductService) UpdateProduct(id int64, quantity int64) (*dto.Product, error, string) {
	err := p.ProductRepo.UpdateProductQuantityByID(id, quantity)
	if err != nil {
		return nil, fmt.Errorf("error updating product quantity: %v", err), ""
	}
	product, err := p.ProductRepo.GetProductByID(id)
	if err != nil {
		return nil, nil, "updated"
	}
	return product, nil, "updated"
}

func (p *ProductService) AddProduct(product *dto.Product) (int64, error) {
	return p.ProductRepo.AddProduct(product)
}
