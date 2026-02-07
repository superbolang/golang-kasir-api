package handler

import (
	"encoding/json"
	"gokasir-api/models"
	"gokasir-api/service"
	"io"
	"net/http"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/checkout" {
		switch r.Method {
		case http.MethodPost:
			h.handleCheckout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	if r.URL.Path == "/api/v1/report" {
		switch r.Method {
		case http.MethodGet:
			h.handleGetTransaction(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	if r.URL.Path == "/api/v1/report/today" {
		switch r.Method {
		case http.MethodGet:
			h.handleGetTodaysTransaction(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	http.NotFound(w, r)
}

func (h *TransactionHandler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error handling reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var req models.CheckoutRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error handling unmarshal", http.StatusInternalServerError)
		return
	}
	checkout, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(checkout)
}

func (h *TransactionHandler) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start_date")
	end := r.URL.Query().Get("end_date")
	if start != "" && end != "" {
		rangeTransaction, err := h.service.RangeTransaction(start, end)
		if err != nil {
			http.Error(w, "Error handling get range transaction", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rangeTransaction)
	} else {
		transactions, err := h.service.GetAllTransaction()
		if err != nil {
			http.Error(w, "Error handling get all transactions", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transactions)
	}
}

func (h *TransactionHandler) handleGetTodaysTransaction(w http.ResponseWriter, r *http.Request) {
	todaysTransaction, err := h.service.TodaysTransaction()
	if err != nil {
		http.Error(w, "Error handling get today's transaction", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todaysTransaction)
}
