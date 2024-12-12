package dal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"hot-coffee/models"
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
	queryOrder := `
	insert into orders (CustomerName, CreatedAt, Status) values
	($1, $2, $3)
	`
	_, err := repo.db.Query(queryOrder, order.CustomerName, order.CreatedAt, order.Status)

	for _, v := range order.Items {
		queryOrderItems := `
		insert into order_items (ProductID, Quantity, OrderId) values
		($1, $2, $3)
		`
		repo.db.Query(queryOrderItems, v.ProductID, v.Quantity, order.ID)
	}

	return err
}

func (repo *OrderRepository) SaveUpdatedOrder(updatedOrder models.Order, OrderID string) error {
	queryCheckID := `
	select max(ID) from orders
	`
	var tempID string
	row, _ := repo.db.Query(queryCheckID)
	row.Scan(&tempID)
	temp1, _ := strconv.Atoi(OrderID)
	temp2, _ := strconv.Atoi(tempID)
	if temp1 > temp2 {
		return errors.New("the requested order does not exists")
	}
	queryCheckStatus := `
	select Status from orders where ID = $1
	`
	row, _ = repo.db.Query(queryCheckStatus, OrderID)
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
	_, err := repo.db.Query(queryUpdateOrder, updatedOrder.CustomerName, updatedOrder.Status, updatedOrder.ID)
	for _, v := range updatedOrder.Items {
		queryUpdateOrderItems := `
		update order_items set ProductID = $1, Quantity = $2 where OrderID = $3
		`
		repo.db.Query(queryUpdateOrderItems, v.ProductID, v.Quantity, updatedOrder.ID)
	}
	return err
}

func (repo *OrderRepository) DeleteOrder(OrderID string) error {
	queryCheckID := `
	select max(ID) from orders
	`
	var tempID string
	row, _ := repo.db.Query(queryCheckID)
	row.Scan(&tempID)
	temp1, _ := strconv.Atoi(OrderID)
	temp2, _ := strconv.Atoi(tempID)
	if temp1 > temp2 {
		return errors.New("the requested order does not exists")
	}
	queryDeleteOrder := `
	delete from orders
	where ID = $1
	`
	repo.db.Query(queryDeleteOrder, OrderID)
	queryDeleteOrderItems := `
	delete from order_items
	where ID = $1
	`
	repo.db.Query(queryDeleteOrderItems, OrderID)
}

func (repo *OrderRepository) GetID() (int, error) {
	configPath := filepath.Join(filepath.Dir(repo.path), "config.json")

	ConfigContent, err := os.ReadFile(configPath)
	if err != nil {
		return -1, err
	}

	var ID models.OrderID
	err = json.Unmarshal(ConfigContent, &ID)
	if err != nil {
		return -1, err
	}

	i := ID.ID
	ID.ID++
	NewContent, err := json.MarshalIndent(ID, "", "    ")
	if err != nil {
		// TODO
	}
	os.WriteFile(configPath, NewContent, os.ModePerm)
	return i, nil
}

func getOrders(db *sql.DB) ([]models.Order, error) {
	query := `
	 SELECT ID, CustomerName, Status, CreatedAt
	 FROM orders`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.Status, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		// Получение элементов заказа
		items, err := getOrderItems(db, order.ID)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения элементов заказа: %w", err)
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
