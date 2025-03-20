package services

import (
	"braip/internal/models"
	"braip/internal/repository"
)

// GetProducts retorna todos os produtos
func GetProducts() ([]models.Product, error) {
	return repository.GetProducts()
}

// CreateProduct cria um novo produto
func CreateProduct(product models.Product) error {
	return repository.CreateProduct(product)
}

// GetProductByID retorna um produto pelo ID
func GetProductByID(id int) (*models.Product, error) {
	return repository.GetProductByID(id)
}

// UpdateProduct atualiza um produto
func UpdateProduct(id int, product models.Product) error {
	return repository.UpdateProduct(id, product)
}

// DeleteProduct remove um produto
func DeleteProduct(id int) error {
	return repository.DeleteProduct(id)
}

// SearchProductsByNameAndCategory busca produtos por nome e categoria
func SearchProductsByNameAndCategory(name, category string) ([]models.Product, error) {
	return repository.SearchProductsByNameAndCategory(name, category)
}

// SearchProductsByCategory busca produtos por categoria
func SearchProductsByCategory(category string) ([]models.Product, error) {
	return repository.SearchProductsByCategory(category)
}

// SearchProductsByImage busca produtos com ou sem imagem
func SearchProductsByImage(hasImage bool) ([]models.Product, error) {
	return repository.SearchProductsByImage(hasImage)
}