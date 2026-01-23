package handler

import (
	"encoding/json"
	"gokasir-api/models"
	"io"
	"net/http"
)

var kategori = []models.Kategori{
	{ID: 1, Nama: "Minuman Sachet", Deskripsi: "Minuman instan dengan kemasan plastik"},
	{ID: 2, Nama: "Snack", Deskripsi: "Makanan ringan dengan kemasan plastik"},
	{ID: 3, Nama: "Minuman Botol", Deskripsi: "Minuman dengan kemasan berbentuk botol maupun gelas"},
}

// Helper
func FindKategoriByID(id int) bool {
	for _, val := range kategori {
		if val.ID == id {
			return true
		}
	}
	return false
}

func FindKategoriIndex(id int) int {
	for index, val := range kategori {
		if val.ID == id {
			return index
		}
	}
	return -1
}

// kategori handler
func GetAllKategori(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kategori)
}

func CreateKategori(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var newKategori models.Kategori
	if err := json.Unmarshal(body, &newKategori); err != nil {
		http.Error(w, "Error unmarshal JSON", http.StatusInternalServerError)
		return
	}
	if newKategori.Nama == "" || newKategori.Deskripsi == "" {
		http.Error(w, "Nama dan deskripsi tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// Find max ID
	maxID := 0
	for _, val := range kategori {
		if val.ID > maxID {
			maxID = val.ID
		}
	}
	newKategori.ID = maxID + 1
	kategori = append(kategori, newKategori)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newKategori)
}

func GetKategoriByID(w http.ResponseWriter, r *http.Request, id int) {
	// Find kategori
	exist := FindKategoriByID(id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, val := range kategori {
		if val.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(val)
		}
	}
}

func UpdateKategori(w http.ResponseWriter, r *http.Request, id int) {
	// Find kategori
	exist := FindKategoriByID(id)
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
	var updateKategori models.Kategori
	if err := json.Unmarshal(body, &updateKategori); err != nil {
		http.Error(w, "Error unmarshal JSON", http.StatusInternalServerError)
		return
	}
	if updateKategori.Nama == "" || updateKategori.Deskripsi == "" {
		http.Error(w, "Nama, harga dan stok tidak boleh kosong", http.StatusBadRequest)
		return
	}

	updateKategori.ID = id
	for index, val := range kategori {
		if val.ID == id {
			kategori[index] = updateKategori
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(kategori[index])
		}
	}
}

func PatchKategori(w http.ResponseWriter, r *http.Request, id int) {
	// Find kategori
	exist := FindKategoriByID(id)
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
	var patch models.KategoriPatch
	if err := json.Unmarshal(body, &patch); err != nil {
		http.Error(w, "Error unmarshal JSON", http.StatusInternalServerError)
		return
	}
	for index, val := range kategori {
		if val.ID == id {
			if patch.Nama != nil {
				kategori[index].Nama = *patch.Nama
			}
			if patch.Deskripsi != nil {
				kategori[index].Deskripsi = *patch.Deskripsi
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(kategori[index])
		}
	}
}

func DeleteKategori(w http.ResponseWriter, r *http.Request, id int) {
	// Find kategori
	exist := FindKategoriByID(id)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find index
	index := FindKategoriIndex(id)
	kategori = append(kategori[:index], kategori[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}
