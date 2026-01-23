package handler

import (
	"encoding/json"
	"gokasir-api/models"
	"io"
	"net/http"
)

var produk = []models.Produk{
	{ID: 1, Nama: "Kopi kapal api", Harga: 1500, Stok: 10},
	{ID: 2, Nama: "Indomie goreng", Harga: 3500, Stok: 20},
	{ID: 3, Nama: "Kacang garuda", Harga: 15000, Stok: 40},
}

// Helper
func FindProdukByID(id int) bool {
	for _, val := range produk {
		if val.ID == id {
			return true
		}
	}
	return false
}

func FindIndex(id int) int {
	for index, val := range produk {
		if val.ID == id {
			return index
		}
	}
	return -1
}

// Produk handler
func GetAllProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(produk)
}

func CreateProduk(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var newProduk models.Produk
	if err := json.Unmarshal(body, &newProduk); err != nil {
		http.Error(w, "Error unmarshal JSON", http.StatusInternalServerError)
		return
	}
	if newProduk.Nama == "" || newProduk.Harga == 0 || newProduk.Stok == 0 {
		http.Error(w, "Nama, harga dan stok tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// Find max ID
	maxID := 0
	for _, val := range produk {
		if val.ID > maxID {
			maxID = val.ID
		}
	}
	newProduk.ID = maxID + 1
	produk = append(produk, newProduk)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduk)
}

func GetProdukByID(w http.ResponseWriter, r *http.Request, id int) {
	// Find produk
	exist := FindProdukByID(id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, val := range produk {
		if val.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(val)
		}
	}
}

func UpdateProduk(w http.ResponseWriter, r *http.Request, id int) {
	// Find produk
	exist := FindProdukByID(id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var updateProduk models.Produk
	if err := json.Unmarshal(body, &updateProduk); err != nil {
		http.Error(w, "Error unmarshal JSON", http.StatusInternalServerError)
		return
	}
	if updateProduk.Nama == "" || updateProduk.Harga == 0 || updateProduk.Stok == 0 {
		http.Error(w, "Nama, harga dan stok tidak boleh kosong", http.StatusBadRequest)
		return
	}

	updateProduk.ID = id
	for index, val := range produk {
		if val.ID == id {
			produk[index] = updateProduk
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(produk[index])
		}
	}
}

func PatchProduk(w http.ResponseWriter, r *http.Request, id int) {
	// Find produk
	exist := FindProdukByID(id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var patch models.ProdukPatch
	if err := json.Unmarshal(body, &patch); err != nil {
		http.Error(w, "Error unmarshal JSON", http.StatusInternalServerError)
		return
	}
	for index, val := range produk {
		if val.ID == id {
			if patch.Nama != nil {
				produk[index].Nama = *patch.Nama
			}
			if patch.Harga != nil {
				produk[index].Harga = *patch.Harga
			}
			if patch.Stok != nil {
				produk[index].Stok = *patch.Stok
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(produk[index])
		}
	}
}

func DeleteProduk(w http.ResponseWriter, r *http.Request, id int) {
	// Find produk
	exist := FindProdukByID(id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find index
	index := FindIndex(id)
	produk = append(produk[:index], produk[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}
