package services

import (
	"mrizalrizky/sesi-2/models"
	"mrizalrizky/sesi-2/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) CreateCategory(category *models.Category) (error) {
	return s.repo.Create(category)
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) UpdateCategoryByID(category *models.Category) error {
	return s.repo.UpdateByID(category)
}

func (s *CategoryService) DeleteCategoryByID(id int) error {
	return s.repo.DeleteByID(id)
}