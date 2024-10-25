package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
	logger           *slog.Logger
}

func NewInventoryHandler(inventoryService *service.InventoryService, logger *slog.Logger) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService, logger: logger}
}

func (h *InventoryHandler) PostInventory(w http.ResponseWriter, r *http.Request) {
	var newItem models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		h.logger.Error("Could not decode json data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusInternalServerError)
		return
	}

	// Checking for empty fieldss
	if newItem.Name == "" || newItem.IngredientID == "" || newItem.Unit == "" || newItem.Quantity <= 0 {
		h.logger.Error("Some fields are empty, equal or less than zero", "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Some fields are empty, equal or less than zero", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.AddInventoryItem(newItem)
	if err != nil {
		h.logger.Error("Could not add new inventory item", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not add new inventory item Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
	w.WriteHeader(http.StatusCreated)
}

func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	inventoryItems, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		h.logger.Error("Could not get inventory items", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not get inventory items", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(inventoryItems)
	if err != nil {
		h.logger.Error("Could not encode json data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not encode request json data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	inventoryItem, err := h.inventoryService.GetItem(id)
	if err != nil {
		h.logger.Error("Could not get inventory item", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not get inventory item", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(inventoryItem)
	if err != nil {
		h.logger.Error("Could not encode json data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not encode json data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *InventoryHandler) PutInventoryItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		h.logger.Error("Could not decode json data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	if newItem.Name == "" || newItem.IngredientID == "" || newItem.Unit == "" || newItem.Quantity <= 0 {
		h.logger.Error("Some fields are empty", "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Some fields are empty, equal or less than zero", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.UpdateItem(r.PathValue("id"), newItem)
	if err != nil {
		h.logger.Error("Error updating inventory item", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Error updating inventory item Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.inventoryService.DeleteItem(id)
	if err != nil {
		h.logger.Error("Could not delete inventory item", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not delete inventory item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}
