# MallBots API Testing Guide

A comprehensive guide for testing the MallBots e-commerce platform's event-driven architecture using REST API endpoints.

## 🏗️ Architecture Overview

MallBots implements a **monolithic modular architecture** with the following services:

- **🛒 Baskets** - Shopping cart management (`/api/baskets`)
- **👥 Customers** - Customer management (`/api/customers`)
- **📦 Depot** - Inventory and warehouse operations
- **🔔 Notifications** - User notifications
- **📋 Ordering** - Order processing (`/api/ordering`)
- **💳 Payments** - Payment processing (`/api/payments`)
- **🏪 Stores** - Store management (`/api/stores`)

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make (optional, for convenience)

### 1. Install Development Tools

```bash
make install-tools
```

### 2. Generate Code

```bash
make generate
```

### 3. Start the Application

```bash
# Start all services
make up

# Wait for services to be ready (check logs)
make logs
```

### 4. Access Points

- **Web UI**: `http://localhost:8080`
- **gRPC Server**: localhost:8085
- **PostgreSQL**: localhost:5432

## 🧪 Complete E-commerce Flow Testing

### Step 1: Create a Customer

```bash
curl -X 'POST' \
  'http://localhost:8080/api/customers' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "John Doe",
  "smsNumber": "+1234567890"
}'
```

**Expected Response:**
```json
{
  "id": "27dc567d-524d-450b-beae-70e5628d0de8"
}
```

### Step 2: Create a Store

```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "Tech Store",
  "location": "San Francisco, CA"
}'
```

**Expected Response:**
```json
{
  "id": "ca1d1edc-0c9d-4aec-ac51-b9ee28c0c484"
}
```

### Step 3: Enable Store Participation

```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/{store-id}/enable-participation' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

### Step 4: Add Products to Store

```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/{store-id}/products' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "iPhone 15",
  "description": "Latest iPhone model",
  "sku": "IPHONE-15-128GB",
  "price": 999.99
}'
```

**Expected Response:**
```json
{
  "id": "product-uuid"
}
```

Add another product:

```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/{store-id}/products' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "MacBook Pro",
  "description": "Professional laptop",
  "sku": "MBP-14-M2",
  "price": 1999.99
}'
```

### Step 5: Start Shopping Basket

```bash
curl -X 'POST' \
  'http://localhost:8080/api/baskets' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "customerId": "{customer-id}"
}'
```

**Expected Response:**
```json
{
  "id": "basket-uuid"
}
```

### Step 6: Add Items to Basket

```bash
curl -X 'PUT' \
  'http://localhost:8080/api/baskets/{basket-id}/addItem' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "productId": "{product-id}",
  "quantity": 1
}'
```

### Step 7: Get Basket Contents

```bash
curl -X 'GET' \
  'http://localhost:8080/api/baskets/{basket-id}' \
  -H 'accept: application/json'
```

**Expected Response:**
```json
{
  "basket": {
    "id": "basket-uuid",
    "items": [
      {
        "store_id": "store-uuid",
        "product_id": "product-uuid",
        "store_name": "Tech Store",
        "product_name": "iPhone 15",
        "product_price": 999.99,
        "quantity": 1
      }
    ]
  }
}
```

### Step 8: Authorize Payment

```bash
curl -X 'POST' \
  'http://localhost:8080/api/payments' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "customerId": "{customer-id}",
  "amount": 999.99
}'
```

**Expected Response:**
```json
{
  "id": "payment-uuid"
}
```

### Step 9: Checkout Basket

```bash
curl -X 'PUT' \
  'http://localhost:8080/api/baskets/{basket-id}/checkout' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "paymentId": "{payment-id}"
}'
```

### Step 10: Create Order

```bash
curl -X 'POST' \
  'http://localhost:8080/api/ordering' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "items": [
    {
      "storeId": "{store-id}",
      "productId": "{product-id}",
      "storeName": "Tech Store",
      "productName": "iPhone 15",
      "price": 999.99,
      "quantity": 1
    }
  ],
  "customerId": "{customer-id}",
  "paymentId": "{payment-id}"
}'
```

**Expected Response:**
```json
{
  "id": "order-uuid"
}
```

### Step 11: Confirm Payment

```bash
curl -X 'POST' \
  'http://localhost:8080/api/payments/{payment-id}/confirm' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

