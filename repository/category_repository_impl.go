package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"gokasir-api/models"
	"strings"
)

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) ExistCategoryID(id int) (bool, error) {
	var exist bool
	err := r.db.QueryRow("SELECT COUNT(*) FROM category WHERE id = $1", id).Scan(&exist)
	return exist, err
}

func (r *CategoryRepositoryImpl) FindAllCategory() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM category ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var category []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		category = append(category, cat)
	}
	return category, nil
}

func (r *CategoryRepositoryImpl) CreateCategory(req *models.Category) error {
	err := r.db.QueryRow("INSERT INTO category(name, description) VALUES($1, $2) RETURNING id", req.Name, req.Description).Scan(&req.ID)
	return err
}

func (r *CategoryRepositoryImpl) FindCategoryByID(id int) (*models.Category, error) {
	exist, err := r.ExistCategoryID(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("Category ID not found")
	}
	var category models.Category
	if err := r.db.QueryRow("SELECT id, name, description FROM category WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepositoryImpl) UpdateCategory(id int, req *models.Category) error {
	exist, err := r.ExistCategoryID(id)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Category ID not found")
	}
	result, err := r.db.Exec("UPDATE category SET name = $1, description = $2 WHERE id = $3", req.Name, req.Description, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		if rows == 0 {
			return errors.New("Category not found")
		}
		return err
	}
	return nil
}

func (r *CategoryRepositoryImpl) PatchCategory(id int, name, description *string) (*models.Category, error) {
	exist, err := r.ExistCategoryID(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("Category ID not found")
	}

	// Dynamic query
	query := "UPDATE category SET "
	var args []any
	var updates []string
	argCount := 1

	if name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, name)
		argCount++
	}
	if description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		args = append(args, description)
		argCount++
	}
	if len(updates) == 0 {
		return r.FindCategoryByID(id)
	}
	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return r.FindCategoryByID(id)
}

func (r *CategoryRepositoryImpl) DeleteCategory(id int) error {
	exist, err := r.ExistCategoryID(id)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Category ID not found")
	}
	result, err := r.db.Exec("DELETE FROM category WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		if rows == 0 {
			return errors.New("Category not found")
		}
		return err
	}
	return nil
}
