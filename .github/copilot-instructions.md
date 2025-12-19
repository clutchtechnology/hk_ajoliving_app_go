# Copilot Instructions - AJO Living 房产平台 (后端 API)

## 项目概述

这是 **AJO Living** 房产服务平台的后端 API 服务，使用 Go 语言开发，为 Flutter Web 前端提供 RESTful API 支持。目标市场为香港地区，提供租房、买房、卖房等功能。

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin / Echo / Fiber (推荐 Gin)
- **ORM**: GORM
- **数据库**: PostgreSQL (主数据库) / Redis (缓存)
- **认证**: JWT (JSON Web Token)
- **API 文档**: Swagger / OpenAPI 3.0
- **日志**: Zap / Logrus
- **配置管理**: Viper
- **验证**: go-playground/validator
- **测试**: Go testing + testify
- **容器化**: Docker + Docker Compose

## 项目结构

```
.
├── cmd/                         # 应用程序入口
│   └── api/
│       └── main.go              # API 服务入口
├── internal/                    # 私有应用代码
│   ├── config/                  # 配置
│   │   └── config.go
│   ├── handler/                 # HTTP 处理器 (Controller 层)
│   │   ├── auth_handler.go
│   │   ├── property_handler.go
│   │   ├── user_handler.go
│   │   ├── transaction_handler.go
│   │   ├── valuation_handler.go
│   │   ├── furniture_handler.go
│   │   ├── mortgage_handler.go
│   │   ├── news_handler.go
│   │   ├── school_handler.go
│   │   ├── agent_handler.go
│   │   └── price_index_handler.go
│   ├── service/                 # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── property_service.go
│   │   ├── user_service.go
│   │   ├── transaction_service.go
│   │   ├── valuation_service.go
│   │   ├── furniture_service.go
│   │   ├── mortgage_service.go
│   │   ├── news_service.go
│   │   ├── school_service.go
│   │   ├── agent_service.go
│   │   └── price_index_service.go
│   ├── repository/              # 数据访问层
│   │   ├── property_repo.go
│   │   ├── user_repo.go
│   │   ├── transaction_repo.go
│   │   └── ...
│   ├── model/                   # 数据模型
│   │   ├── property.go
│   │   ├── user.go
│   │   ├── transaction.go
│   │   ├── valuation.go
│   │   ├── furniture.go
│   │   ├── news.go
│   │   ├── school.go
│   │   ├── agent.go
│   │   └── price_index.go
│   ├── dto/                     # 数据传输对象
│   │   ├── request/
│   │   │   ├── property_request.go
│   │   │   ├── auth_request.go
│   │   │   └── ...
│   │   └── response/
│   │       ├── property_response.go
│   │       ├── auth_response.go
│   │       └── ...
│   ├── middleware/              # 中间件
│   │   ├── auth.go              # JWT 认证中间件
│   │   ├── cors.go              # CORS 中间件
│   │   ├── logger.go            # 日志中间件
│   │   ├── ratelimit.go         # 限流中间件
│   │   └── recovery.go          # 异常恢复中间件
│   ├── router/                  # 路由配置
│   │   └── router.go
│   └── pkg/                     # 内部共享包
│       ├── errors/              # 错误处理
│       ├── response/            # 统一响应格式
│       ├── utils/               # 工具函数
│       └── validator/           # 自定义验证器
├── pkg/                         # 公共包（可被外部引用）
│   ├── logger/
│   ├── database/
│   └── cache/
├── migrations/                  # 数据库迁移文件
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   └── ...
├── docs/                        # API 文档
│   └── swagger/
├── scripts/                     # 脚本文件
│   ├── build.sh
│   └── deploy.sh
├── deployments/                 # 部署配置
│   ├── docker/
│   │   └── Dockerfile
│   └── k8s/
├── configs/                     # 配置文件
│   ├── config.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 开发流程规范

### 数据库开发策略

**先完成代码开发，最后统一建数据库**

1. **代码优先**：先完成所有 API 的 Handler、Service、Repository、Model 代码
2. **Model 定义**：在 `internal/model/` 中定义完整的数据模型
3. **AutoMigrate**：开发阶段使用 GORM 的 `AutoMigrate` 自动创建/更新表结构
4. **迁移文件**：所有 API 开发完成后，再统一编写 `migrations/` 迁移文件
5. **生产部署**：生产环境使用手动迁移文件，确保精确控制 schema 变更

```go
// main.go 中的 AutoMigrate 用于开发阶段
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.User{},
        &model.Property{},
        // ... 其他 model
    )
}
```

**优点**：
- 避免频繁修改迁移文件
- Model 变更时自动同步表结构
- 开发效率更高

---

## 编码规范

### 命名约定

- **包名**: 使用小写单词，如 `handler`, `service`, `repository`
- **文件名**: 使用 snake_case，如 `property_handler.go`
- **结构体**: 使用 PascalCase，如 `PropertyService`
- **接口**: 使用 PascalCase，以 `er` 结尾或描述性命名，如 `PropertyRepository`
- **方法**: 使用 PascalCase（公开）或 camelCase（私有）
- **常量**: 使用 PascalCase 或 SCREAMING_SNAKE_CASE
- **变量**: 使用 camelCase

### 代码结构规范

```go
// service/property_service.go
package service

