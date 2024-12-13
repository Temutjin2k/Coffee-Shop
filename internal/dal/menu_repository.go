package dal

import (
	"database/sql"
	"encoding/json"
	"hot-coffee/models"
	"os"
)

// MenuRepository implements MenuRepository using JSON files
type MenuRepository struct {
	db *sql.DB
}

// NewMenuRepository creates a new FileMenuRepository
func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (repo *MenuRepository) GetAll() ([]models.MenuItem, error) {
	queryMenuItems := `
	select ID, Name, Description, Price from menu_items
	`
	rows, err := repo.db.Query(queryMenuItems)
	if err != nil {
		return []models.MenuItem{}, err
	}
	var MenuItems []models.MenuItem
	for rows.Next() {
		var MenuItem models.MenuItem
		rows.Scan(&MenuItem.ID, &MenuItem.Name, &MenuItem.Description, &MenuItem.Price)
		var MenuItemIngredients []models.MenuItemIngredient
		queryMenuItemIngredients := `
	        select IngredientID, Quantity from menu_item_ingredients where MenuID = $1
	    `
		rows1, err := repo.db.Query(queryMenuItemIngredients, MenuItem.ID)
		if err != nil {
			return []models.MenuItem{}, err
		}
		for rows1.Next() {
			var MenuItemIngredient models.MenuItemIngredient
			rows1.Scan(&MenuItemIngredient.IngredientID, &MenuItemIngredient.Quantity)
			MenuItemIngredients = append(MenuItemIngredients, MenuItemIngredient)
		}
		MenuItem.Ingredients = MenuItemIngredients
		MenuItems = append(MenuItems, MenuItem)
	}
	return MenuItems, nil
}

func (repo *MenuRepository) Exists(itemID string) bool {
	items, _ := repo.GetAll()
	for _, item := range items {
		if item.ID == itemID {
			return true
		}
	}
	return false
}

func (repo *MenuRepository) SaveAll(menuItems []models.MenuItem) error {
	jsonData, err := json.MarshalIndent(menuItems, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("qwe", jsonData, 0o644)
}
