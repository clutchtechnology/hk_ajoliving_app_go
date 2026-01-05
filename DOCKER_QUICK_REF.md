# Docker 快速参考

## 一键启动

```bash
./scripts/quick_start.sh
```

## 常用命令速查

### 使用部署脚本
```bash
./scripts/deploy_docker.sh build     # 构建镜像
./scripts/deploy_docker.sh up        # 启动服务
./scripts/deploy_docker.sh down      # 停止服务
./scripts/deploy_docker.sh status    # 查看状态
./scripts/deploy_docker.sh logs      # 查看所有日志
./scripts/deploy_docker.sh logs api  # 仅看 API 日志
./scripts/deploy_docker.sh exec api  # 进入 API 容器
```

### 使用 Makefile
```bash
make docker-up           # 启动
make docker-down         # 停止
make docker-logs         # 查看日志
make docker-ps           # 查看状态
make docker-shell-api    # 进入 API
make docker-shell-db     # 进入数据库
```

### 使用 Docker Compose
```bash
docker-compose up -d         # 后台启动
docker-compose down          # 停止
docker-compose ps            # 状态
docker-compose logs -f api   # 实时日志
docker-compose restart api   # 重启 API
```

## 服务地址

- API: http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## 测试

```bash
./scripts/test_api.sh
# 或
curl http://localhost:8080/health
```

## 调试

```bash
# 查看 API 日志
docker-compose logs -f api

# 进入 API 容器
docker-compose exec api sh

# 连接数据库
docker-compose exec postgres psql -U ajoliving -d ajoliving_db

# 连接 Redis
docker-compose exec redis redis-cli
```

## 完整文档

详见 `docs/DOCKER_DEPLOYMENT.md`
