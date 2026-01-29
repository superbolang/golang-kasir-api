package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"gokasir-api/models"
	"strings"
)

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) ExistID(id int) (bool, error) {
	var exist bool
	err := r.db.QueryRow("SELECT COUNT(*) FROM product WHERE id = $1", id).Scan(&exist)
	return exist, err
}

func (r *ProductRepositoryImpl) FindAllProduct() ([]models.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock FROM product ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepositoryImpl) CreateProduct(req *models.Product) error {
	err := r.db.QueryRow("INSERT INTO product(name, price, stock) VALUES($1, $2, $3) RETURNING id", req.Name, req.Price, req.Stock).Scan(&req.ID)
	return err
}

func (r *ProductRepositoryImpl) FindProductByID(id int) (*models.Product, error) {
	exist, err := r.ExistID(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("Product ID not found")
	}
	var product models.Product
	if err := r.db.QueryRow("SELECT id, name, price, stock FROM product WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) UpdateProduct(id int, req *models.Product) error {
	exist, err := r.ExistID(id)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Product ID not found")
	}
	result, err := r.db.Exec("UPDATE product SET name = $1, price = $2, stock = $3 WHERE id = $4", req.Name, req.Price, req.Stock, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		if rows == 0 {
			return errors.New("Product not found")
		}
		return err
	}
	return nil
}

func (r *ProductRepositoryImpl) PatchProduct(id int, name *string, price, stock *int) (*models.Product, error) {
	exist, err := r.ExistID(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("Product ID not found")
	}

	// Dynamic query
	query := "UPDATE product SET "
	var args []any
	var updates []string
	argCount := 1

	if name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, name)
		argCount++
	}
	if price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argCount))
		args = append(args, price)
		argCount++
	}
	if stock != nil {
		updates = append(updates, fmt.Sprintf("stock = $%d", argCount))
		args = append(args, stock)
		argCount++
	}
	if len(updates) == 0 {
		return r.FindProductByID(id)
	}
	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return r.FindProductByID(id)
}

func (r *ProductRepositoryImpl) DeleteProduct(id int) error {
	exist, err := r.ExistID(id)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Product ID not found")
	}
	result, err := r.db.Exec("DELETE FROM product WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		if rows == 0 {
			return errors.New("Product not found")
		}
		return err
	}
	return nil
}
