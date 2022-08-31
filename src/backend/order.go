package backend

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Total        int         `json:"total"`
	Status       string      `json:"status"`
	Items        []orderItem `json:"items"`
}

type orderItem struct {
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func getOrders(db *sql.DB) ([]order, error) {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	orders := []order{}

	for rows.Next() {
		var o order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Total, &o.Status); err != nil {
			return nil, err
		}

		err = o.getOrderItems(db)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}

func (o *order) getOrder(db *sql.DB) error {
	err := db.QueryRow("SELECT customerName, total, status FROM orders WHERE id = ?", o.ID).Scan(&o.CustomerName, &o.Total, &o.Status)
	if err != nil {
		return err
	}

	err = o.getOrderItems(db)
	if err != nil {
		return err
	}
	return nil
}

func (o *order) getOrderItems(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM order_items WHERE order_id = ?", o.ID)
	if err != nil {
		return err
	}

	defer rows.Close()
	orderItems := []orderItem{}

	for rows.Next() {
		var oi orderItem
		if err := rows.Scan(&oi.OrderID, &oi.ProductID, &oi.Quantity); err != nil {
			return err
		}
		orderItems = append(orderItems, oi)
	}

	o.Items = orderItems
	return nil
}

func (o *order) createOrder(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO orders(customerName, total, status) VALUES(?, ?, ?)", o.CustomerName, o.Total, o.Status)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	o.ID = int(id)

	return nil
}

func (oi *orderItem) createOrderItem(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO order_items(order_id, product_id, quantity) VALUES(?, ?, ?)", oi.OrderID, oi.ProductID, oi.Quantity)
	if err != nil {
		return err
	}

	return nil
}
