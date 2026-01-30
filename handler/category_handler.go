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

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/category")
	if r.URL.Path == "/api/v1/category" || r.URL.Path == "/api/v1/category/" {
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
			http.Error(w, "Invalid ID", http.StatusBadRequest)
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
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	http.NotFound(w, r)
}

func (h *CategoryHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	category, err := h.service.GetAllCategory()
	if err != nil {
		http.Error(w, "Failed to get all category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&category)
}

func (h *CategoryHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var req models.CreateCategoryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error unmarshal", http.StatusInternalServerError)
		return
	}
	category, err := h.service.CreateCategory(&req)
	if err != nil {
		http.Error(w, "Error creating category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&category)
}

func (h *CategoryHandler) handleGetByID(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "Error getting category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&category)
}

func (h *CategoryHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var req models.UpdateCategoryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error unmarshal", http.StatusInternalServerError)
		return
	}
	category, err := h.service.UpdateCategory(id, &req)
	if err != nil {
		http.Error(w, "Error updating category", http.StatusInternalServerError)
		return
	}
	category.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&category)
}

func (h *CategoryHandler) handlePatch(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var req models.PatchCategoryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error unmarshal", http.StatusInternalServerError)
		return
	}
	category, err := h.service.PatchCategory(id, &req)
	if err != nil {
		http.Error(w, "Error patching category", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&category)
}

func (h *CategoryHandler) handleDelete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.DeleteCategory(id); err != nil {
		http.Error(w, "Error deleting category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
