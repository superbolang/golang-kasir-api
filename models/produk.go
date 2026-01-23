package models

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type ProdukPatch struct {
	Nama  *string `json:"nama,omitempty"`
	Harga *int    `json:"harga,omitempty"`
	Stok  *int    `json:"stok,omitempty"`
}