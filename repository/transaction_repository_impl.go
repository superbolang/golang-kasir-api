package repository

import (
	"database/sql"
	"gokasir-api/models"
	"log"
)

type TransactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	// Crerate db transaction
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
		err := tx.QueryRow("SELECT name, price, stock FROM product WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Product %d not found", item.ProductID)
				return nil, err
			}
			return nil, err
		}
		if stock == 0 || stock < item.Quantity {
			log.Print("Stock quantity is less than required quantity")
			return nil, err
		}
		subTotal := productPrice * item.Quantity
		totalAmount += subTotal
		_, err = tx.Exec("UPDATE product SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			SubTotal:    subTotal,
		})
	}
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions(total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, sub_total) VALUES ($1,$2,$3,$4)", transactionID, details[i].ProductID, details[i].Quantity, details[i].SubTotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, err

}
