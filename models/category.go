package models

import "errors"

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Description string    `json:"description"`
}

type CreateCategoryRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Description string    `json:"description"`
}

type UpdateCategoryRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Description string    `json:"description"`
}

type PatchCategoryRequest struct {
	ID    *int    `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
}

func (p *CreateCategoryRequest) Validate() error {
	if p.Name == "" || p.Description == "" {
		return errors.New("Name and description are required")
	}
	return nil
}

func (p *UpdateCategoryRequest) Validate() error {
	if p.Name == "" || p.Description == "" {
		return errors.New("Name and description are required")
	}
	return nil
}
