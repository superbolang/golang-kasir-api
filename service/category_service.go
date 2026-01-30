package service

import "gokasir-api/models"

type CategoryService interface {
	GetAllCategory() ([]models.Category, error)
	CreateCategory(req *models.CreateCategoryRequest) (*models.Category, error)
	GetCategoryByID(id int) (*models.Category, error)
	UpdateCategory(id int, req *models.UpdateCategoryRequest) (*models.Category, error)
	PatchCategory(id int, req *models.PatchCategoryRequest) (*models.Category, error)
	DeleteCategory(id int) error
}
