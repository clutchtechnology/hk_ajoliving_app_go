#!/bin/bash

# AJO Living API - 快速启动脚本
# 一键构建、启动和初始化所有服务

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

info "=========================================="
info "  AJO Living API - 快速启动"
info "=========================================="
echo ""

# 1. 检查环境
info "步骤 1/5: 检查 Docker 环境..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装"
    exit 1
fi
echo "✅ Docker 已安装"

# 2. 创建环境变量文件
info "步骤 2/5: 准备环境变量..."
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ 已创建 .env 文件"
    warn "请在生产环境中修改 JWT_SECRET 等敏感信息！"
else
    echo "✅ .env 文件已存在"
fi

# 3. 构建镜像
info "步骤 3/5: 构建 Docker 镜像..."
docker-compose build
echo "✅ 镜像构建完成"

# 4. 启动服务
info "步骤 4/5: 启动服务..."
docker-compose up -d
echo "✅ 服务已启动"

# 5. 等待服务就绪
info "步骤 5/5: 等待服务就绪..."
echo -n "等待 PostgreSQL 就绪"
for i in {1..30}; do
    if docker-compose exec -T postgres pg_isready -U ajoliving &> /dev/null; then
        echo ""
        echo "✅ PostgreSQL 已就绪"
        break
    fi
    echo -n "."
    sleep 1
done

sleep 2

# 6. 显示服务状态
info "=========================================="
info "  服务启动成功！"
info "=========================================="
echo ""
docker-compose ps
echo ""
info "🚀 API 服务地址: http://localhost:8080"
info "📊 PostgreSQL: localhost:5432"
info "💾 Redis: localhost:6379"
echo ""
info "查看日志: docker-compose logs -f"
info "停止服务: docker-compose down"
info "或使用: ./scripts/deploy_docker.sh [命令]"
echo ""
