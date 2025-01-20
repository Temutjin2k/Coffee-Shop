package handler

import (
	"encoding/json"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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

func (h *AggregationHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	querySrting := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")
	if querySrting == "" {
		h.logger.Error("Search query string is required", "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Search query string is required", http.StatusBadRequest)
		return
	}
	var args []string
	if filter != "" {
		args = strings.Split(filter, ",")
	}
	for _, v := range args {
		if v != "orders" && v != "menu" && v != "inventory" && v != "all" {
			h.logger.Error("Incorrect search arguments", "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, "Incorrect search arguments", http.StatusBadRequest)
			return
		}
	}
	var MinPrice int
	if minPrice != "" {
		MinPriceTemp, err := strconv.Atoi(minPrice)
		if err != nil {
			h.logger.Error("Min Price should be number", "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, "Min Price should be number", http.StatusBadRequest)
			return
		}
		MinPrice = MinPriceTemp
	} else {
		MinPrice = 0
	}

	var MaxPrice int
	if minPrice != "" {
		MaxPriceTemp, err := strconv.Atoi(minPrice)
		if err != nil {
			h.logger.Error("Max Price should be number", "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, "Max Price should be number", http.StatusBadRequest)
			return
		}
		MinPrice = MaxPriceTemp
	} else {
		MaxPrice = 999999
	}

	MaxPrice, err := strconv.Atoi(maxPrice)
	if err != nil {
		h.logger.Error("Max Price should be number", "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Max Price should be number", http.StatusBadRequest)
		return
	}
	err = h.orderService.SearchService(MinPrice, MaxPrice, args, querySrting)
	if err != nil {
		h.logger.Error(err.Error(), "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Can not searched", http.StatusInternalServerError)
		return
	}
}

/*
	2
	GET /reports/search: Search through orders, menu items, and customers with partial matching and ranking.
	##### Parameters:

		q (required): Search query string
		filter (optional): What to search through, can be multiple values comma-separated:
			orders (search in customer names and order details)
			menu (search in item names and descriptions)
			all (default, search everywhere)
		minPrice (optional): Minimum order/item price to include
		maxPrice (optional): Maximum order/item price to include
*/

func (h *AggregationHandler) Search(w http.ResponseWriter, r *http.Request) {
}

/*
	3
	GET /reports/orderedItemsByPeriod?period={day|month}&month={month}: Returns the number of orders for the specified period, grouped by day within a month or by month within a year. The period parameter can take the value day or month. The month parameter is optional and used only when period=day.

##### Parameters:

	period (required):
	    day: Groups data by day within the specified month.
	    month: Groups data by month within the specified year.
	month (optional): Specifies the month (e.g., october). Used only if period=day.
	year (optional): Specifies the year. Used only if period=month.
*/
func (h *AggregationHandler) OrderByPeriod(w http.ResponseWriter, r *http.Request) {
}
