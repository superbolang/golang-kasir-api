package models

import "errors"

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type CreateProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type UpdateProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type PatchProductRequest struct {
	Name  *string `json:"name,omitempty"`
	Price *int    `json:"price,omitempty"`
	Stock *int    `json:"stock,omitempty"`
}

func (p *CreateProductRequest) Validate() error {
	if p.Name == "" || p.Price == 0 || p.Stock == 0 {
		return errors.New("Name, price and stock are required")
	}
	return nil
}

func (p *UpdateProductRequest) Validate() error {
	if p.Name == "" || p.Price == 0 || p.Stock == 0 {
		return errors.New("Name, price and stock are required")
	}
	return nil
}
