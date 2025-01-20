-- Create the database if it doesn't exist
DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'frappuccino') THEN
      EXECUTE 'CREATE DATABASE frappuccino';
   END IF;
END
$$;

CREATE TYPE order_status AS ENUM ('open', 'closed');
CREATE TYPE unit_types AS ENUM ('ml', 'shots', 'g');

CREATE TABLE menu_items (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(50),
    Description TEXT,
    Price NUMERIC(10, 2)
);

CREATE TABLE inventory (
    IngredientID SERIAL PRIMARY KEY,
    Name VARCHAR(50),
    Quantity INT,
    Unit unit_types
);

CREATE TABLE orders (
    ID SERIAL PRIMARY KEY,
    CustomerName VARCHAR(50),
    Status order_status DEFAULT 'open',
    Notes JSONB, -- To store special client's wish
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    OrderID INT,
    ProductID INT,
    Quantity INT,
    PRIMARY KEY (OrderID, ProductID),
    FOREIGN KEY (OrderID) REFERENCES orders(ID),
    FOREIGN KEY (ProductID) REFERENCES menu_items(ID)
);

CREATE TABLE price_history (
    Menu_ItemID INT,
    Price NUMERIC(10, 2),
    Date DATE,
    PRIMARY KEY (Menu_ItemID, Date),
    FOREIGN KEY (Menu_ItemID) REFERENCES menu_items(ID)
);

CREATE TABLE menu_item_ingredients (
    MenuID INT,
    IngredientID INT,
    Quantity INT,
    PRIMARY KEY (MenuID, IngredientID),
    FOREIGN KEY (MenuID) REFERENCES menu_items(ID) ON DELETE CASCADE,
    FOREIGN KEY (IngredientID) REFERENCES inventory(IngredientID)
);

CREATE TABLE order_status_history (
    ID SERIAL PRIMARY KEY,
    OrderID INT,
    OpenedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ClosedAt TIMESTAMP,
    FOREIGN KEY (OrderID) REFERENCES orders(ID)
);

