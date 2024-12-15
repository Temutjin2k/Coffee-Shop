package handler

import (
	"encoding/json"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"log/slog"
	"net/http"
)

type AggregationHandler struct {
	orderService *service.OrderService
	logger       *slog.Logger
}

func NewAggregationHandler(orderService *service.OrderService, logger *slog.Logger) *AggregationHandler {
	return &AggregationHandler{orderService: orderService, logger: logger}
}

// Return all saled items as key and quantity as value in JSON
func (h *AggregationHandler) TotalSalesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	totalSales, err := h.orderService.GetTotalSales()
	if err != nil {
		h.logger.Error("Error getting data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Error getting data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalSales)

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

// Returns Each item as key and quatity as value
func (h *AggregationHandler) PopularItemsHandler(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.orderService.GetPopularItems(3)
	if err != nil {
		h.logger.Error("Error in orderService GetPopularItems", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Error getting data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(popularItems)

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}
