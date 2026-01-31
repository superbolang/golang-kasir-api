package handler

import (
	"encoding/json"
	"gokasir-api/models"
	"gokasir-api/service"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/product")
	if r.URL.Path == "/api/v1/product" || r.URL.Path == "/api/v1/product/" {
		switch r.Method {
		case http.MethodGet:
			h.handleGetAll(w, r)
		case http.MethodPost:
			h.handleCreate(w, r)
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
			h.handleGetByID(w, r, id)
		case http.MethodPut:
			h.handleUpdate(w, r, id)
		case http.MethodPatch:
			h.handlePatch(w, r, id)
		case http.MethodDelete:
			h.handleDelete(w, r, id)
		}
		return
	}
	http.NotFound(w, r)
}

func (h *ProductHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProduct()
	if err != nil {
		log.Printf("Error handling get product: %v", err)
		http.Error(w, "Error handling get product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&products)
}

func (h *ProductHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error handling reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var req models.CreateProductRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error handling unmarshal", http.StatusInternalServerError)
		return
	}
	product, err := h.service.CreateProduct(&req)
	if err != nil {
		log.Printf("Error handling creating product: %v", err)
		http.Error(w, "Error handling creating product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&product)
}

func (h *ProductHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id int) {
	product, err := h.service.GetProductByID(id)
	if err != nil {
		log.Printf("Error handling getting single product: %v", err)
		http.Error(w, "Error handling getting single product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&product)
}

func (h *ProductHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error handling reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var req models.UpdateProductRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error handling unmarshal", http.StatusInternalServerError)
		return
	}
	product, err := h.service.UpdateProduct(id, &req)
	if err != nil {
		log.Printf("Error handling updating product: %v", err)
		http.Error(w, "Error handling updating product", http.StatusInternalServerError)
		return
	}
	product.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&product)
}

func (h *ProductHandler) handlePatch(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error handling reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var req models.PatchProductRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error handling unmarshal", http.StatusInternalServerError)
		return
	}
	product, err := h.service.PatchProduct(id, &req)
	if err != nil {
		log.Printf("Error handling patching product: %v", err)
		http.Error(w, "Error handling patching product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&product)
}

func (h *ProductHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.DeleteProduct(id); err != nil {
		log.Printf("Error handling deleting product: %v", err)
		http.Error(w, "Error handling deleting product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