### Step 12: Create Invoice

```bash
curl -X 'POST' \
  'http://localhost:8080/api/payments/invoices' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "orderId": "{order-id}",
  "paymentId": "{payment-id}",
  "amount": 999.99
}'
```

**Expected Response:**
```json
{
  "id": "invoice-uuid"
}
```

### Step 13: Pay Invoice

```bash
curl -X 'POST' \
  'http://localhost:8080/api/payments/invoices/{invoice-id}/pay' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

### Step 14: Complete Order

```bash
curl -X 'POST' \
  'http://localhost:8080/api/ordering/orders/{order-id}/complete' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "invoiceId": "{invoice-id}"
}'
```

## 📋 Additional API Endpoints

### Customer Management

#### Get Customer Details
```bash
curl -X 'GET' \
  'http://localhost:8080/api/customers/{customer-id}' \
  -H 'accept: application/json'
```

#### Enable Customer
```bash
curl -X 'POST' \
  'http://localhost:8080/api/customers/{customer-id}/enable' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

#### Authorize Customer
```bash
curl -X 'POST' \
  'http://localhost:8080/api/customers/{customer-id}/authorize' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

### Store Management

#### Get Store Details
```bash
curl -X 'GET' \
  'http://localhost:8080/api/stores/{store-id}' \
  -H 'accept: application/json'
```

#### Get Store Catalog
```bash
curl -X 'GET' \
  'http://localhost:8080/api/stores/{store-id}/catalog' \
  -H 'accept: application/json'
```

#### Get All Participating Stores
```bash
curl -X 'GET' \
  'http://localhost:8080/api/stores/participating' \
  -H 'accept: application/json'
```

#### Get All Stores
```bash
curl -X 'GET' \
  'http://localhost:8080/api/stores' \
  -H 'accept: application/json'
```

#### Rebrand Store
```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/{store-id}/rebrand' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "New Store Name"
}'
```

#### Disable Store Participation
```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/{store-id}/disable-participation' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

### Product Management

#### Get Product Details
```bash
curl -X 'GET' \
  'http://localhost:8080/api/stores/products/{product-id}' \
  -H 'accept: application/json'
```

#### Rebrand Product
```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/products/{product-id}/rebrand' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "New Product Name",
  "description": "Updated description"
}'
```

#### Increase Product Price
```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/products/{product-id}/increase-price' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 1099.99
}'
```

#### Decrease Product Price
```bash
curl -X 'POST' \
  'http://localhost:8080/api/stores/products/{product-id}/decrease-price' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "price": 899.99
}'
```

#### Remove Product
```bash
curl -X 'DELETE' \
  'http://localhost:8080/api/stores/products/{product-id}' \
  -H 'accept: application/json'
```

### Basket Management

#### Cancel Basket
```bash
curl -X 'POST' \
  'http://localhost:8080/api/baskets/{basket-id}/cancel' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

#### Remove Item from Basket
```bash
curl -X 'DELETE' \
  'http://localhost:8080/api/baskets/{basket-id}/items' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "productId": "{product-id}",
  "quantity": 1
}'
```

### Order Management

#### Get Order Details
```bash
curl -X 'GET' \
  'http://localhost:8080/api/ordering/orders/{order-id}' \
  -H 'accept: application/json'
```

#### Cancel Order
```bash
curl -X 'POST' \
  'http://localhost:8080/api/ordering/orders/{order-id}/cancel' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

#### Ready Order
```bash
curl -X 'POST' \
  'http://localhost:8080/api/ordering/orders/{order-id}/ready' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

### Payment Management

#### Adjust Invoice
```bash
curl -X 'POST' \
  'http://localhost:8080/api/payments/invoices/{invoice-id}/adjust' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "amount": 899.99
}'
```

#### Cancel Invoice
```bash
curl -X 'POST' \
  'http://localhost:8080/api/payments/invoices/{invoice-id}/cancel' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

