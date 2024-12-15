package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type InventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) *InventoryService {
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

func (s *InventoryService) GetItem(id string) (models.InventoryItem, error) {
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

func (s *InventoryService) UpdateItem(id string, newItem models.InventoryItem) error {
	if !s.inventoryRepo.Exists(id) {
		return errors.New("inventory item does not exists")
	}
	return s.inventoryRepo.UpdateItemRepo(id, newItem)
}

func (s *InventoryService) DeleteItem(id string) error {
	if !s.inventoryRepo.Exists(id) {
		return errors.New("inventory item does not exists")
	}
	return s.inventoryRepo.DeleteItemRepo(id)
}

func (s *InventoryService) Exists(id string) bool {
	return s.inventoryRepo.Exists(id)
}
