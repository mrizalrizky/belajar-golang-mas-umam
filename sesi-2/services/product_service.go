package services

import (
	"mrizalrizky/sesi-2/models"
	"mrizalrizky/sesi-2/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) CreateProduct(product *models.Product) (error) {
	return s.repo.Create(product)
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) UpdateProductByID(product *models.Product) (error) {
	return s.repo.UpdateByID(product)
}

func (s *ProductService) DeleteProductByID(id int) error {
	return s.repo.DeleteByID(id)
}