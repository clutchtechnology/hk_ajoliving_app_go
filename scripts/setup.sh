#!/bin/bash

# AJO Living API - 项目设置脚本

set -e

echo "🚀 开始设置 AJO Living API 项目..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 未检测到 Go。请先安装 Go 1.21 或更高版本。"
    echo "访问: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go 版本: $(go version)"

# 下载依赖
echo ""
echo "📦 下载 Go 依赖..."
go mod download
go mod tidy

echo ""
echo "✅ 依赖下载完成"

# 检查 Docker 是否安装
if command -v docker &> /dev/null; then
    echo ""
    echo "🐳 检测到 Docker"
    read -p "是否启动数据库容器? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker-compose up -d postgres redis
        echo "⏳ 等待数据库启动..."
        sleep 5
        echo "✅ 数据库容器已启动"
    fi
else
    echo "⚠️  未检测到 Docker。请手动配置 PostgreSQL 和 Redis。"
fi

# 创建 .env 文件
if [ ! -f .env ]; then
    echo ""
    echo "📝 创建 .env 配置文件..."
    cp .env.example .env
    echo "✅ .env 文件已创建，请根据实际情况修改配置"
else
    echo ""
    echo "ℹ️  .env 文件已存在"
fi

echo ""
echo "✨ 设置完成！"
echo ""
echo "下一步："
echo "  1. 编辑 .env 文件，配置数据库连接等信息"
echo "  2. 运行 'make run' 启动服务"
echo "  3. 访问 http://localhost:8080/api/v1/health 测试"
echo ""
echo "更多信息请查看 docs/API_USAGE.md"
