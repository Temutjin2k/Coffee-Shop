package service

import (
	"errors"
	"strconv"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type IInventoryService interface {
	AddInventoryItem(item models.InventoryItem) error
	GetAllInventoryItems() ([]models.InventoryItem, error)
	GetItem(id int) (models.InventoryItem, error)
	UpdateItem(id int, newItem models.InventoryItem) error
	DeleteItem(id int) error
	Exists(id int) bool
	GetLeftOvers(sortBy, page, pageSize string) (map[string]any, error)
}

type InventoryService struct {
	inventoryRepo dal.IInventoryRepository
}

func NewInventoryService(inventoryRepo dal.IInventoryRepository) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

func (s *InventoryService) AddInventoryItem(item models.InventoryItem) error {
	return s.inventoryRepo.AddInventoryItemRepo(item)
}

func (s *InventoryService) GetAllInventoryItems() ([]models.InventoryItem, error) {
	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		return []models.InventoryItem{}, nil
	}
	return items, nil
}

func (s *InventoryService) GetItem(id int) (models.InventoryItem, error) {
	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}

	for _, item := range items {
		if item.IngredientID == id {
			return item, nil
		}
	}
	return models.InventoryItem{}, errors.New("inventory item does not exists")
}

func (s *InventoryService) UpdateItem(id int, newItem models.InventoryItem) error {
	if !s.inventoryRepo.Exists(id) {
		return errors.New("inventory item does not exists")
	}
	return s.inventoryRepo.UpdateItemRepo(id, newItem)
}

func (s *InventoryService) DeleteItem(id int) error {
	if !s.inventoryRepo.Exists(id) {
		return errors.New("inventory item does not exists")
	}
	return s.inventoryRepo.DeleteItemRepo(id)
}

func (s *InventoryService) Exists(id int) bool {
	return s.inventoryRepo.Exists(id)
}

func (s *InventoryService) GetLeftOvers(sortBy, page, pageSize string) (map[string]any, error) {
	if sortBy == "" {
		sortBy = "price"
	}

	if page == "" {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 0 {
		return nil, errors.New("invalid page parametr. page must be positive integer")
	}

	if pageSize == "" {
		pageSize = "10"
	}
	pageSizeNum, err := strconv.Atoi(pageSize)

	if err != nil || pageSizeNum < 0 {
		return nil, errors.New("invalid pageSize parametr. pageSize must be positive integer")
	}

	if sortBy != "price" && sortBy != "quantity" {
		return nil, errors.New("invalid sortBy value, must be 'price' or 'quantity")
	}

	return s.inventoryRepo.GetLeftOvers(sortBy, page, pageSize)
}