## 🧪 Error Handling Tests

### Invalid Basket Operations

#### Try to Add Item to Non-existent Basket
```bash
curl -X 'POST' \
  'http://localhost:8080/api/baskets/non-existent-id/items' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "productId": "{product-id}",
  "quantity": 1
}'
```

#### Try to Checkout Empty Basket
```bash
curl -X 'POST' \
  'http://localhost:8080/api/baskets/{basket-id}/checkout' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "paymentId": "{payment-id}"
}'
```

### Invalid Order Operations

#### Try to Complete Non-existent Order
```bash
curl -X 'POST' \
  'http://localhost:8080/api/ordering/orders/non-existent-id/complete' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "invoiceId": "{invoice-id}"
}'
```

### Invalid Payment Operations

#### Try to Confirm Non-existent Payment
```bash
curl -X 'POST' \
  'http://localhost:8080/api/payments/non-existent-id/confirm' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json'
```

## 📊 Monitoring and Debugging

### View Application Logs
```bash
make logs
```

### Access Swagger Documentation

Each service provides its own Swagger UI documentation:

- **Stores API**: http://localhost:8080/stores-spec/
- **Baskets API**: http://localhost:8080/baskets-spec/
- **Customers API**: http://localhost:8080/customers-spec/
- **Ordering API**: http://localhost:8080/ordering-spec/
- **Payments API**: http://localhost:8080/payments-spec/

### Database Management

```bash
# Connect to database
docker exec -it postgres psql -U mallbots_user -d mallbots

# Useful PostgreSQL commands:
# \dt - list tables
# \d table_name - describe table
# \d+ table_name - describe table with details
# \q - quit
# \h - help
# \? - help
# \l - list databases
```

## 🔄 Event-Driven Architecture Flow

This application demonstrates a complete event-driven architecture:

1. **Commands** trigger domain events
2. **Domain events** trigger integration events
3. **Integration events** coordinate between bounded contexts
4. **Saga orchestrator** manages distributed transactions
5. **Event sourcing** maintains audit trail
6. **CQRS** separates read and write operations

### Flow Explanation

The e-commerce flow shows how a simple transaction involves multiple services:

1. **Customer Creation** → Customer domain events
2. **Store Setup** → Store domain events
3. **Product Addition** → Product domain events
4. **Basket Operations** → Basket domain events
5. **Payment Authorization** → Payment domain events
6. **Order Creation** → Order domain events + Integration events
7. **Payment Confirmation** → Payment domain events
8. **Invoice Creation** → Invoice domain events
9. **Order Completion** → Order domain events + Integration events

Each step triggers events that other services can react to, demonstrating the power of event-driven architecture in handling complex business processes.

## 🧹 Cleanup

### Stop All Services
```bash
make down
```

### Remove All Data (Volumes)
```bash
make down -v
```

## 📝 Notes

- Replace `{customer-id}`, `{store-id}`, `{product-id}`, `{basket-id}`, `{payment-id}`, `{order-id}`, and `{invoice-id}` with actual UUIDs returned from previous API calls
- All timestamps are in ISO 8601 format
- The application uses event sourcing, so all operations are auditable
- The saga orchestrator manages distributed transactions across services
- Integration events coordinate between different bounded contexts

## 🔧 Troubleshooting

### Common Issues

1. **Services not starting**: Check if Docker and Docker Compose are running
2. **Database connection errors**: Ensure PostgreSQL container is healthy
3. **NATS connection errors**: Check if NATS container is running
4. **API not responding**: Verify the application is running on port 8080

### Debug Commands

```bash
# Check container status
docker-compose ps

# Check application logs
docker-compose logs monolith

# Check database logs
docker-compose logs postgres

# Check NATS logs
docker-compose logs nats

# Restart services
docker-compose restart
```
