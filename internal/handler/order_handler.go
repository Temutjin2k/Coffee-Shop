package handler

import (
	"encoding/json"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type OrderHandler struct {
	orderService *service.OrderService
	logger       *slog.Logger
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *service.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, logger: logger}
}

// PostOrder creates new Order
func (h *OrderHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	err = h.orderService.AddOrder(newOrder)
	if err != nil {
		ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	Orders, err := h.orderService.GetAllOrders()
	if err != nil {
		ErrorHandler.Error(w, "Can not read order data from server", http.StatusInternalServerError)
	}
	var RequestedOrder models.Order
	for i, Order := range Orders {
		if Order.ID == r.PathValue("id") {
			RequestedOrder = Orders[i]
		}
	}
	jsonData, err := json.MarshalIndent(RequestedOrder, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Can not convert order data to json", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	Orders, err := h.orderService.GetAllOrders()
	if err != nil {
		ErrorHandler.Error(w, "Can not read order data from server", http.StatusInternalServerError)
	}
	jsonData, err := json.MarshalIndent(Orders, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Can not convert order data to json", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (h *OrderHandler) PutOrder(w http.ResponseWriter, r *http.Request) {
	h.menuService.MenuCheckByID()
	Requestcontent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	var RequestrOrder models.Order
	err = json.Unmarshal(Requestcontent, &RequestrOrder)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	h.orderService.UpdateOrder(RequestrOrder)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
}
