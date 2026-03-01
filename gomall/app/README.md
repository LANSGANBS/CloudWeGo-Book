# Gomall 应用服务

本目录包含 Gomall 电商平台的所有微服务应用。

## 📁 服务目录

### HTTP 服务
- **frontend/** - Web 前端服务 (端口: 8080)
  - 基于 Hertz 框架
  - 提供用户界面和 API 网关功能

### RPC 服务 (基于 Kitex)

| 服务 | 目录 | 端口 | 功能描述 |
|------|------|------|----------|
| User | user/ | 8881 | 用户注册、登录、认证 |
| Product | product/ | 8882 | 商品管理、分类、搜索 |
| Cart | cart/ | 8883 | 购物车增删改查 |
| Payment | payment/ | 8884 | 支付处理 |
| Checkout | checkout/ | 8885 | 结账流程编排 |
| Order | order/ | 8886 | 订单管理 |
| Email | email/ | 8887 | 邮件通知 (NATS) |

## 🚀 快速启动

### 启动顺序
服务之间有依赖关系,建议按以下顺序启动:

```bash
# 1. 基础服务
cd user && go run .
cd product && go run .

# 2. 业务服务
cd cart && go run .
cd payment && go run .

# 3. 编排服务
cd checkout && go run .
cd order && go run .
cd email && go run .

# 4. 前端服务
cd frontend && go run .
```

### 单个服务启动
```bash
cd <service-name>
go run .
```

## 🔧 开发

### 环境配置
每个服务需要 `.env` 文件,参考各服务目录下的 `.env.example`

### 依赖管理
```bash
# 修复所有服务依赖
cd .. && .\fix-deps.ps1

# 单个服务
cd <service-name>
go mod tidy
```

### 构建
```bash
# 构建单个服务
cd <service-name>
go build -o bin/<service-name> .
```

## 📝 服务详情

详见各服务目录下的 README.md 文件。
