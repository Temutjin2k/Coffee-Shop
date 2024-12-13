package dal

import (
	"database/sql"
	"encoding/json"
	"hot-coffee/models"
	"os"
)

// InventoryRepository implements InventoryRepository using JSON files
type InventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository creates a new FileInventoryRepository
func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (repo *InventoryRepository) GetAll() ([]models.InventoryItem, error) {
	content, err := os.ReadFile("qwe")
	if err != nil {
		return nil, err
	}

	var inventoryItems []models.InventoryItem
	err = json.Unmarshal(content, &inventoryItems)
	return inventoryItems, err
}

// Check if ingridient by given ID exists
func (repo *InventoryRepository) Exists(ID string) bool {
	inventoryItems, err := repo.GetAll()
	if err != nil {
		return false
	}

	for _, item := range inventoryItems {
		if item.IngredientID == ID {
			return true
		}
	}
	return false
}

func (repo *InventoryRepository) SubtractIngredients(ingredients map[string]float64) error {
	inventoryItems, err := repo.GetAll()
	if err != nil {
		return err
	}

	for i, inventoryItem := range inventoryItems {
		if value, exists := ingredients[inventoryItem.IngredientID]; exists {
			inventoryItems[i].Quantity -= value
		}
	}

	jsonData, err := json.MarshalIndent(inventoryItems, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile("qwe", jsonData, 0o644)
}

func (repo *InventoryRepository) SaveAll(items []models.InventoryItem) error {
	jsonData, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("qwe", jsonData, 0o644)
}
