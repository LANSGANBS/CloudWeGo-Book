#!/bin/bash

# Gomall 服务启动脚本 - macOS 版本

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

echo "=== Gomall 服务启动脚本 ==="
echo ""

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "错误: Docker 未运行，请先启动 Docker Desktop"
    exit 1
fi

# 启动基础设施服务
echo "1. 启动基础设施服务..."
docker-compose up -d

echo ""
echo "2. 等待 MySQL 初始化..."
sleep 30

# 测试 MySQL 连接
echo "3. 测试 MySQL 连接..."
if docker exec gomall-mysql-1 mysql -uroot -p123456 -e "SELECT 1" > /dev/null 2>&1; then
    echo "   MySQL 连接成功"
else
    echo "   警告: MySQL 连接失败，请检查日志"
    docker logs gomall-mysql-1
fi

echo ""
echo "4. 下载 Go 依赖..."
make tidy

echo ""
echo "=== 基础设施服务已启动 ==="
echo ""
echo "请在新终端中启动应用服务："
echo ""
echo "  cd app/user && go run ."
echo "  cd app/product && go run ."
echo "  cd app/cart && go run ."
echo "  cd app/payment && go run ."
echo "  cd app/checkout && go run ."
echo "  cd app/order && go run ."
echo "  cd app/email && go run ."
echo "  cd app/frontend && go run ."
echo ""
echo "或者使用 make 命令："
echo ""
echo "  make run svc=user"
echo "  make run svc=product"
echo "  make run svc=cart"
echo "  make run svc=payment"
echo "  make run svc=checkout"
echo "  make run svc=order"
echo "  make run svc=email"
echo "  make run svc=frontend"
echo ""
echo "启动顺序: user → product → cart → payment → checkout → order → email → frontend"
echo ""
echo "访问地址:"
echo "  前端: http://localhost:8080"
echo "  Consul: http://localhost:8500"
echo "  Jaeger: http://localhost:16686"
echo "  Grafana: http://localhost:3000"
