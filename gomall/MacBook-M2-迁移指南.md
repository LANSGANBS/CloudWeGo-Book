# Gomall MacBook M2 迁移指南

## 概述

本文档指导你将 Gomall 项目从 Windows + Ubuntu VM 环境迁移到 MacBook M2 本地开发环境。

---

## 第一步：安装必要软件

### 1.1 安装 Homebrew（如果还没有）

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 1.2 安装 Go 1.23+

```bash
brew install go
```

验证安装：

```bash
go version
```

### 1.3 安装 Docker Desktop for Mac

1. 下载地址：https://www.docker.com/products/docker-desktop
2. 选择 **Apple Silicon (M1/M2/M3)** 版本
3. 安装后启动 Docker Desktop
4. 等待 Docker 图标显示为运行状态

验证安装：

```bash
docker --version
docker-compose --version
```

### 1.4 安装可选工具

```bash
# 代码生成工具
brew install cwgo

# 热重载工具
brew install air

# Kitex 工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

---

## 第二步：克隆项目

```bash
# 如果是从 GitHub 克隆
git clone https://github.com/cloudwego/biz-demo.git
cd biz-demo/gomall

# 如果是从 Windows 迁移，直接复制项目文件夹即可
```

---

## 第三步：启动基础设施服务

### 3.1 方式一：使用已有独立容器（推荐）

如果您已有独立运行的 Docker 容器，只需确保以下服务运行中：

```bash
# 检查已有容器
docker ps

# 启动缺失的服务（如 Jaeger）
docker run -d --name jaeger-all-in-one -p 16686:16686 -p 4317:4317 -p 4318:4318 jaegertracing/all-in-one:latest

# 启动 RocketMQ（独立管理）
cd rocketmq && docker-compose up -d
```

### 3.2 方式二：使用 gomall 的 docker-compose

```bash
cd gomall
docker-compose up -d

# 单独启动 RocketMQ（独立管理）
cd ../rocketmq && docker-compose up -d
```

### 3.3 验证容器状态

```bash
docker ps
```

必需服务：
| 容器 | 端口 | 用途 |
|------|------|------|
| mysql | 3306 | 数据库 |
| redis | 6379 | 缓存 |
| consul | 8500 | 服务发现 |
| jaeger-all-in-one | 16686, 4317, 4318 | 链路追踪 |
| rocketmq-namesrv | 9876 | RocketMQ NameServer |
| rocketmq-broker | 10911, 10909 | RocketMQ Broker |

可选服务：
| 容器 | 端口 | 用途 |
|------|------|------|
| prometheus | 9090 | 监控指标 |
| grafana | 3000 | 监控面板 |
| etcd | 2379, 2380 | 键值存储 |
| nats | 4222, 8222 | 消息队列 |

### 3.4 创建 RocketMQ Topics

```bash
# 等待 RocketMQ 启动后创建 topics
docker exec -it rocketmq-broker sh -c "cd /home/rocketmq/rocketmq-5.1.0/bin && ./mqadmin updateTopic -n rocketmq-namesrv:9876 -t stock_deduct -c DefaultCluster"

docker exec -it rocketmq-broker sh -c "cd /home/rocketmq/rocketmq-5.1.0/bin && ./mqadmin updateTopic -n rocketmq-namesrv:9876 -t stock_restore -c DefaultCluster"

docker exec -it rocketmq-broker sh -c "cd /home/rocketmq/rocketmq-5.1.0/bin && ./mqadmin updateTopic -n rocketmq-namesrv:9876 -t stock_sync -c DefaultCluster"
```

### 3.4 等待 MySQL 初始化

首次启动时，MySQL 需要初始化数据库。等待约 30-60 秒：

```bash
# 检查 MySQL 日志
docker logs gomall-mysql-1

# 测试连接
docker exec -it gomall-mysql-1 mysql -uroot -p123456 -e "SELECT 1"
```

### 3.4 导入完整数据库（重要！）

默认只创建空数据库，需要导入完整数据：

```bash
# 方式一：使用合并的初始化脚本
docker exec -i gomall-mysql-1 mysql -uroot -p123456 < sql/init_all.sql

