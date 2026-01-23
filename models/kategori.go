package models

type Kategori struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
}

type KategoriPatch struct {
	Nama  *string `json:"nama,omitempty"`
	Deskripsi *string    `json:"kategori,omitempty"`
}
