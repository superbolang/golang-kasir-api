package repository

import "gokasir-api/models"

type TransactionRepository interface {
	CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error)
	FindAllTransaction() ([]models.TransactionDetail, error)
	TodaysTransaction() (*models.Report, error)
	RangeTransaction(start, end string) (*models.Report, error)
}
