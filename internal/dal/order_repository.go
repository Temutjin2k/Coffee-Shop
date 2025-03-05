package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"hot-coffee/models"
)

type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new FileOrderRepository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) Add(order models.Order) (models.BatchOrderInfo, []models.BatchOrderInventoryUpdate, error) {
	processInfo := models.BatchOrderInfo{
		CustomerName: order.CustomerName,
		Status:       models.StatusOrderRejected,
	}
	inventoryInfo := []models.BatchOrderInventoryUpdate{}

	tx, err := repo.db.Begin()
	if err != nil {
		processInfo.Reason = "internal server error. Failed to start transaction."
		return processInfo, inventoryInfo, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Inserting order and getting ID
	queryOrder := `
        INSERT INTO orders (CustomerName)
        VALUES ($1)
        RETURNING ID
    `
	var ID int
	err = tx.QueryRow(queryOrder, order.CustomerName).Scan(&ID)
	if err != nil {
		processInfo.Reason = "internal server error. Failed to scan ID"
		return processInfo, inventoryInfo, err
	}
	processInfo.OrderID = ID

	// Inserting order items. in case when same product id is given, it check on conflict, if so it's just adding quantity for previus row.
	queryOrderItems := `
		INSERT INTO order_items (ProductID, Quantity, OrderID) VALUES
		($1, $2, $3)
		ON CONFLICT (OrderID, ProductID)
		DO UPDATE SET Quantity = order_items.Quantity + EXCLUDED.Quantity;
	`

	// Getting price of menu item
	queryGetPrice := `
		SELECT price FROM menu_items WHERE id = $1
	`

	for _, v := range order.Items {
		_, err = tx.Exec(queryOrderItems, v.ProductID, v.Quantity, ID)
		if err != nil {
			processInfo.Reason = "internal server error. " + err.Error()
			processInfo.Total = 0
			return processInfo, inventoryInfo, err
		}

		var price float64
		err = tx.QueryRow(queryGetPrice, v.ProductID).Scan(&price)
		if err != nil {
			processInfo.Reason = "internal server error." + err.Error()
			processInfo.Total = 0
			return processInfo, inventoryInfo, err
		}
		processInfo.Total += float64(v.Quantity) * price
	}

	// Reducing ingredients from inventory
	queryGetIngredients := `
	SELECT IngredientID, Quantity FROM menu_item_ingredients WHERE MenuID = $1
`

	queryUpdateInventory := `
	UPDATE inventory SET Quantity = Quantity - $1 WHERE IngredientID = $2 AND Quantity >= $1
`

	for _, v := range order.Items {
		// Загружаем все ингредиенты перед обработкой
		var ingredients []struct {
			IngredientID     int
			RequiredQuantity int
		}

		rows, err := tx.Query(queryGetIngredients, v.ProductID)
		if err != nil {
			processInfo.Reason = "internal server error. Failed to get ingredients."
			processInfo.Total = 0
			return processInfo, inventoryInfo, err
		}
		defer rows.Close()

		for rows.Next() {
			var ingredient struct {
				IngredientID     int
				RequiredQuantity int
			}
			if err := rows.Scan(&ingredient.IngredientID, &ingredient.RequiredQuantity); err != nil {
				processInfo.Reason = "internal server error. Failed to scan ingredient."
				processInfo.Total = 0
				return processInfo, inventoryInfo, err
			}
			ingredients = append(ingredients, ingredient)
		}

		// Проверяем и обновляем инвентарь
		for _, ing := range ingredients {
			totalRequired := ing.RequiredQuantity * v.Quantity

			var availableQuantity int
			var InvName string

			err = tx.QueryRow("SELECT quantity, name FROM inventory WHERE IngredientID = $1", ing.IngredientID).Scan(&availableQuantity, &InvName)
			if err != nil {
				processInfo.Reason = fmt.Sprintf("internal server error. Failed to check inventory. ID=%d", ing.IngredientID)
				processInfo.Total = 0
				return processInfo, inventoryInfo, err
			}

			if availableQuantity < totalRequired {
				processInfo.Reason = fmt.Sprintf("insufficient_inventory. IngredientID: %d. Required: %d, Available: %d", ing.IngredientID, totalRequired, availableQuantity)
				processInfo.Total = 0
				return processInfo, inventoryInfo, fmt.Errorf(processInfo.Reason)
			}

			// Обновляем инвентарь
			_, err = tx.Exec(queryUpdateInventory, totalRequired, ing.IngredientID)
			if err != nil {
				processInfo.Reason = "internal server error. Failed to update inventory."
				processInfo.Total = 0
				return processInfo, inventoryInfo, err
			}

			InvInfo := models.BatchOrderInventoryUpdate{
				IngredientID:  ing.IngredientID,
				Name:          InvName,
				Quantity_used: totalRequired,
				Remaining:     availableQuantity - totalRequired,
			}
			inventoryInfo = append(inventoryInfo, InvInfo)
		}
	}

	// Commiting transaction
	err = tx.Commit()
	if err != nil {
		processInfo.Reason = "Internal server error. Error commiting transaction."
		processInfo.Total = 0
		return processInfo, inventoryInfo, err
	}
	processInfo.Status = models.StatusOrderAccepted
	processInfo.Reason = "OK"
	return processInfo, inventoryInfo, nil
}

func (repo *OrderRepository) GetAll() ([]models.Order, error) {
	query := `
	 SELECT ID, CustomerName, Status, CreatedAt
	 FROM orders`

	rows, err := repo.db.Query(query)
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
		items, err := getOrderItems(repo.db, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *OrderRepository) SaveUpdatedOrder(updatedOrder models.Order, OrderID string) error {
	queryCheckStatus := `
	select Status from orders where ID = $1
	`
	var Status string
	// Use QueryRow instead of Query for single-row results
	err := repo.db.QueryRow(queryCheckStatus, OrderID).Scan(&Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no order found with the given ID")
		}
		return err // Return any other error
	}

	if Status == "closed" {
		return errors.New("the requested order is already closed")
	}
	queryUpdateOrder := `
	update orders 
	set CustomerName = $1
	where ID = $2
	`
	_, err = repo.db.Query(queryUpdateOrder, updatedOrder.CustomerName, OrderID)
	if err != nil {
		return err
	}
	for _, v := range updatedOrder.Items {
		queryUpdateOrderItems := `
		update order_items set ProductID = $1, Quantity = $2 where OrderID = $3
		`
		_, err = repo.db.Query(queryUpdateOrderItems, v.ProductID, v.Quantity, OrderID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *OrderRepository) DeleteOrder(OrderID int) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// checking if orders exist
	var orderExists bool
	queryCheckOrder := `SELECT EXISTS(SELECT 1 FROM orders WHERE ID = $1)`
	err = tx.QueryRow(queryCheckOrder, OrderID).Scan(&orderExists)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !orderExists {
		tx.Rollback()
		return fmt.Errorf("order with ID %d not found", OrderID)
	}

	queryDeleteOrderItems := `
	DELETE FROM order_items
	WHERE OrderID = $1
	`
	_, err = tx.Exec(queryDeleteOrderItems, OrderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete order items: %w", err)
	}

	queryDeleteOrder := `
	DELETE FROM orders
	WHERE ID = $1
	`
	_, err = tx.Exec(queryDeleteOrder, OrderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return tx.Commit()
}

func (repo *OrderRepository) CloseOrderRepo(id int) error {
	// Check current order status
	var status string
	queryCheckStatus := `
		SELECT status FROM orders WHERE ID = $1
	`
	err := repo.db.QueryRow(queryCheckStatus, id).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order with ID %d not found", id)
		}
		return err
	}

	if status == "closed" {
		return fmt.Errorf("order with ID %d is already closed", id)
	}

	// Update order status to "closed"
	queryToClose := `
		UPDATE orders SET status = 'closed' WHERE ID = $1
	`
	_, err = repo.db.Exec(queryToClose, id)
	if err != nil {
		return err
	}

	return nil
}

func getOrderItems(db *sql.DB, orderID int) ([]models.OrderItem, error) {
	query := `
	 SELECT ProductID, Quantity
	 FROM order_items
	 WHERE OrderID = $1`

	rows, err := db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed request for order_items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("error scanning row in order_items: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (repo *OrderRepository) GetNumberOfItems(startDate, endDate time.Time) (map[string]int, error) {
	// Query to fetch the number of items ordered in the given date range
	query := `
		SELECT
			m.Name,
			COALESCE(SUM(oi.Quantity), 0) AS total_quantity
		FROM
			menu_items m
		LEFT JOIN
			order_items oi ON m.ID = oi.ProductID
		LEFT JOIN
			orders o ON oi.OrderID = o.ID
		WHERE
			(o.CreatedAt BETWEEN $1 AND $2) AND o.Status = 'closed'
		GROUP BY
			m.Name
		ORDER BY
			total_quantity DESC;
	`

	// Execute the query with parameters
	rows, err := repo.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Create a map to store the results
	result := make(map[string]int)

	// Iterate over the rows and populate the result map
	for rows.Next() {
		var name string
		var quantity int
		if err := rows.Scan(&name, &quantity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result[name] = quantity
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating rows: %v", err)
	}

	return result, nil
}

func (repo *OrderRepository) GetEarliestDate() string {
	query := `
	SELECT min(CreatedAt) FROM orders
	`
	row, _ := repo.db.Query(query)
	var Date string
	row.Scan(&Date)
	return Date
}

// Returns the number of orders for the specified period, grouped by day within a month.
// month from 1 to 12. Year as year.
func (repo *OrderRepository) OrderedItemsByDay(month, year int) (map[string]interface{}, error) {
	query := `
		SELECT EXTRACT(DAY FROM createdat) AS day, COUNT(*) AS order_count
		FROM orders
		WHERE EXTRACT(MONTH FROM createdat) = $1 `

	if year != -1 {
		query += `AND EXTRACT(YEAR FROM createdat) = $2 `
	}
	query += `GROUP BY day ORDER BY day`

	var rows *sql.Rows
	var err error

	if year != -1 {
		rows, err = repo.db.Query(query, month, year)
	} else {
		rows, err = repo.db.Query(query, month)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch ordered items by day: %w", err)
	}

	defer rows.Close()

	orderedItems := make([]map[string]int, 0)
	for rows.Next() {
		var day int
		var orderCount int
		if err := rows.Scan(&day, &orderCount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		orderedItems = append(orderedItems, map[string]int{fmt.Sprintf("%d", day): orderCount})
	}

	var result map[string]interface{}
	if year != -1 {
		result = map[string]interface{}{
			"period":       "day",
			"month":        month,
			"year":         year,
			"orderedItems": orderedItems,
		}
	} else {
		result = map[string]interface{}{
			"period":       "day",
			"month":        month,
			"orderedItems": orderedItems,
		}
	}

	return result, nil
}

// Returns the number of orders for the specified period, grouped by month within a year
func (repo *OrderRepository) OrderedItemsByMonth(year int) (map[string]interface{}, error) {
	query := `
		SELECT 
			TO_CHAR(o.createdat, 'Month') AS month,
			COUNT(o.ID) AS total_orders
		FROM 
			orders o
		WHERE 
			EXTRACT(YEAR FROM o.createdat) = $1
			AND o.status = 'closed'
		GROUP BY 
			TO_CHAR(o.createdat, 'Month'), EXTRACT(MONTH FROM o.createdat)
		ORDER BY 
			EXTRACT(MONTH FROM o.createdat);
	`
	rows, err := repo.db.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderedItems := make(map[string]interface{})

	// Iterate over the rows and populate the map
	for rows.Next() {
		var month string
		var orderCount int

		if err := rows.Scan(&month, &orderCount); err != nil {
			return nil, err
		}

		month = strings.TrimSpace(month)

		orderedItems[month] = orderCount
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"period":       "month",
		"year":         year,
		"orderedItems": orderedItems,
	}

	return result, nil
}
