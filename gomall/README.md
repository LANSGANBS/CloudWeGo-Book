# Gomall - Microservices E-commerce Platform

[中文](README_cn.md) | [English](#english)

> A modern microservices-based e-commerce platform built with CloudWeGo ecosystem, demonstrating best practices for building high-performance distributed systems.

## 📖 Overview

Gomall is a teaching project that showcases how to build a production-ready microservices architecture using CloudWeGo's Kitex (RPC) and Hertz (HTTP) frameworks. It implements a complete e-commerce system with 8 independent services.

## 🏗️ Architecture

### Services Overview

| Service | Port | Type | Description |
|---------|------|------|-------------|
| **frontend** | 8080 | HTTP | Web frontend service (Hertz) |
| **user** | 8881 | RPC | User authentication & management (Kitex) |
| **product** | 8882 | RPC | Product catalog & inventory (Kitex) |
| **cart** | 8883 | RPC | Shopping cart management (Kitex) |
| **payment** | 8884 | RPC | Payment processing (Kitex) |
| **checkout** | 8885 | RPC | Checkout orchestration (Kitex) |
| **order** | 8886 | RPC | Order management (Kitex) |
| **email** | 8887 | RPC | Email notifications (Kitex) |

### Infrastructure Components

| Component | Port | Purpose |
|-----------|------|---------|
| MySQL | 3306 | Primary database |
| Redis | 6379 | Cache & session storage |
| Consul | 8500 | Service discovery & registry |
| Jaeger | 16686 | Distributed tracing |
| Prometheus | 9090 | Metrics collection |
| Grafana | 3000 | Monitoring dashboards |
| NATS | 4222 | Message queue |
| Loki | 3100 | Log aggregation |

## 🛠️ Technology Stack

| Technology | Purpose | Documentation |
|------------|---------|---------------|
| [Kitex](https://github.com/cloudwego/kitex) | High-performance RPC framework | [Docs](https://www.cloudwego.io/docs/kitex/) |
| [Hertz](https://github.com/cloudwego/hertz) | High-performance HTTP framework | [Docs](https://www.cloudwego.io/docs/hertz/) |
| [cwgo](https://github.com/cloudwego/cwgo) | CloudWeGo code generation tool | [Docs](https://www.cloudwego.io/docs/cwgo/) |
| [Bootstrap](https://getbootstrap.com/) | Frontend UI toolkit | [Docs](https://getbootstrap.com/docs/) |
| MySQL | Relational database | - |
| Redis | In-memory data store | - |
| Consul | Service mesh & discovery | - |
| Prometheus | Monitoring & alerting | - |
| Jaeger | Distributed tracing | - |
| NATS | Message streaming | - |
| Docker | Containerization | - |


## ✨ Features

### Implemented Business Logic
- [x] User authentication & authorization
- [x] User registration & login/logout
- [x] Product catalog with categories
- [x] Product search & filtering
- [x] Shopping cart management
- [x] Real-time cart badge updates
- [x] Checkout process
- [x] Payment processing
- [x] Order management & history
- [x] Email notifications (async via NATS)
- [x] Session management
- [x] Distributed tracing
- [x] Metrics & monitoring

## 🚀 Quick Start

### Prerequisites

**Required:**
- Go 1.23+ ([Download](https://go.dev/dl/))
- Docker & Docker Compose ([Download](https://www.docker.com/))
- IDE / Code Editor (VS Code, GoLand, etc.)

**Optional (for development):**
- [cwgo](https://github.com/cloudwego/cwgo) - Code generation tool
- [kitex](https://github.com/cloudwego/kitex) - `go install github.com/cloudwego/kitex/tool/cmd/kitex@latest`
- [Air](https://github.com/cosmtrek/air) - Hot reload for Go apps

### Installation Steps

#### 1. Clone Repository
```bash
git clone https://github.com/cloudwego/biz-demo.git
cd biz-demo/gomall
```

#### 2. Fix Go Module Dependencies
```bash
# Run the dependency fix script
.\fix-deps.ps1

# Or manually for each service
cd app/user && go mod tidy
cd app/product && go mod tidy
# ... repeat for all services
```

#### 3. Start Infrastructure Services
```bash
# Start MySQL, Redis, Consul, Jaeger, etc.
docker-compose up -d

# Verify all containers are running
docker ps

# Check logs if needed
docker-compose logs -f
```

#### 4. Configure Environment Variables
```bash
# Initialize .env files for all services
make init

# Or manually create .env files (see .env.example)
```

**Important:** Generate a random SESSION_SECRET value for the frontend service.

#### 5. Download Dependencies
```bash
# Download all Go modules
make tidy

# Or manually
cd app/user && go mod download
# ... repeat for all services
```

#### 6. Start Application Services

**Option A: Using Make (Recommended)**
```bash
# Start a specific service with hot reload
make run svc=user
make run svc=product
make run svc=frontend
```

**Option B: Manual Start (Production)**
```bash
# Start services in dependency order
cd app/user && go run .      # Port 8881
cd app/product && go run .   # Port 8882
cd app/cart && go run .      # Port 8883
cd app/payment && go run .   # Port 8884
cd app/checkout && go run .  # Port 8885
cd app/order && go run .     # Port 8886
cd app/email && go run .     # Port 8887
cd app/frontend && go run .  # Port 8080
```

**Option C: Using PowerShell Script (Windows)**
```powershell
# Check service status
.\check-status.ps1

# See manual startup guide for detailed instructions
# 手动启动指南-完整版.md
```

#### 7. Verify Deployment
```bash
# Check service status
.\check-status.ps1

# Or manually check each port
curl http://localhost:8080  # Frontend
curl http://localhost:8881  # User service
# ... etc
```

### Access Application

- **Frontend Website**: http://localhost:8080
- **Consul UI**: http://localhost:8500
- **Jaeger Tracing**: http://localhost:16686
- **Grafana Dashboards**: http://localhost:3000 (admin/admin)
- **Prometheus Metrics**: http://localhost:9090

### Quick Commands

```bash
# View all available make commands
make

# Open Gomall website
make open-gomall

# Open Consul registry
make open-consul

# Stop infrastructure services
make env-stop
```
## 🐛 Troubleshooting

### Common Issues

**MySQL Connection Failed (Error 1045)**
```bash
# Check MySQL container
docker ps | grep mysql

# Test connection
docker exec -it <mysql-container> mysql -uroot -p123456

# Check .env file configuration
cat app/user/.env
```

**Service Exits Immediately**
- Check if port is already in use: `netstat -ano | findstr "8881"`
- Verify configuration files are correct
- Check service logs for error messages

**Frontend Returns 502 Error**
- Ensure all backend RPC services are running
- Check Consul UI for service registration status
- Wait 1-2 minutes for services to fully initialize
- Verify network connectivity between services

**Go Module Dependency Errors**
```bash
# Run the fix script
.\fix-deps.ps1

# Or manually fix
cd app/user && go mod tidy
```

## 💻 Development

### Project Structure
```
gomall/
├── app/                    # Application services
│   ├── frontend/          # HTTP frontend (Hertz)
│   ├── user/              # User RPC service
│   ├── product/           # Product RPC service
│   ├── cart/              # Cart RPC service
│   ├── payment/           # Payment RPC service
│   ├── checkout/          # Checkout RPC service
│   ├── order/             # Order RPC service
│   └── email/             # Email RPC service
├── common/                # Shared utilities
├── rpc_gen/              # Generated RPC code
├── docker-compose.yaml   # Infrastructure services
└── README.md
```

### Building Services
```bash
# Build all services
make build

# Build specific service
cd app/user && go build -o bin/user .
```

### Running Tests
```bash
# Run all tests
make test

# Run tests for specific service
cd app/user && go test ./...

# Run tests with coverage
go test -cover ./...
```

### Code Generation
```bash
# Generate RPC code using cwgo
cwgo server --type RPC --idl <idl-file> --service <service-name>

# Generate client code
cwgo client --type RPC --idl <idl-file> --service <service-name>
```

## 📚 Documentation

- [手动启动指南](./手动启动指南-完整版.md) - Detailed manual startup guide (Chinese)
- [Architecture Design](./docs/architecture.md) - System architecture details
- [API Documentation](./docs/api.md) - API reference
- [Development Guide](./docs/development.md) - Development setup

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 👥 Contributors
- [rogerogers](https://github.com/rogerogers)
- [baiyutang](https://github.com/baiyutang)
