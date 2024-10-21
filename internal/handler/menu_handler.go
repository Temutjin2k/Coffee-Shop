package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type MenuHandler struct {
	menuService *service.MenuService
}

func (h *MenuHandler) NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

func (h *MenuHandler) MenuHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path[1:], "/")
	switch len(parts) {
	case 1:
		switch r.Method {
		case http.MethodPost:
			h.PostMenu(w, r)
		case http.MethodGet:
			h.GetMenuItems(w, r)
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 2:
		switch r.Method {
		case http.MethodPut:
			h.PutMenu(w, r)
		case http.MethodGet:
			h.GetMenuItem(w, r)
		case http.MethodDelete:
			h.DeleteMenu(w, r)
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	default:
		ErrorHandler.Error(w, "Something wrong with your request", http.StatusBadRequest)
	}
}

func (h *MenuHandler) PostMenu(w http.ResponseWriter, r *http.Request) {
	var newItem models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	// Use the service to check if the item already exists
	if h.menuService.MenuCheck(newItem) {
		ErrorHandler.Error(w, "The requested menu item already exists in current menu", http.StatusBadRequest)
		return
	}

	// Add the new menu item using the service
	if err := h.menuService.AddMenuItem(newItem); err != nil {
		ErrorHandler.Error(w, "Could not add menu item", http.StatusInternalServerError)
		return
	}
}

func (h *MenuHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
}

func (h *MenuHandler) PutMenu(w http.ResponseWriter, r *http.Request) {
}

func (h *MenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
}