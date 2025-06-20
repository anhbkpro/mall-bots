# MallBots - Event-Driven Architecture in Go

A microservices-based e-commerce platform built with Go, gRPC, and event-driven architecture.

## 🏗️ Architecture

This project implements a **monolithic modular architecture** with the following services:

- **🛒 Baskets** - Shopping cart management
- **👥 Customers** - Customer management
- **📦 Depot** - Inventory and warehouse operations
- **🔔 Notifications** - User notifications
- **📋 Ordering** - Order processing
- **💳 Payments** - Payment processing
- **🏪 Stores** - Store management

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make (optional, for convenience)

### 1. Install Development Tools

```bash
make install-tools
```

This installs:

- Protocol Buffer compiler
- gRPC tools
- Swagger tools

### 2. Generate Code

```bash
make generate
# or
cd baskets # or any other service
buf generate
```

This generates:

- Protocol Buffer Go files
- gRPC service definitions
- Swagger documentation

### 3. Start the Application

#### Option A: Using Docker (Recommended)

```bash
# Start all services
make up

# View logs
make logs

# Stop services
make down
```

#### Option B: Local Development

```bash
# Build the application
make build

# Run locally (requires PostgreSQL)
make dev
```

### 4. Access the Application

- **Web UI**: http://localhost:8080
- **gRPC Server**: localhost:8085
- **PostgreSQL**: localhost:5432

## 🛠️ Development

### Project Structure

```text
├── cmd/mallbots/          # Application entry point
├── internal/              # Shared internal packages
│   ├── config/           # Configuration management
│   ├── ddd/              # Domain-driven design utilities
│   ├── logger/           # Logging utilities
│   ├── monolith/         # Monolith interface
│   ├── rpc/              # gRPC configuration
│   ├── waiter/           # Graceful shutdown utilities
│   └── web/              # Web server utilities
├── baskets/              # Baskets service
├── customers/            # Customers service
├── depot/                # Depot service
├── notifications/        # Notifications service
├── ordering/             # Ordering service
├── payments/             # Payments service
├── stores/               # Stores service
├── docker/               # Docker configuration
└── proto/                # Protocol Buffer definitions
```

### Available Make Commands

```bash
# Install development tools
make install-tools

# Generate protocol buffer and swagger files
make generate

# Build the application
make build

# Run tests
make test

# Clean build artifacts
make clean

# Start development environment
make dev

# Stop development environment
make dev-stop

# Start all services
make up

# Stop services
make down

# View logs
make logs
```

### Database Management

```bash
# Reset database (removes all data)
make down -v
make up

# View database logs
docker-compose logs postgres

# Connect to database
docker exec -it postgres psql -U mallbots_user -d mallbots
# password: mallbots_pass
# Once connected, useful PostgreSQL commands:
# \dt - list tables
# \d table_name - describe table
# \d+ table_name - describe table with details
# \q - quit
# \h - help
# \? - help
# \l - list databases
```

### API Documentation

Each service provides its own Swagger UI documentation:

- **Baskets**: http://localhost:8080/baskets-spec/
- **Customers**: http://localhost:8080/customers-spec/
- **Depot**: http://localhost:8080/depot-spec/
- **Ordering**: http://localhost:8080/ordering-spec/
- **Payments**: http://localhost:8080/payments-spec/
- **Stores**: http://localhost:8080/stores-spec/

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
ENVIRONMENT=development
LOG_LEVEL=DEBUG
PG_CONN=postgres://mallbots_user:mallbots_pass@localhost:5432/mallbots?sslmode=disable
WEB_HOST=0.0.0.0
WEB_PORT=8080
RPC_HOST=0.0.0.0
RPC_PORT=8085
SHUTDOWN_TIMEOUT=30s
```

### Docker Configuration

The application uses Docker Compose with:

- **PostgreSQL 12** - Database
- **Custom Go application** - Monolith service

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific service tests
go test ./baskets/...
go test ./customers/...
```

## 📊 Monitoring

### Logs

```bash
# View application logs
docker-compose logs -f monolith

# View database logs
docker-compose logs -f postgres

# View all logs
docker-compose logs -f
```

### Health Checks

- **Web Server**: http://localhost:8080
- **gRPC Server**: Use gRPC reflection on port 8085

## 🚀 Deployment

### Production Build

```bash
# Build production image
docker build -f docker/Dockerfile -t mallbots:latest .

# Run production container
docker run -p 8080:8080 -p 8085:8085 mallbots:latest
```

### Environment-Specific Configurations

- **Development**: Uses local PostgreSQL
- **Production**: Configure external database and secrets

## 🐛 Troubleshooting

### Common Issues

1. **Port 8080 already in use**
   ```bash
   # Find process using port
   lsof -i :8080

   # Kill process or change port in docker-compose.yml
   ```

2. **Database connection issues**
   ```bash
   # Restart database
   docker-compose restart postgres

   # Check database logs
   docker-compose logs postgres
   ```

3. **Protocol buffer generation errors**
   ```bash
   # Reinstall tools
   make install-tools

   # Regenerate
   make generate
   ```

### Debug Mode

```bash
# Run with debug logging
LOG_LEVEL=DEBUG docker-compose up

# Run locally with debug
LOG_LEVEL=DEBUG ./monolith
```

## 📝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make generate` if you modified protobuf files
6. Submit a pull request
