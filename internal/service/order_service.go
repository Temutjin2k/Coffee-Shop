package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OrderService struct {
	orderRepo dal.OrderRepository
	menuRepo  dal.MenuRepository
}

func NewOrderService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

// AddOrder adds a new order to the repository
func (s *OrderService) AddOrder(order models.Order) error {
	if order.Items == nil || strings.TrimSpace(order.CustomerName) == "" {
		return errors.New("Womething wrong with your requested order")
	}
	for _, order := range order.Items {
		if order.Quantity < 1 {
			return errors.New("Womething wrong with your requested order")
		}
	}

	OrderID, err := s.orderRepo.GetID()
	if err != nil {
		return err
	}
	order.ID = strconv.Itoa(OrderID)
	Location, err := time.LoadLocation("Asia/Aqtau")
	if err != nil {
		return err
	}
	timenow := time.Now().In(Location).Format(time.RFC3339)
	order.CreatedAt = timenow
	order.Status = "open"

	return s.orderRepo.Add(order)
}

// GetAllOrders retrieves all orders from the repository
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *OrderService) GetOrder(OrderID string) (models.Order, error) {
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
	return models.Order{}, errors.New("The order does not exist")
}

// DeleteOrder removes an order by ID
func (s *OrderService) DeleteOrder(orderID string) error {
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for _, order := range existingOrders {
		if order.ID == orderID {
			return s.orderRepo.Delete(orderID)
		}
	}

	return errors.New("order not found")
}

// UpdateOrder updates an existing order
func (s *OrderService) UpdateOrder(updatedOrder models.Order, OrderID string) error {
	if updatedOrder.Items == nil || strings.TrimSpace(updatedOrder.CustomerName) == "" {
		return errors.New("Womething wrong with your updated order")
	}
	for _, order := range updatedOrder.Items {
		if order.Quantity < 1 {
			return errors.New("Womething wrong with your updated order")
		}
	}
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for i, order := range existingOrders {
		if order.ID == OrderID {
			existingOrders[i].CustomerName = updatedOrder.CustomerName
			existingOrders[i].ID = OrderID
			existingOrders[i].Items = updatedOrder.Items
			existingOrders[i].Status = "Open"
		}
	}

	s.orderRepo.SaveAll(existingOrders)
	return nil
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

func (s *OrderService) DeleteOrderByID(OrderID string) error {
	// if !s.orderRepo.Exists(id) {
	// 	return errors.New("inventory item does not exists")
	// }

	Orders, err := s.GetAllOrders()
	if err != nil {
		return err
	}
	NewOrders := make([]models.Order, 0)
	for _, order := range Orders {
		if order.ID != OrderID {
			var NewOrder models.Order
			NewOrder.CreatedAt = order.CreatedAt
			NewOrder.CustomerName = order.CustomerName
			NewOrder.ID = order.ID
			NewOrder.Items = order.Items
			NewOrder.Status = order.Status
			NewOrders = append(NewOrders, NewOrder)
		}
	}
	s.orderRepo.SaveAll(NewOrders)
	return nil
}

func (s *OrderService) CloseOrder(OrderID string) error {
	Orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil
	}
	var ClosingOrder models.Order
	for _, order := range Orders {
		if order.ID == OrderID {
			ClosingOrder.CreatedAt = order.CreatedAt
			ClosingOrder.CustomerName = order.CustomerName
			ClosingOrder.ID = OrderID
			ClosingOrder.Items = order.Items
			ClosingOrder.Status = "closed"
		}
	}
	for i, order := range Orders {
		if order.ID == OrderID {
			Orders[i].CreatedAt = ClosingOrder.CreatedAt
			Orders[i].CustomerName = ClosingOrder.CustomerName
			Orders[i].ID = OrderID
			Orders[i].Items = ClosingOrder.Items
			Orders[i].Status = "closed"
		}
	}
	s.orderRepo.SaveAll(Orders)
	return nil
}
