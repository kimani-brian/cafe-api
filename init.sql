-- 1. Create Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'cashier'))
);

-- 2. Create Inventory Items Table (The Menu)
CREATE TABLE inventory_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INT NOT NULL CHECK (stock_quantity >= 0)
);

-- 3. Create Orders Table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    cashier_id INT REFERENCES users(id),
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create Order Items Table (The Receipt/Line Items)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    item_id INT REFERENCES inventory_items(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    price_at_purchase DECIMAL(10, 2) NOT NULL
);

-- 5. Insert a Default Admin User 
-- (Password is 'admin123', hashed using bcrypt for security)
INSERT INTO users (username, password_hash, role) 
VALUES ('manager', '$2a$10$Y1/n2p0O237.Z8o/tZqX1.B0q88.vTf8U1/2L6R6Z3/U1n.u.wY3O', 'admin');

-- 6. Insert a Default Cashier User 
-- (Password is 'cashier123')
INSERT INTO users (username, password_hash, role) 
VALUES ('barista_bob', '$2a$10$tZqX1.B0q88.vTf8U1/2L6R6Z3/U1n.u.wY3OY1/n2p0O237.Z8o/t', 'cashier');