// PropertyService Methods:
// 0. NewPropertyService(repo repository.PropertyRepository, cache cache.Cache, logger *zap.Logger) -> 注入依赖
// 1. ListProperties(ctx context.Context, filter *dto.PropertyFilter) -> 获取房产列表
// 2. GetProperty(ctx context.Context, id uint) -> 获取单个房产详情
// 3. CreateProperty(ctx context.Context, req *dto.CreatePropertyRequest) -> 创建房产
// 4. UpdateProperty(ctx context.Context, id uint, req *dto.UpdatePropertyRequest) -> 更新房产信息
// 5. DeleteProperty(ctx context.Context, id uint) -> 删除房产

import (
    // 标准库
    "context"
    "errors"
    
    // 第三方库
    "go.uber.org/zap"
    
    // 项目内部包
    "github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
    "github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
    "github.com/clutchtechnology/hk_ajoliving_app_go/pkg/cache"
)

// PropertyServiceInterface 定义房产服务接口
type PropertyServiceInterface interface {
    ListProperties(ctx context.Context, filter *dto.PropertyFilter) ([]*model.Property, error)   // 1. 获取房产列表
    GetProperty(ctx context.Context, id uint) (*model.Property, error)                           // 2. 获取单个房产详情
    CreateProperty(ctx context.Context, req *dto.CreatePropertyRequest) (*model.Property, error) // 3. 创建房产
    UpdateProperty(ctx context.Context, id uint, req *dto.UpdatePropertyRequest) (*model.Property, error) // 4. 更新房产信息
    DeleteProperty(ctx context.Context, id uint) error                                           // 5. 删除房产
}

// PropertyService 房产服务
type PropertyService struct {
    repo   repository.PropertyRepository
    cache  cache.Cache
    logger *zap.Logger
}

// 0. NewPropertyService 构造函数
func NewPropertyService(repo repository.PropertyRepository, cache cache.Cache, logger *zap.Logger) *PropertyService {
    return &PropertyService{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}

// 1. ListProperties 获取房产列表
func (s *PropertyService) ListProperties(ctx context.Context, filter *dto.PropertyFilter) ([]*model.Property, error) {
    // 业务逻辑
}

// 2. GetProperty 获取单个房产详情
func (s *PropertyService) GetProperty(ctx context.Context, id uint) (*model.Property, error) {
    // 业务逻辑
}

// 3. CreateProperty 创建房产
func (s *PropertyService) CreateProperty(ctx context.Context, req *dto.CreatePropertyRequest) (*model.Property, error) {
    // 业务逻辑
}

// 4. UpdateProperty 更新房产信息
func (s *PropertyService) UpdateProperty(ctx context.Context, id uint, req *dto.UpdatePropertyRequest) (*model.Property, error) {
    // 业务逻辑
}

// 5. DeleteProperty 删除房产
func (s *PropertyService) DeleteProperty(ctx context.Context, id uint) error {
    // 业务逻辑
}
```

### Handler 规范

```go
// handler/property_handler.go
package handler

// PropertyHandler Methods:
// 0. NewPropertyHandler(service service.PropertyService) -> 注入 PropertyService
// 1. ListProperties(c *gin.Context) -> 获取房产列表
// 2. GetProperty(c *gin.Context) -> 获取单个房产详情
// 3. CreateProperty(c *gin.Context) -> 创建房产
// 4. UpdateProperty(c *gin.Context) -> 更新房产信息
// 5. DeleteProperty(c *gin.Context) -> 删除房产

type PropertyHandler struct {
    service service.PropertyService
}

// 0. NewPropertyHandler -> 注入 PropertyService
func NewPropertyHandler(service service.PropertyService) *PropertyHandler {
    return &PropertyHandler{service: service}
}

// 1. ListProperties -> 获取房产列表
func (h *PropertyHandler) ListProperties(c *gin.Context) {
    // 参数绑定与验证
    var req dto.ListPropertiesRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    
    // 调用服务层
    properties, err := h.service.ListProperties(c.Request.Context(), &req)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    
    // 返回响应
    response.Success(c, properties)
}

// 2. GetProperty -> 获取单个房产详情
func (h *PropertyHandler) GetProperty(c *gin.Context) {
    // 解析路径参数
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.BadRequest(c, "invalid property id")
        return
    }
    
    // 调用服务层
    property, err := h.service.GetProperty(c.Request.Context(), uint(id))
    if err != nil {
        if errors.Is(err, errors.ErrNotFound) {
            response.NotFound(c, "property not found")
            return
        }
        response.InternalError(c, err.Error())
        return
    }
    
    // 返回响应
    response.Success(c, property)
}

