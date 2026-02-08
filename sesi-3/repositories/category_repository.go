package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"mrizalrizky/sesi-3/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
	// _, err := r.db.Exec("INSERT INTO categories (name, description) VALUES (?, ?)", category.Name, category.Description)
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var category models.Category
	err := r.db.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("Category not found")
	}
	return &category, nil
}

func (r *CategoryRepository) UpdateByID(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := r.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	fmt.Println("ROWS", rows, err)
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return errors.New("Product not found")
	}

	return err
}

func (r *CategoryRepository) DeleteByID(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
	}

	return err
}