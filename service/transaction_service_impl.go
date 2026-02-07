package service

import (
	"gokasir-api/models"
	"gokasir-api/repository"
)

type TransactionServiceImpl struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &TransactionServiceImpl{repo: repo}
}

func (s *TransactionServiceImpl) Checkout(item []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(item)
}

func (s *TransactionServiceImpl) GetAllTransaction() ([]models.TransactionDetail, error) {
	return s.repo.FindAllTransaction()
}

func (s *TransactionServiceImpl) TodaysTransaction() (*models.Report, error) {
	return s.repo.TodaysTransaction()
}