// 3. CreateProperty -> 创建房产
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
    // 参数绑定与验证
    var req dto.CreatePropertyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    
    // 获取当前用户
    userID := c.GetUint("user_id")
    
    // 调用服务层
    property, err := h.service.CreateProperty(c.Request.Context(), userID, &req)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    
    // 返回响应
    response.Created(c, property)
}
```

### 统一响应格式

```go
// pkg/response/response.go
package response

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
    Response
    Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
    Page      int   `json:"page"`
    PageSize  int   `json:"page_size"`
    Total     int64 `json:"total"`
    TotalPage int   `json:"total_page"`
}

// 响应码定义
const (
    CodeSuccess       = 0
    CodeBadRequest    = 400
    CodeUnauthorized  = 401
    CodeForbidden     = 403
    CodeNotFound      = 404
    CodeInternalError = 500
)

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    CodeSuccess,
        Message: "success",
        Data:    data,
    })
}

func BadRequest(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, Response{
        Code:    CodeBadRequest,
        Message: message,
    })
}
```

### 错误处理规范

```go
// internal/pkg/errors/errors.go
package errors

import "errors"

var (
    ErrNotFound       = errors.New("resource not found")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrForbidden      = errors.New("forbidden")
    ErrInvalidInput   = errors.New("invalid input")
    ErrAlreadyExists  = errors.New("resource already exists")
    ErrInternalServer = errors.New("internal server error")
)

// 自定义业务错误
type BusinessError struct {
    Code    int
    Message string
    Err     error
}

func (e *BusinessError) Error() string {
    return e.Message
}

func NewBusinessError(code int, message string) *BusinessError {
    return &BusinessError{
        Code:    code,
        Message: message,
    }
}
```

### 数据模型规范

```go
// internal/model/property.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type Property struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Title       string         `gorm:"size:255;not null" json:"title"`
    Description string         `gorm:"type:text" json:"description"`
    Price       float64        `gorm:"not null" json:"price"`
    Area        float64        `json:"area"`                          // 面积（平方尺）
    Bedrooms    int            `json:"bedrooms"`                      // 卧室数
    Bathrooms   int            `json:"bathrooms"`                     // 浴室数
    PropertyType string        `gorm:"size:50" json:"property_type"` // 物业类型
    ListingType  string        `gorm:"size:20" json:"listing_type"`  // 租/售
    Address     string         `gorm:"size:500" json:"address"`
    DistrictID  uint           `json:"district_id"`
    EstateID    *uint          `json:"estate_id"`
    AgentID     uint           `json:"agent_id"`
    Status      string         `gorm:"size:20;default:'active'" json:"status"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    // 关联
    District    *District      `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
    Estate      *Estate        `gorm:"foreignKey:EstateID" json:"estate,omitempty"`
    Agent       *Agent         `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
    Images      []PropertyImage `gorm:"foreignKey:PropertyID" json:"images,omitempty"`
    Facilities  []Facility     `gorm:"many2many:property_facilities" json:"facilities,omitempty"`
}

func (Property) TableName() string {
    return "properties"
}
```

