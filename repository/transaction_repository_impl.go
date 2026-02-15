package repository

import (
	"database/sql"
	"gokasir-api/models"
	"log"
	"time"
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

func (r *TransactionRepositoryImpl) FindAllTransaction() ([]models.TransactionDetail, error) {
	query := "SELECT t.id, t.transaction_id, t.product_id, p.name, t.quantity, t.sub_total FROM transaction_details t INNER JOIN product p ON t.product_id = p.id ORDER BY t.id"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error getting all transaction details: %v", err)
		return nil, err
	}
	defer rows.Close()
	var transactions []models.TransactionDetail
	for rows.Next() {
		var transaction models.TransactionDetail
		if err := rows.Scan(&transaction.ID, &transaction.TransactionID, &transaction.ProductID, &transaction.ProductName, &transaction.Quantity, &transaction.SubTotal); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepositoryImpl) TodaysTransaction() (*models.Report, error) {
	currentTime := time.Now().Format("2006-01-02")

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id, total_amount FROM transactions WHERE created_at::date = $1", currentTime)
	if err != nil {
		log.Printf("Error getting transaction: %v", err)
		return nil, err
	}
	defer rows.Close()
	totalEarning := 0
	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.TotalAmount); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
		totalEarning += transaction.TotalAmount
	}
	var productSold []models.ProductSold
	for i := range transactions {
		rows, err := tx.Query("SELECT p.name, t.quantity FROM transaction_details t INNER JOIN product p ON t.product_id = p.id WHERE t.transaction_id = $1", transactions[i].ID)
		if err != nil {
			log.Printf("Error getting transaction details: %v", err)
			return nil, err
		}
		defer rows.Close()
		var product models.ProductSold
		for rows.Next() {
			if err := rows.Scan(&product.ProductName, &product.ProductQty); err != nil {
				return nil, err
			}
			productSold = append(productSold, product)
		}
	}

	highSelling := 0
	var mostSold string
	for _, v := range productSold {
		if v.ProductQty > highSelling {
			highSelling = v.ProductQty
			mostSold = v.ProductName
		}
	}
	details := models.ProductSold{
		ProductName: mostSold,
		ProductQty:  highSelling,
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Report{
		TotalRevenue:     totalEarning,
		TotalTransaction: len(transactions),
		HighestSelling:   details,
	}, err
}

func (r *TransactionRepositoryImpl) RangeTransaction(start, end string) (*models.Report, error) {
	start_date := start + " 00:00:00"
	end_date := end + " 23:59:50"

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id, total_amount FROM transactions WHERE created_at >= $1 AND created_at < $2", start_date, end_date)
	if err != nil {
		log.Printf("Error getting transaction: %v", err)
		return nil, err
	}
	defer rows.Close()
	totalEarning := 0
	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.TotalAmount); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
		totalEarning += transaction.TotalAmount
	}
	var productSold []models.ProductSold
	for i := range transactions {
		rows, err := tx.Query("SELECT p.name, t.quantity FROM transaction_details t INNER JOIN product p ON t.product_id = p.id WHERE t.transaction_id = $1", transactions[i].ID)
		if err != nil {
			log.Printf("Error getting transaction details: %v", err)
			return nil, err
		}
		defer rows.Close()
		var product models.ProductSold
		for rows.Next() {
			if err := rows.Scan(&product.ProductName, &product.ProductQty); err != nil {
				return nil, err
			}
			productSold = append(productSold, product)
		}
	}
	log.Println("All product sold: ", productSold)

	highSelling := 0
	var mostSold string
	for _, v := range productSold {
		if v.ProductQty > highSelling {
			highSelling = v.ProductQty
			mostSold = v.ProductName
		}
	}
	details := models.ProductSold{
		ProductName: mostSold,
		ProductQty:  highSelling,
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Report{
		TotalRevenue:     totalEarning,
		TotalTransaction: len(transactions),
		HighestSelling:   details,
	}, err
}
