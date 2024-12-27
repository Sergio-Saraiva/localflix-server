package repositories

import (
	"database/sql"
	"fmt"
	"localflix-server/src/models"
)

type CategoriesRepository struct {
	db *sql.DB
}

func NewCategoriesRepository(db *sql.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (c *CategoriesRepository) CreateCategory(name string) (*models.Category, error) {
	result, err := c.db.Exec("INSERT INTO categories (name) VALUES (?)", name)
	if err != nil {
		fmt.Printf("error inserting category: %v\n", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("error getting last insert id: %v\n", err)
		return nil, err
	}

	return &models.Category{
		ID:   int(id),
		Name: name,
	}, nil
}

func (c *CategoriesRepository) GetCategory(id int) (*models.Category, error) {
	row := c.db.QueryRow("SELECT * FROM categories WHERE id = ?", id)
	var category models.Category
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		fmt.Printf("error getting category: %v\n", err)
		return nil, err
	}

	return &category, nil
}

func (c *CategoriesRepository) ListCategories() []*models.Category {
	rows, err := c.db.Query("SELECT * FROM categories")
	if err != nil {
		fmt.Printf("error listing categories: %v\n", err)
		return nil
	}

	var categories []*models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			fmt.Printf("error scanning category: %v\n", err)
			return nil
		}

		categories = append(categories, &category)
	}

	return categories
}

func (c *CategoriesRepository) DeleteCategory(id int) error {
	_, err := c.db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		fmt.Printf("error deleting category: %v\n", err)
		return err
	}

	return nil
}

func (c *CategoriesRepository) UpdateCategory(id int, name string) (*models.Category, error) {
	_, err := c.db.Exec("UPDATE categories SET name = ? WHERE id = ?", name, id)
	if err != nil {
		fmt.Printf("error updating category: %v\n", err)
		return nil, err
	}

	return &models.Category{
		ID:   id,
		Name: name,
	}, nil
}

func (c *CategoriesRepository) GetCategoryByName(name string) (*models.Category, error) {
	row := c.db.QueryRow("SELECT * FROM categories WHERE name = ?", name)
	var category models.Category
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		fmt.Printf("error getting category: %v\n", err)
		return nil, err
	}

	return &category, nil
}
