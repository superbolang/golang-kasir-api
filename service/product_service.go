package service

import "gokasir-api/models"

type ProductService interface {
	GetAllProduct() ([]models.Product, error)
	CreateProduct(req *models.CreateProductRequest) (*models.Product, error)
	GetProductByID(id int) (*models.Product, error)
	UpdateProduct(id int, req *models.UpdateProductRequest) (*models.Product, error)
	PatchProduct(id int, req *models.PatchProductRequest) (*models.Product, error)
	DeleteProduct(id int) error
}
