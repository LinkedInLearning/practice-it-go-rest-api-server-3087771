package backend

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type product struct {
	ID          int    `json:"id"`
	ProductCode string `json:"productCode"`
	Name        string `json:"name"`
	Inventory   int    `json:"inventory"`
	Price       int    `json:"price"`
	Status      string `json:"status"`
}

func getProducts(db *sql.DB) ([]product, error) {
	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.ProductCode, &p.Name, &p.Inventory, &p.Price, &p.Status); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, err
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT productCode, name, inventory, price, status FROM products WHERE id = ?", p.ID).Scan(&p.ProductCode, &p.Name, &p.Inventory, &p.Price, &p.Status)
}

func (p *product) createProduct(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO products(productCode, name, inventory, price, status) VALUES(?, ?, ?, ?, ?)", p.ProductCode, p.Name, p.Inventory, p.Price, p.Status)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)

	return nil
}
