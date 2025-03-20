# Coffee-Shop

## Context

Have you ever wondered how your favorite coffee shop manages to handle a flurry of orders during the morning rush, keep track of inventory so they never run out of your preferred blend, or remember that you like your coffee with an extra shot of espresso?

Behind the scenes, coffee shops rely on sophisticated management systems that coordinate orders, inventory, menu items, and customer preferences in real-time. These systems ensure that baristas can focus on crafting the perfect cup while the technology handles the complexities of order processing, stock management, and data recording.

The `Coffee-shop` (coffee shop management system) project is a simplified version of these real-world applications, designed to give you hands-on experience with the core principles behind such operational software. Imagine an application that allows staff to:

- **Manage Orders:** Create, modify, close, and delete customer orders efficiently.
- **Oversee Inventory:** Track ingredient stock levels to prevent shortages and ensure freshness.
- **Update the Menu:** Add new drinks or pastries, adjust prices as needed, and keep the offerings up to date.


### API Endpoints

#### Orders

| Method | Endpoint          | Description                       
|--------|-------------------|-----------------------------------|
| POST   | `/orders`         | Creates a new order.            
| GET    | `/orders`         | Retrieves all orders.             
| GET    | `/orders/{id}`    | Retrieves a specific order by ID. 
| PUT    | `/orders/{id}`    | Updates an existing order.        | 
| DELETE | `/orders/{id}`    | Deletes an order.                 | 
| POST   | `/orders/{id}/close` | Closes an open order.         
| GET    | `/orders/numberOfOrderedItems` | Returns a list of ordered items and their quantities for a specified time period|  
| POST    | `/orders/batch-process` | Process multiple orders simultaneously while ensuring inventory consistency.|  

#### Menu Items

| Method | Endpoint          | Description                        
|--------|-------------------|------------------------------------|
| POST   | `/menu`           | Adds a new menu item.             |               
| GET    | `/menu`           | Retrieves all menu items.          |                  
| GET    | `/menu/{id}`      | Retrieves a specific menu item.    | 
| PUT    | `/menu/{id}`      | Updates an existing menu item.     | 
| DELETE | `/menu/{id}`      | Deletes a menu item.               | 


#### Inventory

| Method | Endpoint          | Description                        | 
|--------|-------------------|------------------------------------|
| POST   | `/inventory`      | Adds a new inventory item.        | 
| GET    | `/inventory`      | Retrieves all inventory items.     | 
| GET    | `/inventory/{id}`  | Retrieves a specific inventory item. |
| PUT    | `/inventory/{id}`  | Updates an inventory item.         | 
| DELETE | `/inventory/{id}`  | Deletes an inventory item.         | 
| GET | `/inventory/getLeftOvers`      | Returns the inventory leftovers in the coffee shop, including sorting and pagination options              | 

#### Aggregations

| Method | Endpoint                  | Description                        |
|--------|---------------------------|------------------------------------|
| GET    | `/reports/total-sales`    | Retrieves the total sales amount.  |
| GET    | `/reports/popular-items`   | Retrieves a list of popular menu items. | 
| GET    | `/reports/orderedItemsByPeriod`   | Returns the number of orders for the specified period, grouped by day within a month or by month within a year.| 
| GET    | `/reports/search`   | Search through orders, menu items, and customers with partial matching and ranking. | 






## How to Run the Project

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Temutjin2k/Coffee-Shop.git
   cd Coffee-Shop
   ```

2. **Start the server in docker**:
   ```bash
   docker compose up
   ```