### DTO 规范

```go
// internal/dto/request/property_request.go
package request

type ListPropertiesRequest struct {
    DistrictID   *uint    `form:"district_id"`
    EstateID     *uint    `form:"estate_id"`
    MinPrice     *float64 `form:"min_price"`
    MaxPrice     *float64 `form:"max_price"`
    MinArea      *float64 `form:"min_area"`
    MaxArea      *float64 `form:"max_area"`
    Bedrooms     *int     `form:"bedrooms"`
    PropertyType *string  `form:"property_type"`
    ListingType  *string  `form:"listing_type" binding:"omitempty,oneof=rent sale"`
    Page         int      `form:"page,default=1" binding:"min=1"`
    PageSize     int      `form:"page_size,default=20" binding:"min=1,max=100"`
    SortBy       string   `form:"sort_by,default=created_at"`
    SortOrder    string   `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

type CreatePropertyRequest struct {
    Title        string   `json:"title" binding:"required,max=255"`
    Description  string   `json:"description"`
    Price        float64  `json:"price" binding:"required,gt=0"`
    Area         float64  `json:"area" binding:"required,gt=0"`
    Bedrooms     int      `json:"bedrooms" binding:"min=0"`
    Bathrooms    int      `json:"bathrooms" binding:"min=0"`
    PropertyType string   `json:"property_type" binding:"required"`
    ListingType  string   `json:"listing_type" binding:"required,oneof=rent sale"`
    Address      string   `json:"address" binding:"required"`
    DistrictID   uint     `json:"district_id" binding:"required"`
    EstateID     *uint    `json:"estate_id"`
    FacilityIDs  []uint   `json:"facility_ids"`
    ImageURLs    []string `json:"image_urls"`
}
```

### Repository 规范

```go
// internal/repository/property_repo.go
package repository

import (
    "context"
    "gorm.io/gorm"
)

type PropertyRepository interface {
    List(ctx context.Context, filter *dto.ListPropertiesRequest) ([]*model.Property, int64, error)
    GetByID(ctx context.Context, id uint) (*model.Property, error)
    Create(ctx context.Context, property *model.Property) error
    Update(ctx context.Context, property *model.Property) error
    Delete(ctx context.Context, id uint) error
}

type propertyRepository struct {
    db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) PropertyRepository {
    return &propertyRepository{db: db}
}

func (r *propertyRepository) List(ctx context.Context, filter *dto.ListPropertiesRequest) ([]*model.Property, int64, error) {
    var properties []*model.Property
    var total int64
    
    query := r.db.WithContext(ctx).Model(&model.Property{})
    
    // 应用筛选条件
    if filter.DistrictID != nil {
        query = query.Where("district_id = ?", *filter.DistrictID)
    }
    if filter.MinPrice != nil {
        query = query.Where("price >= ?", *filter.MinPrice)
    }
    if filter.MaxPrice != nil {
        query = query.Where("price <= ?", *filter.MaxPrice)
    }
    // ... 更多筛选条件
    
    // 统计总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 分页和排序
    offset := (filter.Page - 1) * filter.PageSize
    query = query.Offset(offset).Limit(filter.PageSize)
    query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder))
    
    // 预加载关联
    query = query.Preload("District").Preload("Estate").Preload("Images")
    
    if err := query.Find(&properties).Error; err != nil {
        return nil, 0, err
    }
    
    return properties, total, nil
}
```

---

## 中间件配置

### JWT 认证中间件

```go
// internal/middleware/auth.go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            response.Unauthorized(c, "missing authorization header")
            c.Abort()
            return
        }
        
        // 解析 Bearer token
        if !strings.HasPrefix(token, "Bearer ") {
            response.Unauthorized(c, "invalid token format")
            c.Abort()
            return
        }
        
        tokenString := strings.TrimPrefix(token, "Bearer ")
        claims, err := jwt.ParseToken(tokenString)
        if err != nil {
            response.Unauthorized(c, "invalid token")
            c.Abort()
            return
        }
        
        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}
```

### CORS 中间件

```go
// internal/middleware/cors.go
func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // 生产环境应限制域名
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

---

## 数据库设计

### 主要数据表

