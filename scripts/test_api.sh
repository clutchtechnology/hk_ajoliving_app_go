#!/bin/bash

# AJO Living API - API 测试脚本
# 用于快速测试 Docker 容器中的 API 是否正常工作

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

API_URL="${API_URL:-http://localhost:8080}"

info "=========================================="
info "  AJO Living API - API 测试"
info "=========================================="
echo ""
info "测试目标: $API_URL"
echo ""

# 测试健康检查
test_health() {
    info "测试 1: 健康检查..."
    response=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/health" || echo "000")
    
    if [ "$response" = "200" ]; then
        echo "✅ 健康检查通过 (HTTP $response)"
    else
        error "健康检查失败 (HTTP $response)"
        return 1
    fi
}

# 测试 API 版本
test_version() {
    info "测试 2: API 版本..."
    response=$(curl -s "$API_URL/api/v1/health" || echo "")
    
    if [ ! -z "$response" ]; then
        echo "✅ API 响应成功"
        echo "   响应: $response"
    else
        warn "API 无响应"
    fi
}

# 测试数据库连接
test_database() {
    info "测试 3: 数据库连接..."
    info "检查 PostgreSQL 容器..."
    
    if docker-compose exec -T postgres pg_isready -U ajoliving &> /dev/null; then
        echo "✅ PostgreSQL 连接正常"
    else
        error "PostgreSQL 连接失败"
        return 1
    fi
}

# 测试 Redis 连接
test_redis() {
    info "测试 4: Redis 连接..."
    
    if docker-compose exec -T redis redis-cli ping &> /dev/null; then
        echo "✅ Redis 连接正常"
    else
        error "Redis 连接失败"
        return 1
    fi
}

# 测试用户注册
test_register() {
    info "测试 5: 用户注册..."
    
    timestamp=$(date +%s)
    response=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"testuser_$timestamp\",
            \"email\": \"test_$timestamp@example.com\",
            \"password\": \"Test123456\",
            \"phone\": \"85298765432\"
        }" || echo "")
    
    if [[ "$response" == *"code"* ]]; then
        echo "✅ 注册接口响应正常"
        echo "   响应: $response"
    else
        warn "注册接口可能未实现或有问题"
    fi
}

# 运行所有测试
run_tests() {
    test_health || exit 1
    test_version
    test_database || exit 1
    test_redis || exit 1
    test_register
    
    echo ""
    info "=========================================="
    info "  测试完成！"
    info "=========================================="
}

# 等待 API 就绪
wait_for_api() {
    info "等待 API 服务就绪..."
    for i in {1..30}; do
        if curl -s -f "$API_URL/health" &> /dev/null; then
            echo "✅ API 服务已就绪"
            return 0
        fi
        echo -n "."
        sleep 1
    done
    echo ""
    error "API 服务未能在 30 秒内就绪"
    exit 1
}

# 主函数
main() {
    wait_for_api
    echo ""
    run_tests
}

main "$@"
