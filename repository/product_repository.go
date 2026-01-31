package repository

import "gokasir-api/models"

type ProductRepository interface {
	FindAllProduct() ([]models.Product, error)
	CreateProduct(req *models.Product) error
	FindProductByID(id int) (*models.Product, error)
	UpdateProduct(id int, req *models.Product) error
	PatchProduct(id int, name *string, price, stock, category_id *int) (*models.Product, error)
	DeleteProduct(id int) error
	ExistID(id int) (bool, error)
}
