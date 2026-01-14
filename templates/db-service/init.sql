-- Initialize database schema
-- This file is automatically executed on database container startup

-- Set timezone
SET timezone = 'UTC';

-- Create extension for UUID support (optional, for future use)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The products table will be auto-created by GORM auto-migration
-- This file is mainly for any initial data or custom setup

-- Insert some sample data for testing (optional)
-- Uncomment the lines below to add sample products

-- INSERT INTO products (name, description, price, sku, stock, created_at, updated_at)
-- VALUES 
--   ('Sample Laptop', 'High-performance laptop', 1299.99, 'LAP-SAMPLE-001', 10, NOW(), NOW()),
--   ('Sample Mouse', 'Wireless mouse', 29.99, 'MOU-SAMPLE-001', 50, NOW(), NOW()),
--   ('Sample Keyboard', 'Mechanical keyboard', 89.99, 'KEY-SAMPLE-001', 30, NOW(), NOW());

-- Create indexes for better performance (GORM also creates these via tags)
-- CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_products_sku ON products(sku);

COMMENT ON DATABASE products_db IS 'Database for Product Catalog Service';
