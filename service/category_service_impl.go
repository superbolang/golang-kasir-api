package service

import (
	"gokasir-api/models"
	"gokasir-api/repository"
)

type CategoryServiceImpl struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{repo: repo}
}

func (s *CategoryServiceImpl) GetAllCategory() ([]models.Category, error) {
	return s.repo.FindAllCategory()
}

func (s *CategoryServiceImpl) CreateCategory(req *models.CreateCategoryRequest) (*models.Category, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryServiceImpl) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.FindCategoryByID(id)
}

func (s *CategoryServiceImpl) UpdateCategory(id int, req *models.UpdateCategoryRequest) (*models.Category, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.UpdateCategory(id, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryServiceImpl) PatchCategory(id int, req *models.PatchCategoryRequest) (*models.Category, error) {
	return s.repo.PatchCategory(id, req.Name, req.Description)
}

func (s *CategoryServiceImpl) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}
