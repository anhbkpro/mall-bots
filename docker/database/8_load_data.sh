#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  -- Stores
  INSERT INTO stores.stores (id, name, location, participating)
    VALUES
      ('bd9e26cd-d861-4b74-862f-06e0825c7af3', 'Waldorf Books', 'West Upper Level', true),
      ('c8b5e3a1-9f2d-4e7c-8b1a-3d4f5e6a7b8c', 'Tech Haven', 'East Ground Floor', true),
      ('d7c6b5a4-3e2f-1d0c-9b8a-7c6d5e4f3a2b', 'Fashion Forward', 'North Wing', true),
      ('e6d5c4b3-2a1f-0e9d-8c7b-6d5e4f3a2b1c', 'Home & Garden', 'South Plaza', false),
      ('f5e4d3c2-1b0a-9f8e-7d6c-5e4f3a2b1c0d', 'Sports Central', 'Central Arena', true);

  -- Products for Waldorf Books
  INSERT INTO stores.products (id, store_id, name, description, sku, price)
    VALUES
      ('557e2f05-b10a-475c-baff-2cea8a86d8df', 'bd9e26cd-d861-4b74-862f-06e0825c7af3', 'EDA with Golang', 'Learn Event-Driven Architecture with Go', 'BOOK-001', 49.99),
      ('668f3g16-c21b-586d-cbgg-3dfb9b97e9eg', 'bd9e26cd-d861-4b74-862f-06e0825c7af3', 'Clean Code', 'A Handbook of Agile Software Craftsmanship', 'BOOK-002', 39.99),
      ('779g4h27-d32c-697e-dchh-4egc0c08f0fh', 'bd9e26cd-d861-4b74-862f-06e0825c7af3', 'Design Patterns', 'Elements of Reusable Object-Oriented Software', 'BOOK-003', 59.99);

  -- Products for Tech Haven
  INSERT INTO stores.products (id, store_id, name, description, sku, price)
    VALUES
      ('88ah5i38-e43d-708f-ecii-5fhd1d19g1gi', 'c8b5e3a1-9f2d-4e7c-8b1a-3d4f5e6a7b8c', 'MacBook Pro 16 inch', 'Latest Apple MacBook Pro with M3 chip', 'TECH-001', 2499.99),
      ('99bi6j49-f54e-819g-fdjj-6gie2e20h2hj', 'c8b5e3a1-9f2d-4e7c-8b1a-3d4f5e6a7b8c', 'iPhone 15 Pro', 'Apple iPhone 15 Pro 256GB', 'TECH-002', 999.99),
      ('aacj7k50-g65f-920h-gekk-7hjf3f21i3ik', 'c8b5e3a1-9f2d-4e7c-8b1a-3d4f5e6a7b8c', 'AirPods Pro', 'Wireless earbuds with noise cancellation', 'TECH-003', 249.99);

  -- Products for Fashion Forward
  INSERT INTO stores.products (id, store_id, name, description, sku, price)
    VALUES
      ('bbdk8l61-h76g-031i-hfll-8ikg4g22j4jl', 'd7c6b5a4-3e2f-1d0c-9b8a-7c6d5e4f3a2b', 'Designer Jeans', 'Premium denim jeans', 'FASH-001', 89.99),
      ('ccel9m62-i87h-142j-igmm-9jlh5h23k5km', 'd7c6b5a4-3e2f-1d0c-9b8a-7c6d5e4f3a2b', 'Leather Jacket', 'Classic leather motorcycle jacket', 'FASH-002', 299.99),
      ('ddfm0n63-j98i-253k-jhnn-0kmi6i24l6ln', 'd7c6b5a4-3e2f-1d0c-9b8a-7c6d5e4f3a2b', 'Running Shoes', 'Professional running shoes', 'FASH-003', 129.99);

  -- Products for Home & Garden
  INSERT INTO stores.products (id, store_id, name, description, sku, price)
    VALUES
      ('eegn1o64-k09j-364l-kioo-1lnj7j25m7mo', 'e6d5c4b3-2a1f-0e9d-8c7b-6d5e4f3a2b1c', 'Garden Tool Set', 'Complete gardening tool kit', 'HOME-001', 79.99),
      ('ffho2p65-l10k-475m-ljpp-2mok8k26n8np', 'e6d5c4b3-2a1f-0e9d-8c7b-6d5e4f3a2b1c', 'Coffee Maker', 'Premium coffee machine', 'HOME-002', 199.99),
      ('ggip3q66-m21l-586n-mkqq-3npl9l27o9oq', 'e6d5c4b3-2a1f-0e9d-8c7b-6d5e4f3a2b1c', 'LED Desk Lamp', 'Modern LED desk lamp', 'HOME-003', 49.99);

  -- Products for Sports Central
  INSERT INTO stores.products (id, store_id, name, description, sku, price)
    VALUES
      ('hhjq4r67-n32m-697o-nlrr-4oqm0m28p0pr', 'f5e4d3c2-1b0a-9f8e-7d6c-5e4f3a2b1c0d', 'Basketball', 'Official NBA basketball', 'SPORT-001', 29.99),
      ('iikr5s68-o43n-708p-omss-5prn1n29q1qs', 'f5e4d3c2-1b0a-9f8e-7d6c-5e4f3a2b1c0d', 'Tennis Racket', 'Professional tennis racket', 'SPORT-002', 159.99),
      ('jjls6t69-p54o-819q-pntt-6qso2o30r2rt', 'f5e4d3c2-1b0a-9f8e-7d6c-5e4f3a2b1c0d', 'Yoga Mat', 'Premium yoga mat', 'SPORT-003', 39.99);

  -- Customers
  INSERT INTO customers.customers (id, name, sms_number, enabled)
    VALUES
      ('f0e2d41a-a485-4008-b578-747732ae1089', 'John Smith', '+1-555-0101', true),
      ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'Sarah Johnson', '+1-555-0102', true),
      ('b2c3d4e5-f6g7-8901-bcde-f23456789012', 'Mike Wilson', '+1-555-0103', true),
      ('c3d4e5f6-g7h8-9012-cdef-345678901234', 'Emily Davis', '+1-555-0104', false),
      ('d4e5f6g7-h8i9-0123-defg-456789012345', 'David Brown', '+1-555-0105', true),
      ('e5f6g7h8-i9j0-1234-efgh-567890123456', 'Lisa Anderson', '+1-555-0106', true);

  -- Sample Baskets (for testing)
  INSERT INTO baskets.baskets (id, customer_id, payment_id, items, status)
    VALUES
      ('basket-001', 'f0e2d41a-a485-4008-b578-747732ae1089', NULL, '[]', 'open'),
      ('basket-002', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', NULL, '[]', 'open'),
      ('basket-003', 'b2c3d4e5-f6g7-8901-bcde-f23456789012', 'payment-001', '[{"product_id": "557e2f05-b10a-475c-baff-2cea8a86d8df", "store_id": "bd9e26cd-d861-4b74-862f-06e0825c7af3", "store_name": "Waldorf Books", "product_name": "EDA with Golang", "quantity": 2, "price": 49.99}]', 'checked_out');

  -- Sample Orders (for testing)
  INSERT INTO ordering.orders (id, customer_id, payment_id, items, status)
    VALUES
      ('order-001', 'f0e2d41a-a485-4008-b578-747732ae1089', 'payment-001', '[{"product_id": "557e2f05-b10a-475c-baff-2cea8a86d8df", "store_id": "bd9e26cd-d861-4b74-862f-06e0825c7af3", "store_name": "Waldorf Books", "product_name": "EDA with Golang", "quantity": 1, "price": 49.99}]', 'confirmed'),
      ('order-002', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'payment-002', '[{"product_id": "88ah5i38-e43d-708f-ecii-5fhd1d19g1gi", "store_id": "c8b5e3a1-9f2d-4e7c-8b1a-3d4f5e6a7b8c", "store_name": "Tech Haven", "product_name": "MacBook Pro 16 inch", "quantity": 1, "price": 2499.99}]', 'pending');

  -- Sample Payments (for testing)
  INSERT INTO payments.payments (id, customer_id, amount, currency, status)
    VALUES
      ('payment-001', 'f0e2d41a-a485-4008-b578-747732ae1089', 49.99, 'USD', 'confirmed'),
      ('payment-002', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 2499.99, 'USD', 'pending'),
      ('payment-003', 'b2c3d4e5-f6g7-8901-bcde-f23456789012', 99.98, 'USD', 'confirmed');

  -- Sample Invoices (for testing)
  INSERT INTO payments.invoices (id, order_id, amount, status)
    VALUES
      ('invoice-001', 'order-001', 49.99, 'paid'),
      ('invoice-002', 'order-002', 2499.99, 'pending');

  -- Sample Depot Items (for testing)
  INSERT INTO depot.shopping_lists (id, order_id, items, status)
    VALUES
      ('list-001', 'order-001', '[{"product_id": "557e2f05-b10a-475c-baff-2cea8a86d8df", "store_id": "bd9e26cd-d861-4b74-862f-06e0825c7af3", "store_name": "Waldorf Books", "product_name": "EDA with Golang", "quantity": 1, "price": 49.99}]', 'ready'),
      ('list-002', 'order-002', '[{"product_id": "88ah5i38-e43d-708f-ecii-5fhd1d19g1gi", "store_id": "c8b5e3a1-9f2d-4e7c-8b1a-3d4f5e6a7b8c", "store_name": "Tech Haven", "product_name": "MacBook Pro 16 inch", "quantity": 1, "price": 2499.99}]', 'pending');

EOSQL
