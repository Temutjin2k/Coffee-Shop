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
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
    OpenedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ClosedAt TIMESTAMP,
    FOREIGN KEY (OrderID) REFERENCES orders(ID)
);

CREATE TABLE inventory_transactions (
    IngredientID INT,
    Quantity INT,
    Menu_ItemID INT,
    OrderID INT,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP,
    DeletedAt TIMESTAMP,
    PRIMARY KEY (IngredientID, Menu_ItemID, OrderID, CreatedAt),
    FOREIGN KEY (IngredientID) REFERENCES inventory(IngredientID),
    FOREIGN KEY (Menu_ItemID) REFERENCES menu_items(ID),
    FOREIGN KEY (OrderID) REFERENCES orders(ID)
);


\c frappuccino;