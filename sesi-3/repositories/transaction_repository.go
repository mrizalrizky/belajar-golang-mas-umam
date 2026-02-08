package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"mrizalrizky/sesi-3/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()


	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string
		
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		fmt.Println("ERR", err)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product ID %d not found", item.ProductID)
		}

		if stock < item.Quantity {
			return nil, errors.New("Stock not enough")
		}

		subTotal := productPrice * item.Quantity
		totalAmount += subTotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:     item.ProductID,
			ProductName:   productName,
			Quantity:      item.Quantity,
			Subtotal:      subTotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		var detailID int
		err = tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id", transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	fmt.Println("DETAILS", details)

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt: time.Now(),
		Details:     details,
	}, nil
}