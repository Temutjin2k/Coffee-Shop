-- Create the database if it doesn't exist
DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'frappuccino') THEN
      EXECUTE 'CREATE DATABASE frappuccino';
   END IF;
END
$$;

-- Connect to the frappuccino database




\c frappuccino;

-- Create ENUM types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        CREATE TYPE order_status AS ENUM ('Pending', 'Preparing', 'Completed', 'Cancelled');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_method') THEN
        CREATE TYPE payment_method AS ENUM ('Cash', 'Card', 'Online');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'item_size') THEN
        CREATE TYPE item_size AS ENUM ('Small', 'Medium', 'Large');
    END IF;
END $$;

-- Create the orders table
CREATE TABLE orders (
    ID SERIAL PRIMARY KEY,
    CustomerName VARCHAR(50),
    Status order_status,
    SpecialInstructions JSONB,
    CreatedAt TIMESTAMPTZ DEFAULT NOW()
);

-- Create the order_items table
CREATE TABLE order_items (
    OrderID INT,
    ProductID INT,
    Quantity INT,
    CustomizationOptions JSONB,
    PRIMARY KEY (OrderID, ProductID),
    FOREIGN KEY (OrderID) REFERENCES orders(ID),
    FOREIGN KEY (ProductID) REFERENCES menu_items(ID)
);

-- Create the menu_items table
CREATE TABLE menu_items (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(50),
    Description TEXT,
    Price NUMERIC(10, 2),
    Field VARCHAR(50),
    Tags TEXT[] DEFAULT '{}',
    Allergens TEXT[] DEFAULT '{}',
    Metadata JSONB
);

-- Create the inventory table
CREATE TABLE inventory (
    IngredientID SERIAL PRIMARY KEY,
    Name VARCHAR(50),
    Quantity INT,
    Unit VARCHAR(10),
    Substitutes TEXT[] DEFAULT '{}'
);

-- Create the price_history table
CREATE TABLE price_history (
    Menu_ItemID INT,
    Price NUMERIC(10, 2),
    Date DATE,
    PRIMARY KEY (Menu_ItemID, Date),
    FOREIGN KEY (Menu_ItemID) REFERENCES menu_items(ID)
);

-- Create the menu_item_ingredients table
CREATE TABLE menu_item_ingredients (
    MenuID INT,
    IngredientID INT,
    Quantity INT,
    PRIMARY KEY (MenuID, IngredientID),
    FOREIGN KEY (MenuID) REFERENCES menu_items(ID),
    FOREIGN KEY (IngredientID) REFERENCES inventory(IngredientID)
);

-- Create the order_status_history table
CREATE TABLE order_status_history (
    ID SERIAL PRIMARY KEY,
    OrderID INT,
    StatusChanges TEXT[],
    OpenedAt TIMESTAMPTZ,
    ClosedAt TIMESTAMPTZ,
    FOREIGN KEY (OrderID) REFERENCES orders(ID)
);

-- Create the inventory_transactions table
CREATE TABLE inventory_transactions (
    IngredientID INT,
    Quantity INT,
    Menu_ItemID INT,
    OrderID INT,
    CreatedAt TIMESTAMPTZ DEFAULT NOW(),
    UpdatedAt TIMESTAMPTZ,
    DeletedAt TIMESTAMPTZ,
    PRIMARY KEY (IngredientID, Menu_ItemID, OrderID, CreatedAt),
    FOREIGN KEY (IngredientID) REFERENCES inventory(IngredientID),
    FOREIGN KEY (Menu_ItemID) REFERENCES menu_items(ID),
    FOREIGN KEY (OrderID) REFERENCES orders(ID)
);
