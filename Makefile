.PHONY: build run test clean migrate swagger lint docker-build docker-run help

# 变量
APP_NAME=ajoliving-api
BUILD_DIR=bin
MAIN_FILE=./cmd/api/main.go

# Go 相关
GO=go
GOFLAGS=-v

# 默认目标
.DEFAULT_GOAL := help

## build: 构建应用程序
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

## run: 运行应用程序
run:
	@echo "Running $(APP_NAME)..."
	$(GO) run $(MAIN_FILE)

## test: 运行所有测试
test:
	@echo "Running tests..."
	$(GO) test -v ./...

## test-cover: 运行测试并生成覆盖率报告
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## deps: 下载依赖
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

## lint: 运行代码检查
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: brew install golangci-lint"; \
	fi

## fmt: 格式化代码
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## swagger: 生成 Swagger 文档
swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_FILE) -o docs/swagger; \
	else \
		echo "swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

## migrate-up: 执行数据库迁移
migrate-up:
	@echo "Running migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL is not set"; \
		exit 1; \
	fi
	migrate -path migrations -database "$(DATABASE_URL)" up

## migrate-down: 回滚数据库迁移
migrate-down:
	@echo "Rolling back migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "DATABASE_URL is not set"; \
		exit 1; \
	fi
	migrate -path migrations -database "$(DATABASE_URL)" down

## migrate-create: 创建新的迁移文件
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)

## docker-build: 构建 Docker 镜像
docker-build:
	@echo "Building Docker image..."
	docker-compose build

## docker-up: 启动所有 Docker 服务
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

## docker-down: 停止并删除 Docker 服务
docker-down:
	@echo "Stopping and removing services..."
	docker-compose down

## docker-stop: 停止 Docker 服务
docker-stop:
	@echo "Stopping services..."
	docker-compose stop

## docker-restart: 重启 Docker 服务
docker-restart:
	@echo "Restarting services..."
	docker-compose restart

## docker-logs: 查看 Docker 日志
docker-logs:
	@echo "Following logs (Ctrl+C to exit)..."
	docker-compose logs -f

## docker-logs-api: 查看 API 日志
docker-logs-api:
	docker-compose logs -f api

## docker-ps: 查看 Docker 容器状态
docker-ps:
	docker-compose ps

## docker-clean: 清理 Docker 资源（不删除数据卷）
docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down
	docker image prune -f

## docker-clean-all: 清理所有 Docker 资源（包括数据卷）
docker-clean-all:
	@echo "WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker image prune -f; \
		echo "All Docker resources cleaned"; \
	else \
		echo "Cancelled"; \
	fi

## docker-shell-api: 进入 API 容器
docker-shell-api:
	docker-compose exec api sh

## docker-shell-db: 进入数据库容器
docker-shell-db:
	docker-compose exec postgres psql -U ajoliving -d ajoliving_db

## docker-deploy: 一键部署到 Docker
docker-deploy:
	@chmod +x scripts/deploy_docker.sh
	@./scripts/deploy_docker.sh

## help: 显示帮助信息
help:
	@echo "AJO Living API - Available Commands:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
