package main

import (
	"fmt"
	"gokasir-api/handler"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// Produk handler
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetAllProduk(w, r)
		case http.MethodPost:
			handler.CreateProduk(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/api/produk/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handler.GetProdukByID(w, r, id)
		case http.MethodPut:
			handler.UpdateProduk(w, r, id)
		case http.MethodPatch:
			handler.PatchProduk(w, r, id)
		case http.MethodDelete:
			handler.DeleteProduk(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Kategori handler
	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetAllKategori(w, r)
		case http.MethodPost:
			handler.CreateKategori(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/api/kategori/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handler.GetKategoriByID(w, r, id)
		case http.MethodPut:
			handler.UpdateKategori(w, r, id)
		case http.MethodPatch:
			handler.PatchKategori(w, r, id)
		case http.MethodDelete:
			handler.DeleteKategori(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API server running OK"))
	})

	fmt.Println("Server running on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
