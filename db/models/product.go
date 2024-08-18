package models

import (
	"database/sql"
	"errors"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Report      string // json containing key: Query, Report
	Price       float64
	Cost        float64
	BusinessID  int
}

type DashboardProduct struct {
	Product
	InEnvironments int
}

type productModel struct{}

var ProductModel *productModel

func (*productModel) Create(p Product) error {
	query := `
		INSERT INTO Products (name, description, price, cost, business_id)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.GetDB().Exec(query, p.Name, p.Description, p.Price, p.Cost, p.BusinessID)
	return err
}

var ErrProductNotFound error = errors.New("product not found")

func (*productModel) GetByID(id int) (*Product, error) {
	query := `
		SELECT id, name, description, price, cost, research, business_id
		FROM Products
		WHERE id = ?
	`
	row := db.GetDB().QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Cost, &product.Report, &product.BusinessID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (*productModel) GetAll() ([]Product, error) {
	query := `
		SELECT id, name, description, price, cost, business_id
		FROM Products
	`
	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Cost, &product.BusinessID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (*productModel) GetAllByBusinessID(id int) ([]Product, error) {
	query := `
        SELECT id, name, description, price, cost, business_id
        FROM Products
        WHERE business_id = ?
    `
	rows, err := db.GetDB().Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Cost, &product.BusinessID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (*productModel) GetDashboardData(businessId int) ([]DashboardProduct, error) {
	query := `
        SELECT 
            p.id, p.name, p.description, p.research, p.price, p.cost, p.business_id,
            COUNT(ep.environment_id) as in_environments
        FROM 
            Products p
        LEFT JOIN 
            EnvironmentProducts ep ON p.id = ep.product_id
		WHERE 
			p.business_id = ?
		GROUP BY 
            p.id
        ORDER BY 
            in_environments DESC
        LIMIT 4;
    `

	rows, err := db.GetDB().Query(query, businessId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []DashboardProduct

	for rows.Next() {
		var dp DashboardProduct
		if err := rows.Scan(&dp.ID, &dp.Name, &dp.Description, &dp.Report, &dp.Price, &dp.Cost, &dp.BusinessID, &dp.InEnvironments); err != nil {
			return nil, err
		}
		products = append(products, dp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// p itself already contains the id for the product to be changed
func (*productModel) Update(p Product) error {
	query := `
		UPDATE Products
		SET name = ?, description = ?, price = ?, cost = ?, research = ?, business_id = ?
		WHERE id = ?
	`
	_, err := db.GetDB().Exec(query, p.Name, p.Description, p.Price, p.Cost, p.Report, p.BusinessID, p.ID)
	return err
}

func (*productModel) Delete(id int) error {
	query := `
		DELETE FROM Products
		WHERE id = ?
	`
	_, err := db.GetDB().Exec(query, id)
	return err
}
