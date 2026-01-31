package models

import "errors"

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
	Category_ID int `json:"category_id"`
	Category_Name string `json:"category_name"`
}

type CreateProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
	Category_ID int `json:"category_id"`
}

type UpdateProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
	Category_ID int `json:"category_id"`
}

type PatchProductRequest struct {
	Name  *string `json:"name,omitempty"`
	Price *int    `json:"price,omitempty"`
	Stock *int    `json:"stock,omitempty"`
	Category_ID *int `json:"category_id,omitempty"`
}

func (p *CreateProductRequest) Validate() error {
	if p.Name == "" || p.Price == 0 || p.Stock == 0 || p.Category_ID == 0 {
		return errors.New("Name, price, stock and category_id are required")
	}
	return nil
}

func (p *UpdateProductRequest) Validate() error {
	if p.Name == "" || p.Price == 0 || p.Stock == 0 || p.Category_ID == 0{
		return errors.New("Name, price, stock and category_id are required")
	}
	return nil
}
