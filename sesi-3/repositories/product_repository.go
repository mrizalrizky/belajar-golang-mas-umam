package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"mrizalrizky/sesi-3/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, price, stock, c.id, c.name, c.description 
		FROM products p
		JOIN categories c on c.id = p.category_id 
		`

	var args []interface{}
	if name != "" {
		query += "WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		p.Category = &models.Category{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category.ID, &p.Category.Name, &p.Category.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
	SELECT p.id, p.name, price, stock, c.id, c.name, c.description
	FROM products AS p
	LEFT JOIN categories AS c ON c.id = p.category_id
	WHERE p.id = $1
	`

	var product models.Product
	product.Category = &models.Category{}
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category.ID, &product.Category.Name, &product.Category.Description)
	fmt.Println("PRODUCT", product)
	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}
	return &product, nil
}

func (r *ProductRepository) UpdateByID(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := r.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
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

	return nil
}

func (r *ProductRepository) DeleteByID(id int) error {
	query := "DELETE FROM products WHERE id = $1"
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