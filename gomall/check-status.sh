#!/bin/bash

# Gomall Service Status Check Script for macOS

LOCALHOST="127.0.0.1"

echo "=== Gomall Service Status Check ==="
echo ""

# Check infrastructure services
echo "Infrastructure Services:"
infra_services=(
    "MySQL:3306"
    "Redis:6379"
    "Consul:8500"
)

for service in "${infra_services[@]}"; do
    name="${service%%:*}"
    port="${service##*:}"
    echo -n "  $name ($port)... "
    
    if nc -z "$LOCALHOST" "$port" 2>/dev/null; then
        echo -e "\033[32mOK\033[0m"
    else
        echo -e "\033[31mFAILED\033[0m"
    fi
done

echo ""
echo "Application Services:"

# Check application services
app_services=(
    "user:8881"
    "product:8882"
    "cart:8883"
    "payment:8884"
    "checkout:8885"
    "order:8886"
    "email:8887"
    "frontend:8080"
)

for service in "${app_services[@]}"; do
    name="${service%%:*}"
    port="${service##*:}"
    echo -n "  $name ($port)... "
    
    if nc -z "$LOCALHOST" "$port" 2>/dev/null; then
        echo -e "\033[32mOK\033[0m"
    else
        echo -e "\033[33mNOT READY\033[0m"
    fi
done

echo ""
echo "Access URLs:"
echo "  Frontend:  http://localhost:8080"
echo "  Consul UI: http://localhost:8500"
echo "  Jaeger UI: http://localhost:16686"
echo "  Grafana:   http://localhost:3000"
echo "  Prometheus: http://localhost:9090"