| 表名 | 说明 |
|------|------|
| `users` | 用户表 |
| `properties` | 房产表 |
| `property_images` | 房产图片表 |
| `districts` | 地区表 |
| `estates` | 屋苑表 |
| `transactions` | 成交记录表 |
| `agents` | 代理人表 |
| `agencies` | 代理公司表 |
| `schools` | 学校表 |
| `school_nets` | 校网表 |
| `furniture` | 家具商品表 |
| `furniture_categories` | 家具分类表 |
| `orders` | 订单表 |
| `order_items` | 订单项表 |
| `cart_items` | 购物车项表 |
| `favorites` | 收藏表 |
| `news` | 新闻表 |
| `news_categories` | 新闻分类表 |
| `price_indices` | 楼价指数表 |
| `mortgage_rates` | 按揭利率表 |
| `facilities` | 设施配套表 |
| `property_facilities` | 房产设施关联表 |

---

## 配置管理

### 配置文件结构

```yaml
# configs/config.yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # debug, release, test

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "ajoliving"
  password: ""
  name: "ajoliving_db"
  sslmode: "disable"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expire_hours: 24
  refresh_expire_hours: 168

log:
  level: "debug"
  format: "json"
  output: "stdout"

cors:
  allowed_origins:
    - "http://localhost:3000"
    - "https://ajoliving.com"
```

### 环境变量

```bash
# .env.example
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=debug

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=ajoliving
DB_PASSWORD=secret
DB_NAME=ajoliving_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-super-secret-key
JWT_EXPIRE_HOURS=24

# 第三方服务
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
FACEBOOK_APP_ID=
FACEBOOK_APP_SECRET=
```

---

## 测试规范

### 单元测试

```go
// internal/service/property_service_test.go
package service_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestPropertyService_ListProperties(t *testing.T) {
    // 设置 mock
    mockRepo := new(mocks.MockPropertyRepository)
    mockCache := new(mocks.MockCache)
    logger := zap.NewNop()
    
    service := NewPropertyService(mockRepo, mockCache, logger)
    
    // 定义测试用例
    tests := []struct {
        name    string
        filter  *dto.ListPropertiesRequest
        want    []*model.Property
        wantErr bool
    }{
        {
            name: "成功获取房产列表",
            filter: &dto.ListPropertiesRequest{
                Page:     1,
                PageSize: 10,
            },
            want:    []*model.Property{{ID: 1, Title: "Test Property"}},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo.On("List", mock.Anything, tt.filter).Return(tt.want, int64(1), nil)
            
            got, err := service.ListProperties(context.Background(), tt.filter)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

### 集成测试

```go
// tests/integration/property_test.go
func TestPropertyAPI_Integration(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    router := setupTestRouter(db)
    
    t.Run("GET /api/v1/properties", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/v1/properties", nil)
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
    })
}
```

---

## 部署

### Docker

```dockerfile
# deployments/docker/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Hong_Kong

WORKDIR /app

COPY --from=builder /api .
COPY configs/config.prod.yaml ./configs/config.yaml

EXPOSE 8080

CMD ["./api"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: ajoliving
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: ajoliving_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

---

## 常用命令

```bash
# 运行开发服务器
go run cmd/api/main.go

# 构建
go build -o bin/api cmd/api/main.go

# 运行测试
go test ./...

# 运行测试（带覆盖率）
go test -cover ./...

# 生成 Swagger 文档
swag init -g cmd/api/main.go -o docs/swagger

# 数据库迁移
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

# 代码格式化
go fmt ./...

# 代码检查
golangci-lint run

# 生成 Mock
mockgen -source=internal/repository/property_repo.go -destination=internal/mocks/property_repo_mock.go
```

---

## Makefile

```makefile
.PHONY: build run test clean migrate swagger lint

# 变量
APP_NAME=ajoliving-api
BUILD_DIR=bin

# 构建
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

# 运行
run:
	go run ./cmd/api

# 测试
test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# 清理
clean:
	rm -rf $(BUILD_DIR)

# 数据库迁移
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

# Swagger 文档
swagger:
	swag init -g cmd/api/main.go -o docs/swagger

# 代码检查
lint:
	golangci-lint run

# Docker
docker-build:
	docker build -f deployments/docker/Dockerfile -t $(APP_NAME) .

docker-run:
	docker-compose up -d
```

---

## Git 提交规范

使用 Conventional Commits：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具相关

示例：`feat(property): add property listing API`

---