# 方式二：逐个导入
docker exec -i gomall-mysql-1 mysql -uroot -p123456 < sql/user.sql
docker exec -i gomall-mysql-1 mysql -uroot -p123456 < sql/product.sql
docker exec -i gomall-mysql-1 mysql -uroot -p123456 < sql/cart.sql
docker exec -i gomall-mysql-1 mysql -uroot -p123456 < sql/order.sql
docker exec -i gomall-mysql-1 mysql -uroot -p123456 < sql/payment.sql
docker exec -i gomall-mysql-1 mysql -uroot -p123456 < sql/checkout.sql
```

验证数据库：

```bash
docker exec -it gomall-mysql-1 mysql -uroot -p123456 -e "SHOW DATABASES; USE user; SELECT * FROM user;"
```

---

## 第四步：配置环境变量

项目已经配置好 `.env` 文件，使用 `127.0.0.1` 作为服务地址。

### 验证 .env 文件

```bash
# 检查 user 服务的 .env
cat app/user/.env
```

应该看到：

```
MYSQL_USER=root
MYSQL_PASSWORD=123456
MYSQL_HOST=127.0.0.1
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://127.0.0.1:4317
OTEL_EXPORTER_OTLP_INSECURE=true
```

---

## 第五步：下载依赖

```bash
# 先进入 gomall 目录
cd gomall

# 在 gomall 目录下执行
make tidy
```

## 第六步：启动应用服务

### 6.1 停止旧服务（如果有端口冲突）

```bash
# 查看占用端口的进程
lsof -i :8080 -i :8881 -i :8882 -i :8883 -i :8884 -i :8885 -i :8886 -i :8887

# 停止所有 gomall 服务进程
pkill -f "user|product|cart|payment|checkout|order|email|frontend" || true
```

### 6.2 启动服务

#### 方式一：使用 Make（推荐）

```bash
# 在 gomall 目录下，不同终端窗口分别启动
make run svc=user       # 终端 1
make run svc=product    # 终端 2
make run svc=cart       # 终端 3
make run svc=payment    # 终端 4
make run svc=checkout   # 终端 5
make run svc=order      # 终端 6
make run svc=email      # 终端 7
make run svc=frontend   # 终端 8
```

#### 方式二：手动启动

```bash
# 终端 1 - User 服务
cd gomall/app/user
go run .

# 终端 2 - Product 服务
cd gomall/app/product
go run .

# 终端 3 - Cart 服务
cd gomall/app/cart
go run .

# 终端 4 - Payment 服务
cd gomall/app/payment
go run .

# 终端 5 - Checkout 服务
cd gomall/app/checkout
go run .

# 终端 6 - Order 服务
cd gomall/app/order
go run .

# 终端 7 - Email 服务
cd gomall/app/email
go run .

# 终端 8 - Frontend 服务
cd gomall/app/frontend
go run .
```

### 启动顺序

按以下顺序启动服务（有依赖关系）：

1. **user** (8881) → 2. **product** (8882) → 3. **cart** (8883) → 4. **payment** (8884) → 5. **checkout** (8885) → 6. **order** (8886) → 7. **email** (8887) → 8. **frontend** (8080)

每个服务启动后等待看到类似日志：

```
[Info] KITEX: server listen at addr=[::]:8881
```

---

## 第七步：验证部署

### 7.1 检查服务状态

```bash
# 使用检查脚本
chmod +x check-status.sh
./check-status.sh

# 或手动检查端口
lsof -i :8881  # user
lsof -i :8882  # product
lsof -i :8080  # frontend
```

### 7.2 访问应用

| 服务           | 地址                                |
| -------------- | ----------------------------------- |
| **前端网站**   | http://localhost:8080               |
| **Consul UI**  | http://localhost:8500               |
| **Jaeger UI**  | http://localhost:16686              |
| **Grafana**    | http://localhost:3000 (admin/admin) |
| **Prometheus** | http://localhost:9090               |

---

## 常见问题

### Q1: Docker Desktop 启动慢

**解决方法**：

- 确保 Docker Desktop 已更新到最新版本
- 在 Docker Desktop 设置中增加内存限制（建议 8GB+）

### Q2: MySQL 连接失败

**错误信息**：`Access denied for user 'root'@'127.0.0.1'`

**解决方法**：

```bash
# 重启 MySQL 容器
docker-compose restart mysql

# 检查 MySQL 日志
docker logs gomall-mysql-1
```

### Q3: 端口被占用

**错误信息**：`bind: address already in use`

**解决方法**：

```bash
# 查找占用端口的进程
lsof -i :8881

