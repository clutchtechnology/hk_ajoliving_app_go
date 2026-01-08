#!/bin/bash

# 数据库初始化脚本

echo "🔧 初始化数据库..."

# 创建数据库用户
docker exec ajoliving-postgres psql -U postgres -c "CREATE USER ajoliving WITH PASSWORD 'secret';" 2>/dev/null || echo "✅ 用户 ajoliving 已存在"

# 创建数据库
docker exec ajoliving-postgres psql -U postgres -c "CREATE DATABASE ajoliving_db OWNER ajoliving;" 2>/dev/null || echo "✅ 数据库 ajoliving_db 已存在"

# 授予权限
docker exec ajoliving-postgres psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE ajoliving_db TO ajoliving;" 2>/dev/null

# 授予 schema 权限（PostgreSQL 15+ 需要）
docker exec ajoliving-postgres psql -U postgres -d ajoliving_db -c "GRANT ALL ON SCHEMA public TO ajoliving;" 2>/dev/null

echo "✅ 数据库初始化完成！"
