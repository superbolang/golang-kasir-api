package repository

import "gokasir-api/models"

type CategoryRepository interface {
	FindAllCategory() ([]models.Category, error)
	CreateCategory(req *models.Category) error
	FindCategoryByID(id int) (*models.Category, error)
	UpdateCategory(id int, req *models.Category) error
	PatchCategory(id int, name, description *string) (*models.Category, error)
	DeleteCategory(id int) error
	ExistCategoryID(id int) (bool, error)
}
