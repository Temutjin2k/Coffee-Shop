package service

import (
	"errors"
	"fmt"
	"log"
	"sort"
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
func (s *OrderService) AddOrder(order models.Order) error {
	if order.Items == nil || strings.TrimSpace(order.CustomerName) == "" {
		return errors.New("something wrong with your requested order")
	}
	for _, order := range order.Items {
		if order.Quantity < 1 {
			return errors.New("something wrong with your requested order")
		}
	}

	return s.orderRepo.Add(order)
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
	if updatedOrder.Items == nil || strings.TrimSpace(updatedOrder.CustomerName) == "" {
		return errors.New("something wrong with your updated order")
	}
	for _, order := range updatedOrder.Items {
		if order.Quantity < 1 {
			return errors.New("something wrong with your updated order")
		}
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
	log.Println(start)
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid time format of endDate")
	}
	log.Println(end)

	return s.orderRepo.GetNumberOfItems(start, end)
}

func (s *OrderService) SearchService(minPrice, maxPrice int, args []string, querySrting string) error {
	return nil
}

func (s *OrderService) GetOrdersByPeriod(period, month, year string) (map[string]interface{}, error) {
	return nil, nil
}
