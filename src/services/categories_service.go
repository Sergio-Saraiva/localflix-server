package services

import (
	"context"
	"database/sql"
	"fmt"
	"localflix-server/src/models"
	"localflix-server/src/repositories"
)

type CategoriesService struct {
	ctx                  context.Context
	categoriesRepository *repositories.CategoriesRepository
}

// NewApp creates a new App application struct
func NewCategoriesService(ctx context.Context, db *sql.DB) *CategoriesService {
	return &CategoriesService{
		ctx:                  ctx,
		categoriesRepository: repositories.NewCategoriesRepository(db),
	}
}

func (c *CategoriesService) CreateCategory(name string) (*models.Category, error) {
	category, err := c.categoriesRepository.CreateCategory(name)
	if err != nil {
		fmt.Printf("error creating category: %v\n", err)
		return nil, err
	}
	fmt.Printf("category created: %v\n", category)
	return category, nil
}

func (c *CategoriesService) GetCategory(id int) (*models.Category, error) {
	category, err := c.categoriesRepository.GetCategory(id)
	if err != nil {
		fmt.Printf("error getting category by id: %v\n", err)
		return nil, err
	}

	return category, nil
}

func (c *CategoriesService) GetCategoryByName(name string) (*models.Category, error) {
	category, err := c.categoriesRepository.GetCategoryByName(name)
	if err != nil {
		fmt.Printf("error getting category by name: %v\n", err)
		return nil, err
	}

	return category, nil
}

func (c *CategoriesService) ListCategories() []models.Category {
	categories := c.categoriesRepository.ListCategories()
	result := make([]models.Category, len(categories))
	for i, category := range categories {
		fmt.Println("category: ", category)
		result[i] = *category
	}
	return result
}

func (c *CategoriesService) DeleteCategory(id int) error {
	err := c.categoriesRepository.DeleteCategory(id)
	if err != nil {
		fmt.Printf("error deleting category: %v\n", err)
		return err
	}

	return nil
}

func (c *CategoriesService) UpdateCategory(id int, name string) (*models.Category, error) {
	category, err := c.categoriesRepository.UpdateCategory(id, name)
	if err != nil {
		fmt.Printf("error updating category: %v\n", err)
		return nil, err
	}

	return category, nil
}
