# Docker 部署指南

本文档详细说明如何使用 Docker 部署 AJO Living API 后端服务。

## 目录

- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [详细步骤](#详细步骤)
- [常用命令](#常用命令)
- [配置说明](#配置说明)
- [故障排除](#故障排除)
- [生产环境部署](#生产环境部署)

---

## 前置要求

### 1. 安装 Docker

**macOS:**
```bash
# 使用 Homebrew 安装
brew install --cask docker

# 或下载 Docker Desktop for Mac
# https://www.docker.com/products/docker-desktop
```

**Linux (Ubuntu/Debian):**
```bash
# 更新包索引
sudo apt-get update

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 添加当前用户到 docker 组
sudo usermod -aG docker $USER

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

**Windows:**
- 下载并安装 [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)

### 2. 验证安装

```bash
docker --version
docker-compose --version
```

---

## 快速开始

### 方式一：一键启动（推荐）

```bash
# 1. 进入项目目录
cd /path/to/hk_ajoliving_app_go

# 2. 赋予脚本执行权限
chmod +x scripts/*.sh

# 3. 运行快速启动脚本
./scripts/quick_start.sh
```

等待几分钟，服务将自动：
- ✅ 创建环境变量文件
- ✅ 构建 Docker 镜像
- ✅ 启动所有服务（API、PostgreSQL、Redis）
- ✅ 等待数据库就绪

启动成功后，访问：
- **API 服务**: http://localhost:8080
- **健康检查**: http://localhost:8080/health
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

### 方式二：使用部署脚本

```bash
# 构建镜像
./scripts/deploy_docker.sh build

# 启动服务
./scripts/deploy_docker.sh up

# 查看状态
./scripts/deploy_docker.sh status
```

### 方式三：手动操作

```bash
# 1. 创建环境变量文件
cp .env.example .env

# 2. 构建并启动
docker-compose up -d

# 3. 查看日志
docker-compose logs -f
```

---

## 详细步骤

### 1. 配置环境变量

复制示例文件并编辑：

```bash
cp .env.example .env
```

编辑 `.env` 文件，配置关键参数：

```env
# 数据库密码（生产环境必须修改）
DB_PASSWORD=your-secure-password

# JWT 密钥（生产环境必须修改）
JWT_SECRET=your-super-secret-jwt-key-min-32-characters

# 运行模式（生产环境改为 release）
SERVER_MODE=debug
```

### 2. 构建 Docker 镜像

```bash
# 构建镜像
docker-compose build

# 或使用部署脚本
./scripts/deploy_docker.sh build
```

### 3. 启动服务

```bash
# 后台启动所有服务
docker-compose up -d

# 或使用部署脚本
./scripts/deploy_docker.sh up
```

### 4. 验证服务

```bash
# 查看服务状态
docker-compose ps

# 查看 API 日志
docker-compose logs -f api

# 测试 API
curl http://localhost:8080/health

# 或使用测试脚本
./scripts/test_api.sh
```

### 5. 数据库初始化（如需要）

```bash
# 进入 API 容器
docker-compose exec api sh

# 运行迁移（如果你的应用支持）
# ./api migrate

# 或手动连接数据库
docker-compose exec postgres psql -U ajoliving -d ajoliving_db
```

---

## 常用命令

### 使用部署脚本（推荐）

```bash
# 构建镜像
./scripts/deploy_docker.sh build

# 启动服务
./scripts/deploy_docker.sh up

# 停止服务
./scripts/deploy_docker.sh down

# 重启服务
./scripts/deploy_docker.sh restart

# 查看日志
./scripts/deploy_docker.sh logs          # 所有服务
./scripts/deploy_docker.sh logs api      # 仅 API
./scripts/deploy_docker.sh logs postgres # 仅数据库

# 查看状态
./scripts/deploy_docker.sh status

# 进入容器
./scripts/deploy_docker.sh exec api      # 进入 API 容器
./scripts/deploy_docker.sh exec postgres # 进入数据库
./scripts/deploy_docker.sh exec redis    # 进入 Redis

# 清理所有数据（危险！）
./scripts/deploy_docker.sh clean
```

### 使用 Docker Compose

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f [service_name]

# 查看状态
docker-compose ps

# 进入容器
docker-compose exec api sh
docker-compose exec postgres psql -U ajoliving -d ajoliving_db
docker-compose exec redis redis-cli

# 重新构建镜像
docker-compose build --no-cache

# 扩展服务（如需要）
docker-compose up -d --scale api=3
```

### 测试 API

```bash
# 运行 API 测试脚本
./scripts/test_api.sh

# 手动测试
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/health

# 测试注册接口
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123456",
    "phone": "85298765432"
  }'
```

---

## 配置说明

### Docker Compose 配置

`docker-compose.yml` 包含三个服务：

1. **api**: Go 后端 API 服务
   - 端口: 8080
   - 依赖: postgres, redis

2. **postgres**: PostgreSQL 数据库
   - 端口: 5432
   - 数据卷: postgres_data

3. **redis**: Redis 缓存
   - 端口: 6379
   - 数据卷: redis_data

### Dockerfile 说明

采用多阶段构建：

1. **构建阶段**: 使用 `golang:1.21-alpine` 编译二进制文件
2. **运行阶段**: 使用 `alpine:3.19` 运行，最小化镜像体积

优势：
- ✅ 最终镜像体积小（约 20MB）
- ✅ 安全性高（最小依赖）
- ✅ 构建缓存优化（依赖层分离）

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SERVER_PORT` | API 服务端口 | 8080 |
| `SERVER_MODE` | 运行模式 | debug |
| `DB_HOST` | 数据库主机 | postgres |
| `DB_PORT` | 数据库端口 | 5432 |
| `DB_USER` | 数据库用户 | ajoliving |
| `DB_PASSWORD` | 数据库密码 | secret |
| `DB_NAME` | 数据库名称 | ajoliving_db |
| `REDIS_HOST` | Redis 主机 | redis |
| `REDIS_PORT` | Redis 端口 | 6379 |
| `JWT_SECRET` | JWT 密钥 | （必须修改） |

---

## 故障排除

### 1. 端口被占用

**错误信息:**
```
Error: bind: address already in use
```

**解决方案:**
```bash
# 查找占用端口的进程
lsof -i :8080
lsof -i :5432
lsof -i :6379

# 杀死进程或修改 docker-compose.yml 中的端口映射
```

### 2. 容器无法启动

```bash
# 查看容器日志
docker-compose logs api

# 查看容器详情
docker inspect ajoliving-api

# 重新构建镜像
docker-compose build --no-cache api
```

### 3. 数据库连接失败

```bash
# 检查数据库容器状态
docker-compose ps postgres

# 查看数据库日志
docker-compose logs postgres

# 测试数据库连接
docker-compose exec postgres pg_isready -U ajoliving

# 手动连接数据库
docker-compose exec postgres psql -U ajoliving -d ajoliving_db
```

### 4. Redis 连接失败

```bash
# 检查 Redis 状态
docker-compose ps redis

# 测试 Redis 连接
docker-compose exec redis redis-cli ping

# 查看 Redis 日志
docker-compose logs redis
```

### 5. API 无响应

```bash
# 查看 API 日志
docker-compose logs -f api

# 进入容器检查
docker-compose exec api sh

# 检查进程
docker-compose exec api ps aux

# 测试网络连接
docker-compose exec api wget -O- http://localhost:8080/health
```

### 6. 清理并重新开始

```bash
# 停止所有容器
docker-compose down

# 删除所有数据（包括数据库数据）
docker-compose down -v

# 删除镜像
docker rmi ajoliving-api

# 清理 Docker 系统
docker system prune -a

# 重新构建和启动
./scripts/quick_start.sh
```

---

## 生产环境部署

### 1. 安全配置

```bash
# 修改 .env 文件
SERVER_MODE=release
DB_PASSWORD=strong-random-password-here
JWT_SECRET=very-long-random-secret-at-least-32-chars
```

### 2. 使用生产配置

修改 `docker-compose.yml`：

```yaml
api:
  environment:
    - SERVER_MODE=release
  restart: always
  
postgres:
  environment:
    POSTGRES_PASSWORD: ${DB_PASSWORD}
  restart: always
```

### 3. 使用外部数据库（可选）

```yaml
api:
  environment:
    - DB_HOST=your-production-db.example.com
    - DB_PASSWORD=${PROD_DB_PASSWORD}
  # 移除 depends_on: postgres
```

### 4. 反向代理（Nginx）

```nginx
server {
    listen 80;
    server_name api.ajoliving.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5. HTTPS 配置（Let's Encrypt）

```bash
# 安装 Certbot
sudo apt-get install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d api.ajoliving.com
```

### 6. 监控和日志

```bash
# 使用 Docker stats 监控资源
docker stats

# 集中日志管理
docker-compose logs --tail=100 -f | tee logs/api.log
```

### 7. 备份策略

```bash
# 备份数据库
docker-compose exec -T postgres pg_dump -U ajoliving ajoliving_db > backup.sql

# 备份 Redis
docker-compose exec redis redis-cli SAVE
docker cp ajoliving-redis:/data/dump.rdb ./backup/

# 定期备份脚本
# 添加到 crontab
0 2 * * * /path/to/backup_script.sh
```

---

## 性能优化

### 1. 资源限制

在 `docker-compose.yml` 中添加：

```yaml
api:
  deploy:
    resources:
      limits:
        cpus: '2'
        memory: 1G
      reservations:
        cpus: '0.5'
        memory: 512M
```

### 2. 数据库连接池

在配置文件中调整：

```yaml
database:
  max_idle_conns: 20
  max_open_conns: 100
  conn_max_lifetime: 3600
```

### 3. Redis 持久化

```yaml
redis:
  command: redis-server --appendonly yes
```

---

## 相关文档

- [API 使用文档](./API_USAGE.md)
- [数据库设计](./DATABASE_DESIGN.md)
- [实现总结](./IMPLEMENTATION_SUMMARY.md)

---

## 技术支持

如遇到问题，请：

1. 查看日志: `docker-compose logs -f`
2. 检查文档: 本文档的[故障排除](#故障排除)部分
3. 提交 Issue: 在项目仓库提交问题

---

**最后更新**: 2026-01-04
