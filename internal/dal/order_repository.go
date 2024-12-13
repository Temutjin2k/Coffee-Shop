package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"hot-coffee/models"
	"log"
)

// OrderRepository implements OrderRepository using JSON files
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new FileOrderRepository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) GetAll() ([]models.Order, error) {
	// Чтение заказов
	orders, err := getOrders(repo.db)
	if err != nil {
		log.Fatalf("Ошибка получения заказов: %v", err)
	}

	// Вывод данных
	return orders, err
}

func (repo *OrderRepository) Add(order models.Order) error {
	// Insert into `orders` and retrieve the ID
	queryOrder := `
        INSERT INTO orders (CustomerName, CreatedAt, Status)
        VALUES ($1, $2, $3)
        RETURNING ID
    `
	var ID int
	repo.db.QueryRow(queryOrder, order.CustomerName, order.CreatedAt, order.Status).Scan(&ID)

	for _, v := range order.Items {
		queryOrderItems := `
		insert into order_items (ProductID, Quantity, OrderID) values
		($1, $2, $3)
		`

		repo.db.Exec(queryOrderItems, v.ProductID, v.Quantity, ID)
	}

	return nil
}

func (repo *OrderRepository) SaveUpdatedOrder(updatedOrder models.Order, OrderID string) error {
	queryCheckStatus := `
	select Status from orders where ID = $1
	`
	row, err := repo.db.Query(queryCheckStatus, OrderID)
	if err != nil {
		return err
	}
	var Status string
	row.Scan(&Status)
	if Status == "closed" {
		return errors.New("the requested order is already closed")
	}

	queryUpdateOrder := `
	update orders 
	set CustomerName = $1, Status = $2
	where ID = $3
	`
	_, err = repo.db.Query(queryUpdateOrder, updatedOrder.CustomerName, updatedOrder.Status, updatedOrder.ID)
	if err != nil {
		return err
	}
	for _, v := range updatedOrder.Items {
		queryUpdateOrderItems := `
		update order_items set ProductID = $1, Quantity = $2 where OrderID = $3
		`
		_, err = repo.db.Query(queryUpdateOrderItems, v.ProductID, v.Quantity, updatedOrder.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *OrderRepository) DeleteOrder(OrderID string) error {
	queryDeleteOrderItems := `
	delete from order_items
	where OrderID = $1
	`
	_, err := repo.db.Exec(queryDeleteOrderItems, OrderID)
	if err != nil {
		return err
	}
	queryDeleteOrder := `
	delete from orders
	where ID = $1
	`
	_, err = repo.db.Exec(queryDeleteOrder, OrderID)
	if err != nil {
		log.Println("qwe")
		return err
	}
	return nil
}

func getOrders(db *sql.DB) ([]models.Order, error) {
	query := `
	 SELECT ID, CustomerName, Status, CreatedAt
	 FROM orders`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}

		// Получение элементов заказа
		items, err := getOrderItems(db, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func getOrderItems(db *sql.DB, orderID string) ([]models.OrderItem, error) {
	query := `
	 SELECT ProductID, Quantity
	 FROM order_items
	 WHERE OrderID = $1`

	rows, err := db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для order_items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки в order_items: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
