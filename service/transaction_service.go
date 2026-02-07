package service

import "gokasir-api/models"

type TransactionService interface {
	Checkout(item []models.CheckoutItem) (*models.Transaction, error)
	GetAllTransaction() ([]models.TransactionDetail, error)
	TodaysTransaction() (*models.Report, error)
	RangeTransaction(start, end string) (*models.Report, error)
}
