package orderHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"hot-coffee/config"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/models"
)

func Deleteorder(w http.ResponseWriter, ObjectID string) {
	OrderContents, err := ioutil.ReadFile(config.BaseDir + "/orders.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read orders from server", http.StatusInternalServerError)
		return
	}
	var Orders []models.Order
	json.Unmarshal(OrderContents, &Orders)

	NewOrders := make([]models.Order, 0)
	for _, order := range Orders {
		if order.ID != ObjectID {
			var NewOrder models.Order
			NewOrder.CreatedAt = order.CreatedAt
			NewOrder.CustomerName = order.CustomerName
			NewOrder.ID = order.ID
			NewOrder.Items = order.Items
			NewOrder.Status = order.Status
			NewOrders = append(NewOrders, NewOrder)
		}
	}
	jsonData, err := json.MarshalIndent(NewOrders, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Could transfer orders to json file", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile(config.BaseDir+"/orders.json", jsonData, os.ModePerm)
	if err != nil {
		ErrorHandler.Error(w, "Could not rewwrite orders in server", http.StatusInternalServerError)
		return
	}
}