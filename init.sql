-- Create the database if it doesn't exist
DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'frappuccino') THEN
      EXECUTE 'CREATE DATABASE frappuccino';
   END IF;
END
$$;

-- Connect to the frappuccino database





-- Create the orders table
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
    Unit VARCHAR(10)
);

CREATE TABLE orders (
    ID SERIAL PRIMARY KEY,
    CustomerName VARCHAR(50),
    Status VARCHAR(50),
    CreatedAt DATE
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
    OpenedAt TIMESTAMP DEFAULT NOW(),
    ClosedAt TIMESTAMP,
    FOREIGN KEY (OrderID) REFERENCES orders(ID)
);

CREATE TABLE inventory_transactions (
    IngredientID INT,
    Quantity INT,
    Menu_ItemID INT,
    OrderID INT,
    CreatedAt TIMESTAMP DEFAULT NOW(),
    UpdatedAt TIMESTAMP,
    DeletedAt TIMESTAMP,
    PRIMARY KEY (IngredientID, Menu_ItemID, OrderID, CreatedAt),
    FOREIGN KEY (IngredientID) REFERENCES inventory(IngredientID),
    FOREIGN KEY (Menu_ItemID) REFERENCES menu_items(ID),
    FOREIGN KEY (OrderID) REFERENCES orders(ID)
);



-- Insert into menu_items (Name, Description, Price, Field) values
-- ("Espresso", "Heavy shot of coffee", 5.99, "Idk what to write");

\c frappuccino;