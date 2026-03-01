#!/bin/bash

# Gomall 服务停止脚本 - macOS 版本

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

echo "=== 停止 Gomall 服务 ==="
echo ""

# 停止 Docker 容器
echo "1. 停止基础设施服务..."
docker-compose down

echo ""
echo "2. 停止所有 Go 进程..."
pkill -f "go run" || true

echo ""
echo "=== 所有服务已停止 ==="
