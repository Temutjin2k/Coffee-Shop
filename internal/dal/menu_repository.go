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

func (repo *MenuRepository) DeleteMenuItemRepo(MenuItemID string) error {
	queryDeleteMenuItem := `
	delete from menu_items
	where ID = $1
	`
	_, err := repo.db.Exec(queryDeleteMenuItem, MenuItemID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MenuRepository) UpdateMenuItemRepo(menuItem models.MenuItem) error {
	queryUpdateMenu := `
	update menu_items
	set Name = $1, Description = $2, Price = $3
	where ID = $4
	`
	_, err := repo.db.Exec(queryUpdateMenu, menuItem.Name, menuItem.Description, menuItem.Price, menuItem.ID)
	if err != nil {
		return err
	}
	for _, v := range menuItem.Ingredients {
		queryUpdateMenuIngredients := `
		update menu_item_ingredients 
		set IngredientID = $1, Quantity = $2
		where MenuID = $3
		`
		_, err = repo.db.Exec(queryUpdateMenuIngredients, v.IngredientID, v.Quantity, menuItem.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *MenuRepository) SaveAll(menuItems []models.MenuItem) error {
	jsonData, err := json.MarshalIndent(menuItems, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("qwe", jsonData, 0o644)
}
