package dal

import (
	"database/sql"
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

func (repo *InventoryRepository) Exists(ID string) bool {
	queryIfExists := `
	select IngredientID from inventory where IngredientID = $1
	`
	rows, err := repo.db.Query(queryIfExists, ID)
	if err != nil {
		return false
	}
	return rows.Next()
}

func (repo *InventoryRepository) SubtractIngredients(ingredients map[string]float64) error {
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

func (repo *InventoryRepository) UpdateItemRepo(id string, newItem models.InventoryItem) error {
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

func (repo *InventoryRepository) DeleteItemRepo(id string) error {
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
