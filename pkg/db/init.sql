-- Auto Parts Inventory Management System Database Schema
-- MVP Version

-- Drop tables if they exist (for clean reinstallation)
DROP TABLE IF EXISTS sales CASCADE;
DROP TABLE IF EXISTS purchases CASCADE;
DROP TABLE IF EXISTS compatibility CASCADE;
DROP TABLE IF EXISTS items CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS vehicle_submodels CASCADE;
DROP TABLE IF EXISTS vehicle_models CASCADE;
DROP TABLE IF EXISTS vehicle_makes CASCADE;
DROP TABLE IF EXISTS suppliers CASCADE;

-- Create extension for UUID generation if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create sequences for IDs
CREATE SEQUENCE IF NOT EXISTS category_id_seq;
CREATE SEQUENCE IF NOT EXISTS make_id_seq;
CREATE SEQUENCE IF NOT EXISTS model_id_seq;
CREATE SEQUENCE IF NOT EXISTS submodel_id_seq;
CREATE SEQUENCE IF NOT EXISTS item_id_seq;
CREATE SEQUENCE IF NOT EXISTS supplier_id_seq;
CREATE SEQUENCE IF NOT EXISTS purchase_id_seq;
CREATE SEQUENCE IF NOT EXISTS sale_id_seq;

