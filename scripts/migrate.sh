#!/bin/bash

# 数据库迁移脚本

set -e

# 加载环境变量
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# 设置默认值
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-ajoliving}
DB_PASSWORD=${DB_PASSWORD:-}
DB_NAME=${DB_NAME:-ajoliving_db}

# 构建连接字符串
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "📊 数据库迁移工具"
echo "连接到: ${DB_HOST}:${DB_PORT}/${DB_NAME}"

# 检查 migrate 工具是否安装
if ! command -v migrate &> /dev/null; then
    echo "❌ 未检测到 migrate 工具"
    echo "安装方法:"
    echo "  macOS: brew install golang-migrate"
    echo "  Linux: 访问 https://github.com/golang-migrate/migrate"
    exit 1
fi

# 执行迁移
case "$1" in
    up)
        echo "⬆️  执行向上迁移..."
        migrate -path migrations -database "$DATABASE_URL" up
        echo "✅ 迁移完成"
        ;;
    down)
        echo "⬇️  执行向下迁移..."
        migrate -path migrations -database "$DATABASE_URL" down
        echo "✅ 回滚完成"
        ;;
    force)
        if [ -z "$2" ]; then
            echo "❌ 请指定版本号: ./scripts/migrate.sh force <version>"
            exit 1
        fi
        echo "🔧 强制设置版本为 $2..."
        migrate -path migrations -database "$DATABASE_URL" force $2
        echo "✅ 版本设置完成"
        ;;
    version)
        echo "📌 当前数据库版本:"
        migrate -path migrations -database "$DATABASE_URL" version
        ;;
    *)
        echo "用法: $0 {up|down|force <version>|version}"
        exit 1
        ;;
esac
