package service

import (
	"gokasir-api/models"
	"gokasir-api/repository"
)

type ProductServiceImpl struct {
	repo repository.ProductRepository
}

func NewProductServiceImpl(repo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{repo: repo}
}

func (s *ProductServiceImpl) GetAllProduct() ([]models.Product, error) {
	return s.repo.FindAllProduct()
}

func (s *ProductServiceImpl) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	product := &models.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	if err := s.repo.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductServiceImpl) GetProductByID(id int) (*models.Product, error) {
	return s.repo.FindProductByID(id)
}

func (s *ProductServiceImpl) UpdateProduct(id int, req *models.UpdateProductRequest) (*models.Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	product := &models.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	if err := s.repo.UpdateProduct(id, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductServiceImpl) PatchProduct(id int, req *models.PatchProductRequest) (*models.Product, error) {
	return s.repo.PatchProduct(id, req.Name, req.Price, req.Stock)
}

func (s *ProductServiceImpl) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
