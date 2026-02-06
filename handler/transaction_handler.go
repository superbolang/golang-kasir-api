package handler

import (
	"encoding/json"
	"gokasir-api/models"
	"gokasir-api/service"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/checkout")
	if r.URL.Path == "/api/v1/checkout" || r.URL.Path == "/api/v1/checkout/" {
		switch r.Method {
		case http.MethodPost:
			h.handleCheckout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	if strings.HasPrefix(path, "/") {
		idStr := strings.TrimPrefix(path, "/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusInternalServerError)
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.handleGet(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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
	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) handleGet(w http.ResponseWriter, r *http.Request, id int) {}
