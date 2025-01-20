package server

import (
	"database/sql"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"log"
	"log/slog"
	"net/http"
)

func ServerLaunch(db *sql.DB, logger *slog.Logger) {
	// config.InitDB(db)
	orderRepo := dal.NewOrderRepository(db)
	menuRepo := dal.NewMenuRepository(db)
	inventoryRepo := dal.NewInventoryRepository(db)

	// Initialize services (Business Logic Layer)
	orderService := service.NewOrderService(*orderRepo, *menuRepo, *inventoryRepo)
	menuService := service.NewMenuService(*menuRepo, *inventoryRepo)
	inventoryService := service.NewInventoryService(*inventoryRepo)

	// Initialize handlers (Presentation Layer)
	orderHandler := handler.NewOrderHandler(orderService, menuService, logger)
	menuHandler := handler.NewMenuHandler(menuService, logger)
	inventoryHandler := handler.NewInventoryHandler(inventoryService, logger)
	reportHandler := handler.NewAggregationHandler(orderService, logger)

	// Setup HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", orderHandler.PostOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrder)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.PutOrder)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrder)
	mux.HandleFunc("POST /orders/{id}/close", orderHandler.CloseOrder)

	mux.HandleFunc("GET /orders/numberOfOrderedItems", orderHandler.GetNumberOfOrdered)

	// TODO
	mux.HandleFunc("POST /orders/batch-process", orderHandler.BatchOrders)

	mux.HandleFunc("POST /menu", menuHandler.PostMenu)
	mux.HandleFunc("GET /menu", menuHandler.GetMenu)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetMenuItem)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.PutMenuItem)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteMenuItem)

	mux.HandleFunc("POST /inventory", inventoryHandler.PostInventory)
	mux.HandleFunc("GET /inventory", inventoryHandler.GetInventory)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.GetInventoryItem)
	mux.HandleFunc("PUT /inventory/{id}", inventoryHandler.PutInventoryItem)
	mux.HandleFunc("DELETE /inventory/{id}", inventoryHandler.DeleteInventoryItem)

	mux.HandleFunc("GET /inventory/getLeftOvers", inventoryHandler.GetLeftOvers)

	mux.HandleFunc("GET /reports/total-sales", reportHandler.TotalSalesHandler)
	mux.HandleFunc("GET /reports/popular-items", reportHandler.PopularItemsHandler)
	mux.HandleFunc("GET /reports/search", reportHandler.SearchHandler)

	// TOOD
	mux.HandleFunc("GET /reports/orderedItemsByPeriod", reportHandler.OrderByPeriod)

	logger.Info("Application started", "Address", "http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
