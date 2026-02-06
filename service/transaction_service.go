package service

import "gokasir-api/models"

type TransactionService interface {
	Checkout(item []models.CheckoutItem) (*models.Transaction, error)
}
