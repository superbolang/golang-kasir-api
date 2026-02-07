package models

import (
	"time"
)

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	SubTotal      int    `json:"sub_total"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type Report struct {
	TotalRevenue     int         `json:"total_revenue"`
	TotalTransaction int         `json:"total_transaction"`
	HighestSelling   ProductSold `json:"highest_selling"`
}

type ProductSold struct {
	ProductName string
	ProductQty  int
}