# 终止进程
kill -9 <PID>
```

### Q4: Go 模块依赖错误

**解决方法**：

```bash
# 清理 Go 缓存
go clean -modcache

# 重新下载依赖
make tidy
```

### Q5: M2 芯片兼容性问题

某些镜像可能不支持 ARM 架构。如果遇到问题：

```bash
# 使用 --platform linux/amd64 强制使用 x86 镜像
docker pull --platform linux/amd64 <image-name>
```

### Q6: RocketMQ 连接失败

**解决方法**：

```bash
# 检查 RocketMQ 容器状态
docker ps | grep rocketmq

# 查看 RocketMQ 日志
docker logs rocketmq-namesrv
docker logs rocketmq-broker

# 重启 RocketMQ
docker-compose restart rocketmq-namesrv rocketmq-broker
```

### Q7: RocketMQ Topic 不存在（Mac 迁移后常见）

**问题现象**：
```
remote or network error: rpc error: code = 13 desc = DeductStock.err
```

**原因**：Mac 迁移后 RocketMQ 需要重新创建 topics

**解决方法**：

```bash
# 创建必要的 topics
docker exec -it rocketmq-broker sh -c "cd /home/rocketmq/rocketmq-5.1.0/bin && ./mqadmin updateTopic -n rocketmq-namesrv:9876 -t stock_deduct -c DefaultCluster"

docker exec -it rocketmq-broker sh -c "cd /home/rocketmq/rocketmq-5.1.0/bin && ./mqadmin updateTopic -n rocketmq-namesrv:9876 -t stock_restore -c DefaultCluster"

docker exec -it rocketmq-broker sh -c "cd /home/rocketmq/rocketmq-5.1.0/bin && ./mqadmin updateTopic -n rocketmq-namesrv:9876 -t stock_sync -c DefaultCluster"

# 验证 topics 是否创建成功
docker exec -it rocketmq-namesrv sh -c "cd /home/rocketmq/rocketmq-5.1.0/bin && ./mqadmin topicList -n localhost:9876"
```

---

## 最小启动配置（内存不足时）

如果内存不足，可以只启动核心服务：

1. **User 服务** (8881) - 必需
2. **Product 服务** (8882) - 必需
3. **Frontend 服务** (8080) - 必需

这样可以访问基本的商品浏览功能。

---

## 停止服务

### 停止应用服务

在每个服务终端按 `Ctrl+C`

### 停止 Docker 容器

```bash
docker-compose down
```

### 清理数据（可选）

```bash
# 删除所有容器和卷
docker-compose down -v
```

---

## 与 Windows 环境的主要区别

| 项目            | Windows + VM   | MacBook M2            |
| --------------- | -------------- | --------------------- |
| Docker 运行位置 | Ubuntu 虚拟机  | 本地 Docker Desktop   |
| 服务地址        | 192.168.63.131 | 127.0.0.1 / localhost |
| 端口映射        | VM → Windows   | 容器 → 主机           |
| 性能            | 网络延迟       | 本地直连，更快        |
| 架构            | x86_64         | ARM64 (Apple Silicon) |

---

## 端口映射汇总

| 服务             | 容器端口     | 主机端口     | 用途                |
| ---------------- | ------------ | ------------ | ------------------- |
| MySQL            | 3306         | 3306         | 数据库              |
| Redis            | 6379         | 6379         | 缓存                |
| Consul           | 8500         | 8500         | 服务发现            |
| Jaeger UI        | 16686        | 16686        | 链路追踪界面        |
| Jaeger OTLP gRPC | 4317         | 4317         | 链路数据接收        |
| Jaeger OTLP HTTP | 4318         | 4318         | 链路数据接收        |
| Prometheus       | 9090         | 9090         | 监控指标            |
| Grafana          | 3000         | 3000         | 监控面板            |
| etcd             | 2379, 2380   | 2379, 2380   | 键值存储            |
| NATS             | 4222, 8222   | 4222, 8222   | 消息队列            |
| Loki             | 3100         | 3100         | 日志聚合            |
| RocketMQ NS      | 9876         | 9876         | RocketMQ NameServer |
| RocketMQ Broker  | 10911, 10909 | 10911, 10909 | RocketMQ Broker     |

---

## 下一步

1. 访问 http://localhost:8080 测试前端
2. 在 Consul UI 查看服务注册情况
3. 在 Jaeger UI 查看链路追踪
4. 在 Grafana 配置监控面板

祝你在 MacBook M2 上开发愉快！
