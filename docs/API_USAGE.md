# API 接口使用指南

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 3. 启动数据库（使用 Docker）

```bash
docker-compose up -d postgres redis
```

### 4. 运行数据库迁移

```bash
make migrate-up
```

### 5. 启动服务

```bash
make run
```

服务将在 `http://localhost:8080` 启动。

---

## API 接口测试示例

### 1. 健康检查

```bash
curl http://localhost:8080/api/v1/health
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "healthy",
    "timestamp": 1702896000,
    "service": "AJO Living API"
  }
}
```

### 2. 版本信息

```bash
curl http://localhost:8080/api/v1/version
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version": "1.0.0",
    "api_version": "v1",
    "build_time": "2025-12-18",
    "go_version": "1.21+"
  }
}
```

### 3. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe",
    "phone": "85298765432",
    "user_type": "buyer"
  }'
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "User registered successfully",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "phone": "85298765432",
      "user_type": "buyer",
      "status": "active",
      "created_at": "2025-12-18T10:00:00Z"
    }
  }
}
```

### 4. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "john@example.com",
    "password": "password123"
  }'
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "user_type": "buyer",
      "status": "active"
    }
  }
}
```

### 5. 刷新令牌

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### 6. 忘记密码

```bash
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

### 7. 重置密码

```bash
curl -X POST http://localhost:8080/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "reset_token_here",
    "new_password": "newpassword123"
  }'
```

### 8. 验证码验证

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-code \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "code": "123456",
    "type": "register"
  }'
```

### 9. 用户登出（需要认证）

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 错误响应格式

当请求失败时，API 将返回如下格式的错误响应：

```json
{
  "code": 400,
  "message": "Invalid input",
  "data": null
}
```

常见错误码：
- `400`: 请求参数错误
- `401`: 未授权（未登录或令牌无效）
- `403`: 禁止访问（权限不足）
- `404`: 资源不存在
- `500`: 服务器内部错误

---

## 认证说明

大多数 API 接口需要在请求头中携带 JWT 令牌：

```
Authorization: Bearer <access_token>
```

令牌在用户登录后获得，有效期默认为 24 小时。

---

## 开发工具

### 使用 Postman

1. 导入 API 文档（TODO: 生成 Swagger/OpenAPI 文档）
2. 设置环境变量 `base_url` 为 `http://localhost:8080`
3. 登录后将 `access_token` 保存为环境变量

### 使用 curl

参考上面的示例命令。

### 使用 HTTPie

```bash
# 安装 HTTPie
brew install httpie

# 示例请求
http POST http://localhost:8080/api/v1/auth/login \
  login=john@example.com \
  password=password123
```

---

## 下一步

- [ ] 实现用户模块接口 (11-18)
- [ ] 实现房产模块接口 (21-48)
- [ ] 实现其他业务模块
- [ ] 添加 Swagger 文档
- [ ] 添加单元测试
- [ ] 添加集成测试
- [ ] 实现 Redis 缓存
- [ ] 实现文件上传
- [ ] 实现邮件发送
- [ ] 部署配置
