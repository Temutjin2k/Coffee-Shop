package server

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

func ServerLaunch(db *sql.DB, logger *slog.Logger) {
	orderRepo := dal.NewOrderRepository(db)
	menuRepo := dal.NewMenuRepository("qwe")
	inventoryRepo := dal.NewInventoryRepository("qwe")

	// Initialize services (Business Logic Layer)
	orderService := service.NewOrderService(*orderRepo, *menuRepo)
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

	mux.HandleFunc("GET /reports/total-sales", reportHandler.TotalSalesHandler)
	mux.HandleFunc("GET /reports/popular-items", reportHandler.PopularItemsHandler)

	address := "http://localhost:" + "8080" + "/"
	fmt.Println("Server launched on address:", address)

	logger.Info("Application started", "Address", address)
	log.Fatal(http.ListenAndServe(":"+"8080", mux))
}