-- Categories table with hierarchical structure
CREATE TABLE categories (
    category_id INTEGER PRIMARY KEY DEFAULT nextval('category_id_seq'),
    category_name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_category_id INTEGER REFERENCES categories(category_id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_category_name UNIQUE (category_name)
);

-- Vehicle Makes
CREATE TABLE vehicle_makes (
    make_id INTEGER PRIMARY KEY DEFAULT nextval('make_id_seq'),
    make_name VARCHAR(100) NOT NULL,
    country VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_make_name UNIQUE (make_name)
);

-- Vehicle Models (Parent models like A3, 418, etc.)
CREATE TABLE vehicle_models (
    model_id INTEGER PRIMARY KEY DEFAULT nextval('model_id_seq'),
    make_id INTEGER NOT NULL REFERENCES vehicle_makes(make_id) ON DELETE CASCADE,
    model_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_make_model UNIQUE (make_id, model_name)
);

-- Vehicle Submodels (like A3 Cabrio, 418d, etc.)
CREATE SEQUENCE IF NOT EXISTS submodel_id_seq;
CREATE TABLE vehicle_submodels (
    submodel_id INTEGER PRIMARY KEY DEFAULT nextval('submodel_id_seq'),
    model_id INTEGER NOT NULL REFERENCES vehicle_models(model_id) ON DELETE CASCADE,
    submodel_name VARCHAR(100) NOT NULL,
    year_from INTEGER NOT NULL,
    year_to INTEGER,
    engine_type VARCHAR(100),
    engine_displacement DECIMAL(3,1),
    fuel_type VARCHAR(50),
    transmission_type VARCHAR(50),
    body_type VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_submodel UNIQUE (model_id, submodel_name, year_from)
);

-- Suppliers
CREATE TABLE suppliers (
    supplier_id INTEGER PRIMARY KEY DEFAULT nextval('supplier_id_seq'),
    name VARCHAR(200) NOT NULL,
    contact_person VARCHAR(200),
    phone VARCHAR(50),
    email VARCHAR(200),
    address TEXT,
    tax_id VARCHAR(100),
    payment_terms VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_supplier_name UNIQUE (name)
);

-- Items (Auto Parts)
CREATE TABLE items (
    item_id INTEGER PRIMARY KEY DEFAULT nextval('item_id_seq'),
    part_number VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    category_id INTEGER REFERENCES categories(category_id) ON DELETE SET NULL,
    buy_price DECIMAL(10,2) NOT NULL,
    sell_price DECIMAL(10,2) NOT NULL,
    current_stock INTEGER NOT NULL DEFAULT 0,
    minimum_stock INTEGER NOT NULL DEFAULT 5,
    barcode VARCHAR(100) UNIQUE,
    supplier_id INTEGER REFERENCES suppliers(supplier_id) ON DELETE SET NULL,
    location_aisle VARCHAR(50),
    location_shelf VARCHAR(50),
    location_bin VARCHAR(50),
    weight_kg DECIMAL(10,3),
    dimensions_cm VARCHAR(50), -- Format: LxWxH
    warranty_period VARCHAR(50),
    image_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_part_number UNIQUE (part_number),
    CONSTRAINT positive_buy_price CHECK (buy_price >= 0),
    CONSTRAINT positive_sell_price CHECK (sell_price >= 0),
    CONSTRAINT non_negative_stock CHECK (current_stock >= 0)
);

-- Compatibility mapping between parts and vehicle submodels
CREATE TABLE compatibility (
    compat_id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(item_id) ON DELETE CASCADE,
    submodel_id INTEGER NOT NULL REFERENCES vehicle_submodels(submodel_id) ON DELETE CASCADE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_item_submodel UNIQUE (item_id, submodel_id)
);

-- Purchases
CREATE TABLE purchases (
    purchase_id INTEGER PRIMARY KEY DEFAULT nextval('purchase_id_seq'),
    date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    supplier_id INTEGER NOT NULL REFERENCES suppliers(supplier_id) ON DELETE RESTRICT,
    item_id INTEGER NOT NULL REFERENCES items(item_id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    cost_per_unit DECIMAL(10,2) NOT NULL,
    total_cost DECIMAL(10,2) NOT NULL,
    invoice_number VARCHAR(100),
    received_by VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT positive_cost_per_unit CHECK (cost_per_unit >= 0),
    CONSTRAINT positive_total_cost CHECK (total_cost >= 0)
);

-- Sales
CREATE TABLE sales (
    sale_id INTEGER PRIMARY KEY DEFAULT nextval('sale_id_seq'),
    date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    item_id INTEGER NOT NULL REFERENCES items(item_id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    price_per_unit DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    transaction_number VARCHAR(100),
    customer_name VARCHAR(200),
    customer_phone VARCHAR(50),
    customer_email VARCHAR(200),
    sold_by VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT positive_price_per_unit CHECK (price_per_unit >= 0),
    CONSTRAINT positive_total_price CHECK (total_price >= 0)
);

-- Create indexes for performance
CREATE INDEX idx_categories_parent ON categories(parent_category_id);
CREATE INDEX idx_vehicle_models_make ON vehicle_models(make_id);
CREATE INDEX idx_items_category ON items(category_id);
CREATE INDEX idx_items_supplier ON items(supplier_id);
CREATE INDEX idx_items_barcode ON items(barcode);
CREATE INDEX idx_compatibility_item ON compatibility(item_id);
CREATE INDEX idx_compatibility_submodel ON compatibility(submodel_id);
CREATE INDEX idx_purchases_supplier ON purchases(supplier_id);
CREATE INDEX idx_purchases_item ON purchases(item_id);
CREATE INDEX idx_purchases_date ON purchases(date);
CREATE INDEX idx_sales_item ON sales(item_id);
CREATE INDEX idx_sales_date ON sales(date);

-- Create triggers for updated_at timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_categories_timestamp
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_vehicle_makes_timestamp
BEFORE UPDATE ON vehicle_makes
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_vehicle_models_timestamp
BEFORE UPDATE ON vehicle_models
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_vehicle_submodels_timestamp
BEFORE UPDATE ON vehicle_submodels
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_suppliers_timestamp
BEFORE UPDATE ON suppliers
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_items_timestamp
BEFORE UPDATE ON items
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_purchases_timestamp
BEFORE UPDATE ON purchases
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE TRIGGER update_sales_timestamp
BEFORE UPDATE ON sales
FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

-- Create a trigger to update inventory on purchase
CREATE OR REPLACE FUNCTION update_inventory_on_purchase()
RETURNS TRIGGER AS $$
BEGIN
   UPDATE items
   SET current_stock = current_stock + NEW.quantity,
       updated_at = CURRENT_TIMESTAMP
   WHERE item_id = NEW.item_id;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_inventory_on_purchase
AFTER INSERT ON purchases
FOR EACH ROW EXECUTE PROCEDURE update_inventory_on_purchase();

-- Create a trigger to update inventory on sale
CREATE OR REPLACE FUNCTION update_inventory_on_sale()
RETURNS TRIGGER AS $$
BEGIN
   UPDATE items
   SET current_stock = current_stock - NEW.quantity,
       updated_at = CURRENT_TIMESTAMP
   WHERE item_id = NEW.item_id;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_inventory_on_sale
AFTER INSERT ON sales
FOR EACH ROW EXECUTE PROCEDURE update_inventory_on_sale();

-- Insert some sample data for categories
INSERT INTO categories (category_name, description) VALUES
('Engine Parts', 'Parts related to the engine system'),
('Brake System', 'Parts related to the braking system'),
('Suspension', 'Parts related to the suspension system'),
('Electrical', 'Electrical components and systems'),
('Body Parts', 'External and structural body components'),
('Filters', 'All types of filters for vehicles');

-- Insert some child categories
INSERT INTO categories (category_name, description, parent_category_id) VALUES
('Pistons', 'Engine pistons and related components', 1),
('Timing Belts', 'Engine timing belts and chains', 1),
('Brake Pads', 'Friction material for brake systems', 2),
('Brake Rotors', 'Rotating discs for brake systems', 2),
('Shock Absorbers', 'Dampers for suspension systems', 3),
('Spring Coils', 'Suspension springs and coils', 3),
('Headlights', 'Front lighting systems', 4),
('Alternators', 'Charging system components', 4),
('Hood', 'Front cover for engine compartment', 5),
('Bumpers', 'Front and rear impact protection', 5),
('Oil Filters', 'Filtration for engine oil', 6),
('Air Filters', 'Filtration for engine air intake', 6);

-- Insert some sample vehicle makes
INSERT INTO vehicle_makes (make_name, country) VALUES
('Toyota', 'Japan'),
('Honda', 'Japan'),
('Ford', 'USA'),
('Chevrolet', 'USA'),
('BMW', 'Germany'),
('Mercedes-Benz', 'Germany'),
('Volkswagen', 'Germany'),
('Hyundai', 'South Korea'),
('Audi', 'Germany'),
('Nissan', 'Japan');

-- Insert sample vehicle models (parent models)
INSERT INTO vehicle_models (make_id, model_name) VALUES
(1, 'Camry'),
(1, 'Corolla'),
(2, 'Civic'),
(2, 'Accord'),
(3, 'F-150'),
(3, 'Mustang'),
(4, 'Silverado'),
(5, '3 Series'),
(5, '4 Series'),  -- Adding BMW 4 Series
(6, 'C-Class'),   -- Adding Mercedes C-Class
(7, 'Golf'),      -- Adding VW Golf
(7, 'Passat');    -- Adding VW Passat

-- Insert sample vehicle submodels
INSERT INTO vehicle_submodels (model_id, submodel_name, year_from, year_to, engine_type, engine_displacement, fuel_type, transmission_type, body_type) VALUES
-- Toyota Camry variants
(1, 'Camry SE', 2018, 2022, 'Inline-4', 2.5, 'Gasoline', 'Automatic', 'Sedan'),
(1, 'Camry XLE', 2018, 2022, 'Inline-4', 2.5, 'Gasoline', 'Automatic', 'Sedan'),
(1, 'Camry Hybrid', 2018, 2022, 'Hybrid Inline-4', 2.5, 'Hybrid', 'CVT', 'Sedan'),

-- Toyota Corolla variants
(2, 'Corolla LE', 2019, 2023, 'Inline-4', 1.8, 'Gasoline', 'CVT', 'Sedan'),
(2, 'Corolla Hatchback', 2019, 2023, 'Inline-4', 2.0, 'Gasoline', 'CVT', 'Hatchback'),
(2, 'Corolla Hybrid', 2020, 2023, 'Hybrid Inline-4', 1.8, 'Hybrid', 'CVT', 'Sedan'),

-- Honda Civic variants
(3, 'Civic Sedan', 2016, 2021, 'Inline-4', 1.5, 'Gasoline', 'CVT', 'Sedan'),
(3, 'Civic Hatchback', 2017, 2021, 'Inline-4', 1.5, 'Gasoline', 'Manual', 'Hatchback'),
(3, 'Civic Type R', 2017, 2021, 'Inline-4 Turbo', 2.0, 'Gasoline', 'Manual', 'Hatchback'),

-- Honda Accord variants
(4, 'Accord Sport', 2018, 2022, 'Inline-4 Turbo', 1.5, 'Gasoline', 'CVT', 'Sedan'),
(4, 'Accord Touring', 2018, 2022, 'Inline-4 Turbo', 2.0, 'Gasoline', 'Automatic', 'Sedan'),
(4, 'Accord Hybrid', 2018, 2022, 'Hybrid Inline-4', 2.0, 'Hybrid', 'eCVT', 'Sedan'),

-- Ford F-150 variants
(5, 'F-150 XLT', 2015, 2020, 'V8', 5.0, 'Gasoline', 'Automatic', 'Pickup Truck'),
(5, 'F-150 Raptor', 2017, 2020, 'V6 Turbo', 3.5, 'Gasoline', 'Automatic', 'Pickup Truck'),
(5, 'F-150 Lariat', 2015, 2020, 'V6', 3.5, 'Gasoline', 'Automatic', 'Pickup Truck'),

-- Ford Mustang variants
(6, 'Mustang GT', 2018, 2023, 'V8', 5.0, 'Gasoline', 'Manual', 'Coupe'),
(6, 'Mustang EcoBoost', 2018, 2023, 'Inline-4 Turbo', 2.3, 'Gasoline', 'Automatic', 'Coupe'),
(6, 'Mustang Convertible', 2018, 2023, 'V8', 5.0, 'Gasoline', 'Automatic', 'Convertible'),

-- Chevy Silverado variants
(7, 'Silverado LT', 2019, 2023, 'V8', 5.3, 'Gasoline', 'Automatic', 'Pickup Truck'),
(7, 'Silverado Custom', 2019, 2023, 'V6', 4.3, 'Gasoline', 'Automatic', 'Pickup Truck'),
(7, 'Silverado RST', 2019, 2023, 'V8', 6.2, 'Gasoline', 'Automatic', 'Pickup Truck'),

-- BMW 3 Series variants
(8, '330i', 2019, 2023, 'Inline-4 Turbo', 2.0, 'Gasoline', 'Automatic', 'Sedan'),
(8, '330i xDrive', 2019, 2023, 'Inline-4 Turbo', 2.0, 'Gasoline', 'Automatic', 'Sedan'),
(8, 'M340i', 2019, 2023, 'Inline-6 Turbo', 3.0, 'Gasoline', 'Automatic', 'Sedan'),

-- BMW 4 Series variants
(9, '418d', 2014, 2020, 'Inline-4 Diesel', 2.0, 'Diesel', 'Automatic', 'Coupe'),
(9, '418d Gran Coupe', 2014, 2020, 'Inline-4 Diesel', 2.0, 'Diesel', 'Automatic', 'Gran Coupe'),
(9, '418i', 2014, 2020, 'Inline-4 Turbo', 2.0, 'Gasoline', 'Automatic', 'Coupe'),
(9, '418i Gran Coupe', 2014, 2020, 'Inline-4 Turbo', 2.0, 'Gasoline', 'Automatic', 'Gran Coupe'),

-- Audi A3 variants (adding these as you specifically mentioned them)
(10, 'A3 Cabrio', 2016, 2020, 'Inline-4 Turbo', 1.4, 'Gasoline', 'Automatic', 'Convertible'),
(10, 'A3 Sedan', 2016, 2020, 'Inline-4 Turbo', 1.4, 'Gasoline', 'Automatic', 'Sedan'),
(10, 'A3 Sportback', 2016, 2020, 'Inline-4 Turbo', 1.4, 'Gasoline', 'Automatic', 'Sportback'),
(10, 'A3 Hatchback', 2016, 2020, 'Inline-4 Turbo', 1.4, 'Gasoline', 'Manual', 'Hatchback');

-- Insert some sample suppliers
INSERT INTO suppliers (name, contact_person, phone, email, address, tax_id, payment_terms) VALUES
('Auto Parts Wholesale Inc.', 'John Smith', '555-123-4567', 'john@apw.com', '123 Main St, Anytown, USA', 'APW-12345', 'Net 30'),
('Quality Parts Supply', 'Jane Doe', '555-234-5678', 'jane@qps.com', '456 Second Ave, Othertown, USA', 'QPS-67890', 'Net 45'),
('Import Auto Parts', 'Bob Johnson', '555-345-6789', 'bob@importauto.com', '789 Third Blvd, Somewhere, USA', 'IAP-24680', 'COD'),
('OEM Suppliers Ltd.', 'Mary Wilson', '555-456-7890', 'mary@oemsuppliers.com', '321 Fourth St, Elsewhere, USA', 'OEM-13579', 'Net 60');

-- Insert some sample items
INSERT INTO items (part_number, description, category_id, buy_price, sell_price, current_stock, minimum_stock, barcode, supplier_id, location_aisle, location_shelf, location_bin) VALUES
('BP-1234', 'Premium Brake Pads - Front', 3, 25.50, 49.99, 45, 10, 'BP1234FRONT', 1, 'A', '1', '3'),
('BP-1235', 'Premium Brake Pads - Rear', 3, 22.75, 45.99, 38, 10, 'BP1235REAR', 1, 'A', '1', '4'),
('OF-2345', 'Oil Filter - Standard', 11, 3.25, 8.99, 120, 30, 'OF2345STD', 2, 'B', '3', '1'),
('AF-3456', 'Air Filter - Performance', 12, 12.50, 24.99, 35, 15, 'AF3456PERF', 3, 'B', '3', '5'),
('SA-4567', 'Shock Absorber - Front', 5, 45.75, 89.99, 18, 8, 'SA4567FRONT', 4, 'C', '2', '2'),
('SC-5678', 'Spring Coils - Lowering Kit', 6, 120.00, 249.99, 7, 4, 'SC5678LOWER', 3, 'C', '2', '6'),
('HL-6789', 'Headlight Assembly - Left', 7, 85.50, 169.99, 12, 6, 'HL6789LEFT', 2, 'D', '1', '1'),
('HL-6790', 'Headlight Assembly - Right', 7, 85.50, 169.99, 11, 6, 'HL6790RIGHT', 2, 'D', '1', '2'),
('AL-7890', 'Alternator - 120A', 8, 65.25, 129.99, 9, 5, 'AL7890120A', 1, 'D', '2', '4'),
('TB-8901', 'Timing Belt Kit', 2, 48.75, 94.99, 22, 8, 'TB8901KIT', 4, 'A', '3', '2');

-- Insert some sample compatibility records
INSERT INTO compatibility (item_id, submodel_id, notes) VALUES
-- Front brake pads
(1, 1, 'Perfect fit for Camry SE'), -- Brake pads for Toyota Camry SE
(1, 4, 'Works with minor modification'), -- Brake pads for Toyota Corolla LE
(1, 7, 'Direct replacement'), -- Brake pads for Honda Civic Sedan

-- Rear brake pads
(2, 1, 'OEM replacement'), -- Rear brake pads for Toyota Camry SE
(2, 4, 'OEM replacement'), -- Rear brake pads for Toyota Corolla LE

-- Oil filters
(3, 1, 'Standard oil filter'), -- Oil filter for Toyota Camry SE
(3, 2, 'Standard oil filter'), -- Oil filter for Toyota Camry XLE
(3, 3, 'Standard oil filter'), -- Oil filter for Toyota Camry Hybrid
(3, 4, 'Standard oil filter'), -- Oil filter for Toyota Corolla LE
(3, 7, 'Standard oil filter'), -- Oil filter for Honda Civic Sedan
(3, 10, 'Standard oil filter'), -- Oil filter for Honda Accord Sport

-- Air filters
(4, 13, 'Performance upgrade'), -- Air filter for Ford F-150 XLT
(4, 19, 'Performance upgrade'), -- Air filter for Chevy Silverado LT

-- Shock absorbers
(5, 1, 'OEM replacement'), -- Shock absorber for Toyota Camry SE
(5, 10, 'OEM replacement'), -- Shock absorber for Honda Accord Sport

-- Headlights
(7, 7, 'Direct replacement'), -- Headlight for Honda Civic Sedan
(8, 7, 'Direct replacement'), -- Right headlight for Honda Civic Sedan

-- Alternators
(9, 13, 'High output replacement'), -- Alternator for Ford F-150 XLT
(9, 19, 'High output replacement'), -- Alternator for Chevy Silverado LT

-- Timing belts
(10, 1, 'OEM quality replacement'), -- Timing belt for Toyota Camry SE
(10, 4, 'OEM quality replacement'), -- Timing belt for Toyota Corolla LE

-- Additional compatibility for BMW 4 Series
(3, 25, 'Compatible with all diesel variants'), -- Oil filter for BMW 418d
(3, 26, 'Compatible with all diesel variants'), -- Oil filter for BMW 418d Gran Coupe
(3, 27, 'Compatible with all gasoline variants'), -- Oil filter for BMW 418i
(3, 28, 'Compatible with all gasoline variants'), -- Oil filter for BMW 418i Gran Coupe

-- Additional compatibility for Audi A3 variants
(3, 29, 'Compatible with all engines'), -- Oil filter for Audi A3 Cabrio
(3, 30, 'Compatible with all engines'), -- Oil filter for Audi A3 Sedan
(3, 31, 'Compatible with all engines'), -- Oil filter for Audi A3 Sportback
(3, 32, 'Compatible with all engines'); -- Oil filter for Audi A3 Hatchback

-- Insert some purchase records
INSERT INTO purchases (supplier_id, item_id, quantity, cost_per_unit, total_cost, invoice_number, received_by) VALUES
(1, 1, 20, 25.50, 510.00, 'INV-2023-001', 'Mike Johnson'),
(1, 2, 15, 22.75, 341.25, 'INV-2023-001', 'Mike Johnson'),
(2, 3, 50, 3.25, 162.50, 'INV-2023-002', 'Sarah Williams'),
(3, 4, 15, 12.50, 187.50, 'INV-2023-003', 'David Brown'),
(4, 5, 10, 45.75, 457.50, 'INV-2023-004', 'Lisa Davis'),
(3, 6, 5, 120.00, 600.00, 'INV-2023-005', 'Robert Wilson'),
(2, 7, 6, 85.50, 513.00, 'INV-2023-006', 'Jennifer Taylor'),
(2, 8, 6, 85.50, 513.00, 'INV-2023-006', 'Jennifer Taylor'),
(1, 9, 4, 65.25, 261.00, 'INV-2023-007', 'Michael Moore'),
(4, 10, 12, 48.75, 585.00, 'INV-2023-008', 'Patricia Martin');

-- Insert some sales records
INSERT INTO sales (item_id, quantity, price_per_unit, total_price, transaction_number, customer_name, customer_phone, sold_by) VALUES
(1, 2, 49.99, 99.98, 'TRX-2023-001', 'James Wilson', '555-111-2222', 'Tom Baker'),
(3, 1, 8.99, 8.99, 'TRX-2023-002', 'Maria Garcia', '555-222-3333', 'Tom Baker'),
(4, 1, 24.99, 24.99, 'TRX-2023-002', 'Maria Garcia', '555-222-3333', 'Tom Baker'),
(7, 1, 169.99, 169.99, 'TRX-2023-003', 'Robert Johnson', '555-333-4444', 'Alice Cooper'),
(8, 1, 169.99, 169.99, 'TRX-2023-003', 'Robert Johnson', '555-333-4444', 'Alice Cooper'),
(2, 1, 45.99, 45.99, 'TRX-2023-004', 'Susan Miller', '555-444-5555', 'Tom Baker'),
(5, 2, 89.99, 179.98, 'TRX-2023-005', 'David Thompson', '555-555-6666', 'Alice Cooper'),
(10, 1, 94.99, 94.99, 'TRX-2023-006', 'Linda Martinez', '555-666-7777', 'Tom Baker'),
(3, 3, 8.99, 26.97, 'TRX-2023-007', 'Michael Brown', '555-777-8888', 'Alice Cooper'),
(6, 1, 249.99, 249.99, 'TRX-2023-008', 'Jennifer Davis', '555-888-9999', 'Tom Baker');

-- Create view for low stock alerts
CREATE OR REPLACE VIEW low_stock_items AS
SELECT
    i.item_id,
    i.part_number,
    i.description,
    i.current_stock,
    i.minimum_stock,
    c.category_name,
    s.name as supplier_name,
    s.phone as supplier_phone,
    s.email as supplier_email
FROM
    items i
JOIN
    categories c ON i.category_id = c.category_id
JOIN
    suppliers s ON i.supplier_id = s.supplier_id
WHERE
    i.current_stock <= i.minimum_stock
ORDER BY
    (i.current_stock::float / i.minimum_stock) ASC;

-- Create view for item sales velocity (items sold per day)
CREATE OR REPLACE VIEW item_sales_velocity AS
SELECT
    i.item_id,
    i.part_number,
    i.description,
    COUNT(s.sale_id) as total_sales,
    SUM(s.quantity) as total_quantity_sold,
    (CURRENT_DATE - DATE '2023-01-01') as days_since_jan1,
    ROUND(SUM(s.quantity)::numeric / ((CURRENT_DATE - DATE '2023-01-01')::numeric), 2) as daily_sales_rate,
    CASE
        WHEN i.current_stock > 0 AND (SUM(s.quantity)::numeric / ((CURRENT_DATE - DATE '2023-01-01')::numeric)) > 0
        THEN ROUND(i.current_stock / (SUM(s.quantity)::numeric / ((CURRENT_DATE - DATE '2023-01-01')::numeric)))
        ELSE NULL
    END as estimated_days_until_stockout
FROM
    items i
LEFT JOIN
    sales s ON i.item_id = s.item_id
WHERE
    s.date >= '2023-01-01'
GROUP BY
    i.item_id, i.part_number, i.description, i.current_stock
ORDER BY
    daily_sales_rate DESC;

-- Create view for top selling items
CREATE OR REPLACE VIEW top_selling_items AS
SELECT
    i.item_id,
    i.part_number,
    i.description,
    c.category_name,
    COUNT(s.sale_id) as number_of_sales,
    SUM(s.quantity) as total_quantity_sold,
    SUM(s.total_price) as total_revenue,
    SUM(s.total_price) - (SUM(s.quantity) * i.buy_price) as estimated_profit
FROM
    items i
JOIN
    sales s ON i.item_id = s.item_id
JOIN
    categories c ON i.category_id = c.category_id
GROUP BY
    i.item_id, i.part_number, i.description, i.buy_price, c.category_name
ORDER BY
    total_revenue DESC;

-- Create view for vehicle compatibility count
CREATE OR REPLACE VIEW part_compatibility_summary AS
SELECT
    i.item_id,
    i.part_number,
    i.description,
    COUNT(c.compat_id) as compatible_vehicle_count,
    COUNT(DISTINCT vmod.model_id) as compatible_model_count,
    STRING_AGG(DISTINCT vm.make_name, ', ') as compatible_makes
FROM
    items i
LEFT JOIN
    compatibility c ON i.item_id = c.item_id
LEFT JOIN
    vehicle_submodels vsub ON c.submodel_id = vsub.submodel_id
LEFT JOIN
    vehicle_models vmod ON vsub.model_id = vmod.model_id
LEFT JOIN
    vehicle_makes vm ON vmod.make_id = vm.make_id
GROUP BY
    i.item_id, i.part_number, i.description
ORDER BY
    compatible_vehicle_count DESC;

-- Create function to get compatible parts for a specific vehicle submodel
CREATE OR REPLACE FUNCTION get_compatible_parts(make_name TEXT, model_name TEXT, submodel_name TEXT, model_year INTEGER)
RETURNS TABLE (
    item_id INTEGER,
    part_number VARCHAR,
    description TEXT,
    category_name VARCHAR,
    sell_price DECIMAL,
    current_stock INTEGER,
    barcode VARCHAR
) AS $
BEGIN
    RETURN QUERY
    SELECT
        i.item_id,
        i.part_number,
        i.description,
        c.category_name,
        i.sell_price,
        i.current_stock,
        i.barcode
    FROM
        items i
    JOIN
        categories c ON i.category_id = c.category_id
    JOIN
        compatibility comp ON i.item_id = comp.item_id
    JOIN
        vehicle_submodels vsub ON comp.submodel_id = vsub.submodel_id
    JOIN
        vehicle_models vm ON vsub.model_id = vm.model_id
    JOIN
        vehicle_makes vma ON vm.make_id = vma.make_id
    WHERE
        vma.make_name = make_name
        AND vm.model_name = model_name
        AND vsub.submodel_name = submodel_name
        AND vsub.year_from <= model_year
        AND (vsub.year_to IS NULL OR vsub.year_to >= model_year)
    ORDER BY
        c.category_name, i.part_number;
END;
$ LANGUAGE plpgsql;

-- Create function to get all compatible parts for a model (across all submodels)
CREATE OR REPLACE FUNCTION get_model_compatible_parts(make_name TEXT, model_name TEXT)
RETURNS TABLE (
    item_id INTEGER,
    part_number VARCHAR,
    description TEXT,
    category_name VARCHAR,
    sell_price DECIMAL,
    current_stock INTEGER,
    compatible_submodels TEXT
) AS $
BEGIN
    RETURN QUERY
    SELECT
        i.item_id,
        i.part_number,
        i.description,
        c.category_name,
        i.sell_price,
        i.current_stock,
        STRING_AGG(DISTINCT vsub.submodel_name, ', ') as compatible_submodels
    FROM
        items i
    JOIN
        categories c ON i.category_id = c.category_id
    JOIN
        compatibility comp ON i.item_id = comp.item_id
    JOIN
        vehicle_submodels vsub ON comp.submodel_id = vsub.submodel_id
    JOIN
        vehicle_models vm ON vsub.model_id = vm.model_id
    JOIN
        vehicle_makes vma ON vm.make_id = vma.make_id
    WHERE
        vma.make_name = make_name
        AND vm.model_name = model_name
    GROUP BY
        i.item_id, i.part_number, i.description, c.category_name, i.sell_price, i.current_stock
    ORDER BY
        c.category_name, i.part_number;
END;
$ LANGUAGE plpgsql;

-- Function to get category prefix
CREATE OR REPLACE FUNCTION get_category_prefix(category_id INTEGER)
RETURNS TEXT AS $$
DECLARE
    category_name TEXT;
BEGIN
    SELECT UPPER(LEFT(category_name, 2))
    INTO category_name
    FROM categories
    WHERE category_id = $1;
    RETURN category_name;
END;
$$ LANGUAGE plpgsql;

-- Function to generate barcode
CREATE OR REPLACE FUNCTION generate_barcode(category_id INTEGER, item_id INTEGER)
RETURNS TEXT AS $$
DECLARE
    category_prefix TEXT;
    item_number TEXT;
    year_suffix TEXT;
    base_code TEXT;
    check_digit INTEGER;
BEGIN
    -- Get category prefix
    category_prefix := get_category_prefix(category_id);

    -- Format item_id to 6 digits
    item_number := LPAD(item_id::TEXT, 6, '0');

    -- Get last 2 digits of current year
    year_suffix := RIGHT(EXTRACT(YEAR FROM CURRENT_DATE)::TEXT, 2);

    -- Create base code
    base_code := category_prefix || '-' || item_number || '-' || year_suffix;

    -- Calculate check digit (sum of ASCII values mod 10)
    SELECT SUM(ASCII(c)) % 10
    INTO check_digit
    FROM regexp_split_to_table(base_code, '') c;

    -- Return complete barcode
    RETURN base_code || '-' || check_digit;
END;
$$ LANGUAGE plpgsql;

-- Trigger to generate barcode before insert
CREATE OR REPLACE FUNCTION generate_barcode_trigger()
RETURNS TRIGGER AS $$
BEGIN
    -- Only generate barcode if one isn't provided
    IF NEW.barcode IS NULL THEN
        NEW.barcode := generate_barcode(NEW.category_id, NEW.item_id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
DROP TRIGGER IF EXISTS trg_generate_barcode ON items;
CREATE TRIGGER trg_generate_barcode
    BEFORE INSERT ON items
    FOR EACH ROW
    EXECUTE FUNCTION generate_barcode_trigger();
