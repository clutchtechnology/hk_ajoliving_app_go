#!/bin/bash

# AJO Living API - Docker 部署脚本
# 用于快速部署和管理 Docker 容器

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印信息
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    info "Docker 环境检查通过"
}

# 创建 .env 文件
create_env() {
    if [ ! -f .env ]; then
        info "创建 .env 文件..."
        cp .env.example .env
        warn "请编辑 .env 文件，设置正确的环境变量（特别是 JWT_SECRET）"
        read -p "按回车键继续..."
    else
        info ".env 文件已存在"
    fi
}

# 构建镜像
build() {
    info "构建 Docker 镜像..."
    docker-compose build --no-cache
    info "镜像构建完成"
}

# 启动服务
up() {
    info "启动服务..."
    docker-compose up -d
    info "等待服务启动..."
    sleep 5
    
    # 检查服务状态
    docker-compose ps
    
    info "服务已启动！"
    info "API 地址: http://localhost:8080"
    info "PostgreSQL: localhost:5432"
    info "Redis: localhost:6379"
}

# 停止服务
down() {
    info "停止服务..."
    docker-compose down
    info "服务已停止"
}

# 重启服务
restart() {
    info "重启服务..."
    docker-compose restart
    info "服务已重启"
}

# 查看日志
logs() {
    if [ -z "$1" ]; then
        docker-compose logs -f
    else
        docker-compose logs -f "$1"
    fi
}

# 查看服务状态
status() {
    info "服务状态："
    docker-compose ps
}

# 进入容器
exec_container() {
    if [ -z "$1" ]; then
        error "请指定容器名称: api, postgres, redis"
        exit 1
    fi
    
    case "$1" in
        api)
            docker-compose exec api sh
            ;;
        postgres)
            docker-compose exec postgres psql -U ajoliving -d ajoliving_db
            ;;
        redis)
            docker-compose exec redis redis-cli
            ;;
        *)
            error "未知的容器: $1"
            exit 1
            ;;
    esac
}

# 运行数据库迁移
migrate() {
    info "运行数据库迁移..."
    docker-compose exec api ./api migrate
    info "数据库迁移完成"
}

# 清理所有数据
clean() {
    warn "警告：此操作将删除所有容器和数据卷！"
    read -p "确认继续？(yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        info "停止并删除容器..."
        docker-compose down -v
        info "清理完成"
    else
        info "取消操作"
    fi
}

# 显示帮助信息
show_help() {
    cat << EOF
AJO Living API - Docker 部署脚本

用法: $0 [命令]

命令:
  build       构建 Docker 镜像
  up          启动所有服务
  down        停止所有服务
  restart     重启所有服务
  logs        查看日志 (可选参数: api, postgres, redis)
  status      查看服务状态
  exec        进入容器 (参数: api, postgres, redis)
  migrate     运行数据库迁移
  clean       清理所有容器和数据（危险操作）
  help        显示此帮助信息

示例:
  $0 build              # 构建镜像
  $0 up                 # 启动服务
  $0 logs api           # 查看 API 日志
  $0 exec postgres      # 进入 PostgreSQL 容器
  $0 down               # 停止服务

EOF
}

# 主函数
main() {
    check_docker
    
    case "${1:-help}" in
        build)
            create_env
            build
            ;;
        up)
            create_env
            up
            ;;
        down)
            down
            ;;
        restart)
            restart
            ;;
        logs)
            logs "$2"
            ;;
        status)
            status
            ;;
        exec)
            exec_container "$2"
            ;;
        migrate)
            migrate
            ;;
        clean)
            clean
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