CREATE TABLE inventory_transactions (
    transactionId SERIAL PRIMARY KEY,
    IngredientID INT REFERENCES inventory(IngredientID) ON DELETE CASCADE,
    quantity_change FLOAT NOT NULL,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для таблицы menu_items
CREATE INDEX idx_menu_items_name ON menu_items (Name);

-- Индексы для таблицы inventory
CREATE INDEX idx_inventory_name ON inventory (Name);

-- Индексы для таблицы orders
CREATE INDEX idx_orders_customer_name ON orders (CustomerName);
CREATE INDEX idx_orders_status ON orders (Status);
CREATE INDEX idx_orders_created_at ON orders (CreatedAt);

-- Индексы для таблицы order_items
CREATE INDEX idx_order_items_order_id ON order_items (OrderID);
CREATE INDEX idx_order_items_product_id ON order_items (ProductID);

-- Индексы для таблицы price_history
CREATE INDEX idx_price_history_menu_item_id ON price_history (Menu_ItemID);
CREATE INDEX idx_price_history_date ON price_history (Date);

-- Индексы для таблицы menu_item_ingredients
CREATE INDEX idx_menu_item_ingredients_menu_id ON menu_item_ingredients (MenuID);
CREATE INDEX idx_menu_item_ingredients_ingredient_id ON menu_item_ingredients (IngredientID);

-- Индексы для таблицы order_status_history
CREATE INDEX idx_order_status_history_order_id ON order_status_history (OrderID);

-- Индексы для таблицы inventory_transactions
CREATE INDEX idx_inventory_transactions_ingredient_id ON inventory_transactions (IngredientID);
CREATE INDEX idx_inventory_transactions_created_at ON inventory_transactions (created_at);




-- Функция для логирования изменения цены в price_history
CREATE OR REPLACE FUNCTION log_price_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.price <> OLD.price THEN
        INSERT INTO price_history (Menu_ItemID, Price, Date)
        VALUES (OLD.ID, OLD.price, CURRENT_DATE);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER price_change_trigger
AFTER UPDATE ON menu_items
FOR EACH ROW
EXECUTE FUNCTION log_price_change();

-- Функция для логирования изменения статуса заказа
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.Status <> OLD.Status THEN
        INSERT INTO order_status_history (OrderID, OpenedAt)
        VALUES (OLD.ID, CURRENT_TIMESTAMP);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER order_status_change_trigger
AFTER UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION log_order_status_change();

-- Функция для проверки отрицательных значений в inventory
CREATE OR REPLACE FUNCTION check_inventory_quantity()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.quantity < 0 THEN
        RAISE EXCEPTION 'Quantity in inventory cannot be negative';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_negative_inventory
BEFORE INSERT OR UPDATE ON inventory
FOR EACH ROW
EXECUTE FUNCTION check_inventory_quantity();

-- Функция для проверки отрицательных цен в menu_items
CREATE OR REPLACE FUNCTION check_menu_item_price()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.price <= 0 THEN
        RAISE EXCEPTION 'Price of menu item must be greater than zero';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_negative_price
BEFORE INSERT OR UPDATE ON menu_items
FOR EACH ROW
EXECUTE FUNCTION check_menu_item_price();

-- Функция для проверки нулевого количества в order_items
CREATE OR REPLACE FUNCTION check_order_item_quantity()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.quantity <= 0 THEN
        RAISE EXCEPTION 'Order item quantity must be greater than zero';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_zero_quantity
BEFORE INSERT OR UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION check_order_item_quantity();

--Автоматическое логирование в inventory_transactions.
CREATE OR REPLACE FUNCTION log_inventory_transaction()
RETURNS TRIGGER AS $$
BEGIN

    IF TG_OP = 'UPDATE' THEN
        IF NEW.quantity <> OLD.quantity THEN
            INSERT INTO inventory_transactions(ingredientId, quantity_change, reason, created_at)
            VALUES (
                OLD.ingredientId,
                NEW.quantity - OLD.quantity,
                'Inventory adjustment',
                CURRENT_TIMESTAMP
            );
        END IF;

    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO inventory_transactions(ingredientId, quantity_change, reason, created_at)
        VALUES (
            NEW.ingredientId,
            NEW.quantity,
            'Initial stock',
            CURRENT_TIMESTAMP
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER inventory_change_trigger
AFTER INSERT OR UPDATE ON inventory
FOR EACH ROW
EXECUTE FUNCTION log_inventory_transaction();





-- Mock data for menu_items
INSERT INTO menu_items (Name, Description, Price) VALUES
('Caffe Latte', 'Espresso with steamed milk', 3.50),
('Blueberry Muffin', 'Freshly baked muffin with blueberries', 2.00),
('Espresso', 'Strong and bold coffee', 2.50),
('Cappuccino', 'Espresso with steamed milk and foam', 3.00),
('Mocha', 'Espresso with steamed milk and chocolate', 3.75),
('Iced Latte', 'Iced espresso with milk', 3.80),
('Americano', 'Espresso diluted with hot water', 2.80),
('Carrot Cake', 'Delicious spiced cake with cream cheese frosting', 2.50),
('Vanilla Latte', 'Espresso with steamed milk and vanilla syrup', 3.60),
('Chocolate Croissant', 'Flaky croissant with chocolate filling', 2.80);


-- Mock data for inventory
INSERT INTO inventory (Name, Quantity, Unit) VALUES
('Espresso Shot', 500, 'shots'),
('Milk', 5000, 'ml'),
('Flour', 10000, 'g'),
('Blueberries', 2000, 'g'),
('Sugar', 5000, 'g'),
('Butter', 3000, 'g'),
('Chocolate', 1500, 'g'),
('Coffee Beans', 2000, 'g'),
('Cocoa Powder', 1000, 'g'),
('Vanilla Syrup', 800, 'ml');


-- Mock data for menu_item_ingredients
INSERT INTO menu_item_ingredients (MenuID, IngredientID, Quantity) VALUES
(1, 1, 1),  -- Caffe Latte: 1 Espresso Shot
(1, 2, 200),  -- Caffe Latte: 200 ml Milk
(2, 3, 100),  -- Blueberry Muffin: 100 g Flour
(2, 4, 20),  -- Blueberry Muffin: 20 g Butter
(2, 5, 30),  -- Blueberry Muffin: 30 g Sugar
(3, 1, 1),  -- Espresso: 1 Espresso Shot
(4, 1, 1),  -- Cappuccino: 1 Espresso Shot
(4, 2, 200),  -- Cappuccino: 200 ml Milk
(5, 1, 1),  -- Mocha: 1 Espresso Shot
(5, 2, 200),  -- Mocha: 200 ml Milk
(5, 6, 30),  -- Mocha: 30 g Chocolate
(6, 1, 1),  -- Iced Latte: 1 Espresso Shot
(6, 2, 200),  -- Iced Latte: 200 ml Milk
(7, 1, 1),  -- Americano: 1 Espresso Shot
(8, 3, 100),  -- Carrot Cake: 100 g Flour
(8, 4, 20),  -- Carrot Cake: 20 g Butter
(9, 1, 1),  -- Vanilla Latte: 1 Espresso Shot
(9, 2, 200),  -- Vanilla Latte: 200 ml Milk
(10, 7, 50);  -- Chocolate Croissant: 50 g Chocolate


-- Mock data for orders
INSERT INTO orders (CustomerName, Notes, CreatedAt) VALUES
('Alice Johnson', '{"notes": "No sugar, extra foam"}', CURRENT_TIMESTAMP),
('Bob Smith', '{"notes": "Add extra espresso shot"}', CURRENT_TIMESTAMP),
('Charlie Brown', '{"notes": "Almond milk instead of regular milk"}', CURRENT_TIMESTAMP),
('David White', '{"notes": "Warm up the milk before adding espresso"}', CURRENT_TIMESTAMP),
('Eve Green', '{"notes": "No whipped cream, add extra chocolate syrup"}', CURRENT_TIMESTAMP),
('Frank Black', '{"notes": "Please make it extra strong"}', CURRENT_TIMESTAMP),
('Grace Blue', '{"notes": "Soy milk and extra sweetener"}', CURRENT_TIMESTAMP),
('Hank Red', '{"notes": "Add a shot of vanilla syrup"}', CURRENT_TIMESTAMP),
('Ivy Gold', '{"notes": "Less ice, more coffee"}', CURRENT_TIMESTAMP),
('Jack Pink', '{"notes": "Double shot of espresso"}', CURRENT_TIMESTAMP);


-- Mock data for order_items
INSERT INTO order_items (OrderID, ProductID, Quantity) VALUES
(1, 1, 1),  -- Alice: 1 Caffe Latte
(1, 2, 2),  -- Alice: 2 Blueberry Muffins
(2, 3, 1),  -- Bob: 1 Espresso
(2, 4, 1),  -- Bob: 1 Cappuccino
(3, 5, 1),  -- Charlie: 1 Mocha
(3, 6, 1),  -- Charlie: 1 Iced Latte
(4, 7, 1),  -- David: 1 Americano
(5, 8, 1),  -- Eve: 1 Carrot Cake
(6, 9, 1),  -- Frank: 1 Vanilla Latte
(7, 10, 2);  -- Grace: 2 Chocolate Croissants

