package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"hot-coffee/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (repo *InventoryRepository) GetAll() ([]models.InventoryItem, error) {
	queryGetIngridients := `
	select IngredientID, Name, Quantity, Unit from inventory
	`
	rows, err := repo.db.Query(queryGetIngridients)
	if err != nil {
		return []models.InventoryItem{}, err
	}
	var InventoryItems []models.InventoryItem

	for rows.Next() {
		var InventoryItem models.InventoryItem
		err = rows.Scan(&InventoryItem.IngredientID, &InventoryItem.Name, &InventoryItem.Quantity, &InventoryItem.Unit)
		if err != nil {
			return []models.InventoryItem{}, nil
		}
		InventoryItems = append(InventoryItems, InventoryItem)
	}
	return InventoryItems, nil
}

func (repo *InventoryRepository) Exists(ID int) bool {
	queryIfExists := `
	select IngredientID from inventory where IngredientID = $1
	`
	rows, err := repo.db.Query(queryIfExists, ID)
	if err != nil {
		return false
	}
	return rows.Next()
}

func (repo *InventoryRepository) SubtractIngredients(ingredients map[int]float64) error {
	for key, value := range ingredients {
		queryToSubtract := `
	        update inventory
	        set Quantity  = Quantity - $1
	        where IngredientID = $2
	    `
		_, err := repo.db.Exec(queryToSubtract, value, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *InventoryRepository) AddInventoryItemRepo(item models.InventoryItem) error {
	queryToAddInventory := `
	insert into inventory (Name, Quantity, Unit) values
	($1, $2, $3)
	`
	_, err := repo.db.Exec(queryToAddInventory, item.Name, item.Quantity, item.Unit)
	if err != nil {
		return err
	}
	return nil
}

func (repo *InventoryRepository) UpdateItemRepo(id int, newItem models.InventoryItem) error {
	queryToUpdate := `
	update inventory
	set Quantity = $1, set Name = $2, set Unit = $3
	where IngredientID = $4
	`
	_, err := repo.db.Exec(queryToUpdate, newItem.Quantity, newItem.Name, newItem.Unit, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *InventoryRepository) DeleteItemRepo(id int) error {
	queryToDelete := `
	delete from inventory
	where ID = $1
	`
	_, err := repo.db.Exec(queryToDelete, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *InventoryRepository) GetLeftOvers(sortBy, page, pageSize string) (map[string]any, error) {
	// Set defaults if parameters are missing
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	// Calculate offset for pagination
	offset := (pageNum - 1) * pageSizeNum

	// Prepare base query for fetching inventory data
	query := `
        SELECT i.IngredientID, i.Name, i.Quantity, i.Unit
        FROM inventory i
    `

	// Add sorting condition based on sortBy parameter
	switch sortBy {
	case "price":
		query += " ORDER BY i.Price"
	case "quantity":
		query += " ORDER BY i.Quantity"
	default:
		return nil, errors.New("invalid sortBy value, must be 'price' or 'quantity'")
	}

	// Add pagination to query
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSizeNum, offset)

	// Prepare response structure
	var leftovers []map[string]any
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	// Fetch inventory data
	for rows.Next() {
		var ingredientID int
		var name string
		var quantity int
		var unit string
		if err := rows.Scan(&ingredientID, &name, &quantity, &unit); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		leftovers = append(leftovers, map[string]any{
			"ingredientID": ingredientID,
			"name":         name,
			"quantity":     quantity,
			"unit":         unit,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	// Get total number of pages by querying the total number of records in the inventory table
	var totalItems int
	err = repo.db.QueryRow("SELECT COUNT(*) FROM inventory").Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to count total items: %v", err)
	}

	totalPages := (totalItems + pageSizeNum - 1) / pageSizeNum
	hasNextPage := pageNum < totalPages

	// Return the response with pagination info and data
	response := map[string]any{
		"currentPage": pageNum,
		"hasNextPage": hasNextPage,
		"pageSize":    pageSizeNum,
		"totalPages":  totalPages,
		"data":        leftovers,
	}

	return response, nil
}
