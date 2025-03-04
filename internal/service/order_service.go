package service

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type OrderService struct {
	orderRepo     dal.OrderRepository
	menuRepo      dal.MenuRepository
	inventoryRepo dal.InventoryRepository
}

func NewOrderService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository, inventoryRepo dal.InventoryRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
	}
}

// AddOrder adds a new order to the repository
func (s *OrderService) AddOrder(order models.Order) (models.BatchOrderInfo, error) {
	if err := validateOrder(order); err != nil {
		return models.BatchOrderInfo{
			OrderID:      order.ID,
			CustomerName: order.CustomerName,
			Status:       models.StatusOrderRejected,
			Reason:       err.Error(),
			Total:        0,
		}, err
	}

	return s.orderRepo.Add(order)
}

//	{
//	    "processed_orders": [
//	        {
//	            "order_id": 123,
//	            "customer_name": "Alice",
//	            "status": "accepted",
//	            "total": 15.50
//	        },
//	        {
//	            "order_id": 124,
//	            "customer_name": "Bob",
//	            "status": "rejected",
//	            "reason": "insufficient_inventory"
//	        }
//	    ],
//	    "summary": {
//	        "total_orders": 2,
//	        "accepted": 1,
//	        "rejected": 1,
//	        "total_revenue": 15.50,
//	        "inventory_updates": [
//	            {
//	                "ingredient_id": 1,
//	                "name": "Coffee Beans",
//	                "quantity_used": 100,
//	                "remaining": 2400
//	            }
//	        ]
//	    }
//	}

func (s *OrderService) BulkOrders(orders []models.Order) (models.BatchOrdersResponce, error) {
	proccesedOrdersInfo := []models.BatchOrderInfo{}
	summary := models.BatchOrderSummary{
		TotalOrders: len(orders),
	}
	// inventoryUpdates := []models.BatchOrderInventoryUpdate{}
	for _, order := range orders {
		orderInfo, err := s.AddOrder(order)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		if orderInfo.Status == "accepted" {
			summary.Accepted++
		} else {
			summary.Rejected++
		}

		summary.TotalRevenue += orderInfo.Total

		proccesedOrdersInfo = append(proccesedOrdersInfo, orderInfo)
	}

	result := models.BatchOrdersResponce{
		Processed_orders: proccesedOrdersInfo,
		Summary:          summary,
	}
	return result, nil
}

// GetAllOrders retrieves all orders from the repository
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *OrderService) GetOrder(OrderID int) (models.Order, error) {
	flag := false
	AllOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return models.Order{}, err
	}
	var NeededOrder models.Order
	for i, Order := range AllOrders {
		if Order.ID == OrderID {
			flag = true
			NeededOrder = AllOrders[i]
		}
	}
	if flag {
		return NeededOrder, nil
	}
	return models.Order{}, errors.New("the order with given ID soes not exist")
}

// UpdateOrder updates an existing order
func (s *OrderService) UpdateOrder(updatedOrder models.Order, OrderID string) error {
	if err := validateOrder(updatedOrder); err != nil {
		return err
	}
	return s.orderRepo.SaveUpdatedOrder(updatedOrder, OrderID)
}

func (s *OrderService) GetTotalSales() (models.TotalSales, error) {
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return models.TotalSales{}, err
	}

	// Counting sales amount
	totalSales := models.TotalSales{}

	for _, order := range existingOrders {
		for _, item := range order.Items {
			totalSales.TotalSales += item.Quantity
		}
	}
	return totalSales, nil
}

// Returns Popular Items sorted in decreasing order. Number of returned items depends on passing value(popularItemsNum)
func (s *OrderService) GetPopularItems(popularItemsNum int) (models.PopularItems, error) {
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return models.PopularItems{}, err
	}

	// Should return sorted decreasing array
	itemMap := make(map[string]int)
	for _, order := range existingOrders {
		for _, item := range order.Items {
			itemMap[item.ProductID] += item.Quantity
		}
	}

	sortedItems := make([]models.OrderItem, 0, len(itemMap))
	for productID, quantity := range itemMap {
		sortedItems = append(sortedItems, models.OrderItem{ProductID: productID, Quantity: quantity})
	}

	// Sorting in decresing order
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].Quantity > sortedItems[j].Quantity
	})

	// To prevent from out of range
	if popularItemsNum > len(sortedItems) {
		popularItemsNum = len(sortedItems)
	}

	popularItems := models.PopularItems{Items: sortedItems[:popularItemsNum]} // potential out of range
	return popularItems, nil
}

func (s *OrderService) DeleteOrderByID(OrderID int) error {
	return s.orderRepo.DeleteOrder(OrderID)
}

func (s *OrderService) CloseOrder(OrderID string) error {
	return s.orderRepo.CloseOrderRepo(OrderID)
}

func (s *OrderService) GetNumberOfItems(startDate, endDate string) (map[string]int, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid time format of startDate")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid time format of endDate")
	}

	return s.orderRepo.GetNumberOfItems(start, end)
}

func (s *OrderService) SearchService(minPrice, maxPrice int, args []string, querySrting string) error {
	return nil
}

func (s *OrderService) GetOrderedItemsByPeriod(period, month, year string) (map[string]interface{}, error) {
	if period == "" {
		return nil, fmt.Errorf("period is required")
	}

	if period == "day" {
		if month == "" {
			return nil, fmt.Errorf("period equal to 'day', but month not provided")
		}
	} else if period == "month" {
		if year == "" {
			return nil, fmt.Errorf("period equal to 'month', but year not provided")
		}
	} else {
		return nil, fmt.Errorf("invalid period value, must be 'day' or 'month'")
	}

	if period == "day" {
		monthInt := getMonthNumber(strings.ToLower(month))
		if monthInt == -1 {
			return nil, fmt.Errorf("%s, month does not exist", month)
		}

		var yearInt int
		if year == "" {
			yearInt = -1
		} else {
			n, err := strconv.Atoi(year)
			if err != nil {
				return nil, fmt.Errorf("year must be a number")
			}
			yearInt = n
		}

		return s.orderRepo.OrderedItemsByDay(monthInt, yearInt)
	} else if period == "month" {

		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return nil, fmt.Errorf("year should be number")
		}
		return s.orderRepo.OrderedItemsByMonth(yearInt)
	}

	return nil, fmt.Errorf("invalid inputs. Period: %v, Month: %s, Year: %s", period, month, year)
}

func getMonthNumber(month string) int {
	months := map[string]int{
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"may":       5,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}

	v, ok := months[month]
	if !ok {
		return -1
	}
	return v
}

func validateOrder(order models.Order) error {
	if order.Items == nil {
		return errors.New("no items provided. Array of items it required")
	}

	if strings.TrimSpace(order.CustomerName) == "" {
		return errors.New("customer name is required")
	}
	for _, order := range order.Items {
		if order.Quantity < 1 {
			return errors.New("quantity a product must be greater than zero")
		}
	}

	return nil
}
