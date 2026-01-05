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

## 项目结构（极简版 - 微小项目推荐）

```
.
├── controllers/                 # HTTP 控制器（Handler 层）
│   ├── auth_controller.go
│   ├── property_controller.go
│   ├── user_controller.go
│   └── ...
├── services/                    # 业务逻辑层
│   ├── auth_service.go
│   ├── property_service.go
│   └── ...
├── models/                      # 数据模型 + DTO（合并）
│   ├── property.go              # 包含 Model、Request、Response
│   ├── user.go
│   └── ...
├── databases/                   # 数据库相关
│   ├── db.go                    # 数据库连接初始化
│   ├── property_repo.go         # Repository 层
│   ├── user_repo.go
│   └── ...
├── middlewares/                 # HTTP 中间件
│   ├── auth.go
│   ├── cors.go
│   └── logger.go
├── routes/                      # 路由配置
│   └── routes.go
├── tools/                       # 工具函数
│   ├── response.go              # 统一响应
│   ├── errors.go                # 错误定义
│   ├── jwt.go                   # JWT 工具
│   └── password.go              # 密码加密
├── sql/                         # SQL 迁移文件（可选）
│   ├── 001_create_users.sql
│   └── ...
├── container/                   # 依赖注入容器（可选）
│   └── container.go
├── .env                         # 环境变量
├── .env.production              # 生产环境变量
├── main.go                      # 应用入口
├── go.mod
├── go.sum
├── Makefile
├── docker-compose.yml
└── README.md
```

**精简设计原则（奥卡姆剃刀）：**
1. **扁平化目录** - 减少嵌套层级，提高可读性
2. **DTO 不分子目录** - 单文件内包含 Request/Response，减少文件跳转
3. **pkg 统一管理** - 不再区分 `internal/pkg` 和顶层 `pkg`，统一放在顶层
4. **业务优先** - Handler → Service → Repository 三层架构清晰，专注业务开发
5. **延迟优化** - migrations 后期统一编写，开发期使用 AutoMigrate

### 目录重构方案（当前项目 → 极简结构）

**重构映射关系：**

```
当前结构                           →    极简结构
─────────────────────────────────────────────────────────
cmd/api/main.go                   →    main.go
internal/handler/*.go             →    controllers/*.go
internal/service/*.go             →    services/*.go
internal/model/*.go               →    models/*.go (Model定义)
internal/dto/request/*.go         →    models/*.go (合并到Model文件)
internal/dto/response/*.go        →    models/*.go (合并到Model文件)
internal/repository/*.go          →    databases/*_repo.go
internal/middleware/*.go          →    middlewares/*.go
internal/router/router.go         →    routes/routes.go
internal/config/*.go              →    删除（用 .env + godotenv）
internal/pkg/errors/*.go          →    tools/errors.go
internal/pkg/response/*.go        →    tools/response.go
internal/pkg/utils/*.go           →    tools/*.go
internal/pkg/validator/*.go       →    tools/validator.go
pkg/database/*.go                 →    databases/db.go
pkg/cache/*.go                    →    databases/cache.go
pkg/logger/*.go                   →    tools/logger.go
configs/*.yaml                    →    .env + .env.production
migrations/*.sql                  →    sql/*.sql (可选)
deployments/                      →    保持不变
docs/                             →    保持不变
tests/                            →    保持不变
```

**迁移步骤：**

```bash
# 0. 提交当前代码
git add .
git commit -m "backup before refactoring"

# 1. 创建新目录结构
mkdir -p controllers services models databases middlewares routes tools sql

# 2. 移动并重命名文件
# Handler → Controller
cp internal/handler/*.go controllers/
# 批量重命名 *_handler.go → *_controller.go
for f in controllers/*_handler.go; do mv "$f" "${f/_handler.go/_controller.go}"; done

# Service → Services
cp internal/service/*.go services/

# Model (保持原样)
cp internal/model/*.go models/

# Repository → databases/*_repo.go
cp internal/repository/*.go databases/

# Middleware → Middlewares
cp internal/middleware/*.go middlewares/

# Router → Routes
cp internal/router/*.go routes/
mv routes/router.go routes/routes.go

# 合并 tools
cp internal/pkg/errors/*.go tools/
cp internal/pkg/response/*.go tools/
cp internal/pkg/utils/*.go tools/
cp internal/pkg/validator/*.go tools/
cp pkg/logger/*.go tools/

# 数据库初始化
cp pkg/database/*.go databases/
mv databases/database.go databases/db.go

# 主入口
cp cmd/api/main.go main.go

# 3. 更新 import 路径（关键步骤）
# 替换所有 Go 文件中的 import 路径
find . -name "*.go" -type f -exec sed -i '' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/internal/handler|github.com/clutchtechnology/hk_ajoliving_app_go/controllers|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/internal/service|github.com/clutchtechnology/hk_ajoliving_app_go/services|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/internal/model|github.com/clutchtechnology/hk_ajoliving_app_go/models|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository|github.com/clutchtechnology/hk_ajoliving_app_go/databases|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/internal/middleware|github.com/clutchtechnology/hk_ajoliving_app_go/middlewares|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/internal/router|github.com/clutchtechnology/hk_ajoliving_app_go/routes|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/|github.com/clutchtechnology/hk_ajoliving_app_go/tools/|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/pkg/database|github.com/clutchtechnology/hk_ajoliving_app_go/databases|g' \
  -e 's|github.com/clutchtechnology/hk_ajoliving_app_go/pkg/logger|github.com/clutchtechnology/hk_ajoliving_app_go/tools|g' \
  {} +

# 4. 更新包名
# controllers
find controllers -name "*.go" -exec sed -i '' 's|^package handler$|package controllers|g' {} +

# services (已经是 service，保持不变)

# databases
find databases -name "*.go" -exec sed -i '' 's|^package repository$|package databases|g' {} +

# middlewares
find middlewares -name "*.go" -exec sed -i '' 's|^package middleware$|package middlewares|g' {} +

# routes
find routes -name "*.go" -exec sed -i '' 's|^package router$|package routes|g' {} +

# tools (合并多个包到一个)
find tools -name "*.go" -exec sed -i '' \
  -e 's|^package errors$|package tools|g' \
  -e 's|^package response$|package tools|g' \
  -e 's|^package utils$|package tools|g' \
  -e 's|^package validator$|package tools|g' \
  -e 's|^package logger$|package tools|g' \
  {} +

# 5. 创建 .env 文件
cat > .env << 'EOF'
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=ajoliving
DB_PASSWORD=secret
DB_NAME=ajoliving_db
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24
EOF

# 6. 清理旧目录 (谨慎操作)
# rm -rf internal/ pkg/ cmd/ configs/

# 7. 整理依赖
go mod tidy

# 8. 测试编译
go build -o bin/api main.go
```

**注意事项：**
- ⚠️ **必须先备份代码或创建新分支**
- ⚠️ 迁移后需要手动合并 DTO 到 models 文件
- ⚠️ 检查 tools/ 包内是否有命名冲突（多个包合并）
- ⚠️ 更新 Dockerfile 和 Makefile 中的路径
- ⚠️ 运行完整测试确保功能正常
- ⚠️ 使用 IDE 的重构工具自动更新 import
- ⚠️ 迁移后运行 `go mod tidy` 和完整测试

## 开发流程规范（精简版）

### 开发优先级（业务优先原则）

1. **核心业务先行** - 优先完成 Handler → Service → Repository 核心业务逻辑
2. **模型定义驱动** - 先定义 Model，使用 GORM AutoMigrate 自动建表
3. **快速迭代** - 边开发边调试，不过度设计
4. **延迟优化** - 缓存、限流等非核心功能后期添加

### 数据库开发策略（代码优先）

**开发阶段：使用 GORM AutoMigrate**
```go
// cmd/api/main.go 开发阶段
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.User{},
        &model.Property{},
        &model.Estate{},
        // ... 其他 model
    )
}
```

**生产部署：统一编写迁移文件**
- 所有 API 开发完成后，再编写 `migrations/` 迁移文件
- 生产环境使用手动迁移，确保精确控制

**优点：**
- ✅ 避免频繁修改迁移文件
- ✅ Model 变更自动同步表结构
- ✅ 开发效率高，专注业务

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
// services/property_service.go
package services

// PropertyService Methods:
// 0. NewPropertyService(repo *databases.PropertyRepo) -> 注入依赖
// 1. ListProperties(ctx context.Context, filter *models.ListPropertiesRequest) -> 获取房产列表
// 2. GetProperty(ctx context.Context, id uint) -> 获取单个房产详情
// 3. CreateProperty(ctx context.Context, req *models.CreatePropertyRequest) -> 创建房产
// 4. UpdateProperty(ctx context.Context, id uint, req *models.UpdatePropertyRequest) -> 更新房产信息
// 5. DeleteProperty(ctx context.Context, id uint) -> 删除房产

import (
    "context"
    "errors"
    
    "github.com/clutchtechnology/hk_ajoliving_app_go/databases"
    "github.com/clutchtechnology/hk_ajoliving_app_go/models"
    "github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// PropertyService 房产服务
type PropertyService struct {
    repo *databases.PropertyRepo
}

// 0. NewPropertyService 构造函数
func NewPropertyService(repo *databases.PropertyRepo) *PropertyService {
    return &PropertyService{repo: repo}
}

// 1. ListProperties 获取房产列表
func (s *PropertyService) ListProperties(ctx context.Context, filter *models.ListPropertiesRequest) ([]*models.Property, error) {
    return s.repo.FindAll(ctx, filter)
}

// 2. GetProperty 获取单个房产详情
func (s *PropertyService) GetProperty(ctx context.Context, id uint) (*models.Property, error) {
    property, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return property, nil
}

// 3. CreateProperty 创建房产
func (s *PropertyService) CreateProperty(ctx context.Context, req *models.CreatePropertyRequest) (*models.Property, error) {
    property := &models.Property{
        Title:       req.Title,
        Price:       req.Price,
        // ... 映射字段
    }
    if err := s.repo.Create(ctx, property); err != nil {
        return nil, err
    }
    return property, nil
}
```

### Controller 规范

```go
// controllers/property_controller.go
package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/clutchtechnology/hk_ajoliving_app_go/services"
    "github.com/clutchtechnology/hk_ajoliving_app_go/models"
    "github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// PropertyController Methods:
// 0. NewPropertyController(service *services.PropertyService) -> 注入 PropertyService
// 1. ListProperties(c *gin.Context) -> 获取房产列表
// 2. GetProperty(c *gin.Context) -> 获取单个房产详情
// 3. CreateProperty(c *gin.Context) -> 创建房产
// 4. UpdateProperty(c *gin.Context) -> 更新房产信息
// 5. DeleteProperty(c *gin.Context) -> 删除房产

type PropertyController struct {
    service *services.PropertyService
}

// 0. NewPropertyController -> 注入 PropertyService
func NewPropertyController(service *services.PropertyService) *PropertyController {
    return &PropertyController{service: service}
}

// 1. ListProperties -> 获取房产列表
func (ctrl *PropertyController) ListProperties(c *gin.Context) {
    var req models.ListPropertiesRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        tools.BadRequest(c, err.Error())
        return
    }
    
    properties, err := ctrl.service.ListProperties(c.Request.Context(), &req)
    if err != nil {
        tools.InternalError(c, err.Error())
        return
    }
    
    tools.Success(c, properties)
}

// 2. GetProperty -> 获取单个房产详情
func (ctrl *PropertyController) GetProperty(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        tools.BadRequest(c, "invalid property id")
        return
    }
    
    property, err := ctrl.service.GetProperty(c.Request.Context(), uint(id))
    if err != nil {
        if errors.Is(err, tools.ErrNotFound) {
            tools.NotFound(c, "property not found")
            return
        }
        tools.InternalError(c, err.Error())
        return
    }
    
    tools.Success(c, property)
}

// 3. CreateProperty -> 创建房产
func (ctrl *PropertyController) CreateProperty(c *gin.Context) {
    var req models.CreatePropertyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        tools.BadRequest(c, err.Error())
        return
    }
    
    property, err := ctrl.service.CreateProperty(c.Request.Context(), &req)
    if err != nil {
        tools.InternalError(c, err.Error())
        return
    }
    
    tools.Created(c, property)
}
```

### 统一响应格式（精简版）

```go
// pkg/response/response.go
package response

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// Response 标准响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
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

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    CodeSuccess,
        Message: "success",
        Data:    data,
    })
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
    c.JSON(http.StatusCreated, Response{
        Code:    CodeSuccess,
        Message: "created",
        Data:    data,
    })
}

// BadRequest 错误请求
func BadRequest(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, Response{
        Code:    CodeBadRequest,
        Message: message,
    })
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
    c.JSON(http.StatusUnauthorized, Response{
        Code:    CodeUnauthorized,
        Message: message,
    })
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
    c.JSON(http.StatusForbidden, Response{
        Code:    CodeForbidden,
        Message: message,
    })
}

// NotFound 资源未找到
func NotFound(c *gin.Context, message string) {
    c.JSON(http.StatusNotFound, Response{
        Code:    CodeNotFound,
        Message: message,
    })
}

// InternalError 服务器内部错误
func InternalError(c *gin.Context, message string) {
    c.JSON(http.StatusInternalServerError, Response{
        Code:    CodeInternalError,
        Message: message,
    })
}
```

### 错误处理规范（精简版）

```go
// pkg/errors/errors.go
package errors

import (
    "errors"
    "fmt"
)

// 预定义错误（标准库 errors）
var (
    ErrNotFound       = errors.New("resource not found")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrForbidden      = errors.New("forbidden")
    ErrInvalidInput   = errors.New("invalid input")
    ErrAlreadyExists  = errors.New("resource already exists")
    ErrInternalServer = errors.New("internal server error")
)

// BusinessError 业务错误
type BusinessError struct {
    Code    int
    Message string
    Err     error
}

func (e *BusinessError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *BusinessError) Unwrap() error {
    return e.Err
}

// New 创建业务错误
func New(code int, message string) *BusinessError {
    return &BusinessError{
        Code:    code,
        Message: message,
    }
}

// Wrap 包装底层错误
func Wrap(code int, message string, err error) *BusinessError {
    return &BusinessError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}
```

**使用示例：**
```go
// Service 层
func (s *PropertyService) GetProperty(ctx context.Context, id uint) (*model.Property, error) {
    property, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.ErrNotFound
        }
        return nil, errors.Wrap(500, "failed to get property", err)
    }
    return property, nil
}

// Handler 层
func (h *PropertyHandler) GetProperty(c *gin.Context) {
    property, err := h.service.GetProperty(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, errors.ErrNotFound) {
            response.NotFound(c, "property not found")
            return
        }
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, property)
}
```

### 数据模型规范（Model + DTO 合并）

```go
// models/property.go
package models

import (
    "time"
    "gorm.io/gorm"
)

// ============ GORM Model ============

type Property struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Title       string         `gorm:"size:255;not null" json:"title"`
    Description string         `gorm:"type:text" json:"description"`
    Price       float64        `gorm:"not null" json:"price"`
    Area        float64        `json:"area"`
    Bedrooms    int            `json:"bedrooms"`
    Bathrooms   int            `json:"bathrooms"`
    PropertyType string        `gorm:"size:50" json:"property_type"`
    ListingType  string        `gorm:"size:20" json:"listing_type"`
    Address     string         `gorm:"size:500" json:"address"`
    DistrictID  uint           `json:"district_id"`
    EstateID    *uint          `json:"estate_id"`
    Status      string         `gorm:"size:20;default:'active'" json:"status"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    // 关联
    District    *District      `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
    Estate      *Estate        `gorm:"foreignKey:EstateID" json:"estate,omitempty"`
}

func (Property) TableName() string {
    return "properties"
}

// ============ Request DTO ============

// ListPropertiesRequest 获取房产列表请求
type ListPropertiesRequest struct {
    DistrictID   *uint    `form:"district_id"`
    MinPrice     *float64 `form:"min_price"`
    MaxPrice     *float64 `form:"max_price"`
    PropertyType *string  `form:"property_type"`
    ListingType  *string  `form:"listing_type" binding:"omitempty,oneof=rent sale"`
    Page         int      `form:"page,default=1" binding:"min=1"`
    PageSize     int      `form:"page_size,default=20" binding:"min=1,max=100"`
}

// CreatePropertyRequest 创建房产请求
type CreatePropertyRequest struct {
    Title        string  `json:"title" binding:"required,max=255"`
    Description  string  `json:"description"`
    Price        float64 `json:"price" binding:"required,gt=0"`
    Area         float64 `json:"area" binding:"required,gt=0"`
    Bedrooms     int     `json:"bedrooms" binding:"min=0"`
    PropertyType string  `json:"property_type" binding:"required"`
    ListingType  string  `json:"listing_type" binding:"required,oneof=rent sale"`
    Address      string  `json:"address" binding:"required"`
    DistrictID   uint    `json:"district_id" binding:"required"`
}

// UpdatePropertyRequest 更新房产请求
type UpdatePropertyRequest struct {
    Title       *string  `json:"title" binding:"omitempty,max=255"`
    Description *string  `json:"description"`
    Price       *float64 `json:"price" binding:"omitempty,gt=0"`
    Status      *string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ============ Response DTO ============

// PropertyResponse 房产响应（列表用）
type PropertyResponse struct {
    ID           uint    `json:"id"`
    Title        string  `json:"title"`
    Price        float64 `json:"price"`
    Area         float64 `json:"area"`
    PropertyType string  `json:"property_type"`
    Address      string  `json:"address"`
    Status       string  `json:"status"`
}
```

**精简原则：**
- ✅ 单文件包含 Model + Request + Response
- ✅ 用注释清晰分隔三个区块
- ✅ 减少文件跳转，提高开发效率
- ✅ 适合微小项目快速迭代
```

### DTO 规范（扁平化设计）

```go
// internal/dto/property_dto.go - 单文件包含所有相关的请求和响应
package dto

// ============ 请求 DTO ============

// ListPropertiesRequest 获取房产列表请求
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

// CreatePropertyRequest 创建房产请求
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

// UpdatePropertyRequest 更新房产请求
type UpdatePropertyRequest struct {
    Title        *string  `json:"title" binding:"omitempty,max=255"`
    Description  *string  `json:"description"`
    Price        *float64 `json:"price" binding:"omitempty,gt=0"`
    Area         *float64 `json:"area" binding:"omitempty,gt=0"`
    Bedrooms     *int     `json:"bedrooms" binding:"omitempty,min=0"`
    Bathrooms    *int     `json:"bathrooms" binding:"omitempty,min=0"`
    PropertyType *string  `json:"property_type"`
    Status       *string  `json:"status" binding:"omitempty,oneof=active inactive sold rented"`
}

// ============ 响应 DTO ============

// PropertyResponse 房产响应（精简版 - 列表用）
type PropertyResponse struct {
    ID           uint     `json:"id"`
    Title        string   `json:"title"`
    Price        float64  `json:"price"`
    Area         float64  `json:"area"`
    Bedrooms     int      `json:"bedrooms"`
    Bathrooms    int      `json:"bathrooms"`
    PropertyType string   `json:"property_type"`
    ListingType  string   `json:"listing_type"`
    Address      string   `json:"address"`
    District     string   `json:"district"`
    CoverImage   string   `json:"cover_image"`
    Status       string   `json:"status"`
    CreatedAt    string   `json:"created_at"`
}

// PropertyDetailResponse 房产详情响应（完整版）
type PropertyDetailResponse struct {
    PropertyResponse                     // 嵌入基础信息
    Description      string              `json:"description"`
    Estate           *EstateInfo         `json:"estate,omitempty"`
    Agent            *AgentInfo          `json:"agent,omitempty"`
    Images           []PropertyImage     `json:"images"`
    Facilities       []FacilityInfo      `json:"facilities"`
    ViewCount        int                 `json:"view_count"`
    FavoriteCount    int                 `json:"favorite_count"`
    UpdatedAt        string              `json:"updated_at"`
}

// PropertyListResponse 房产列表响应（带分页）
type PropertyListResponse struct {
    Items      []*PropertyResponse `json:"items"`
    Pagination PaginationInfo      `json:"pagination"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
    Page      int   `json:"page"`
    PageSize  int   `json:"page_size"`
    Total     int64 `json:"total"`
    TotalPage int   `json:"total_page"`
}

// ============ 辅助结构 ============

type EstateInfo struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

type AgentInfo struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

type PropertyImage struct {
    URL       string `json:"url"`
    IsCover   bool   `json:"is_cover"`
    SortOrder int    `json:"sort_order"`
}

type FacilityInfo struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
    Icon string `json:"icon"`
}
```

**精简原则：**
- ✅ 单文件包含所有 Request 和 Response，减少文件跳转
- ✅ 复用结构体（嵌入 PropertyResponse）
- ✅ 清晰的注释分隔不同区块
- ✅ 使用指针表示可选字段（Update 场景）

### Repository 规范

```go
// databases/property_repo.go
package databases

import (
    "context"
    "gorm.io/gorm"
    "github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

type PropertyRepo struct {
    db *gorm.DB
}

func NewPropertyRepo(db *gorm.DB) *PropertyRepo {
    return &PropertyRepo{db: db}
}

func (r *PropertyRepo) FindAll(ctx context.Context, filter *models.ListPropertiesRequest) ([]*models.Property, error) {
    var properties []*models.Property
    query := r.db.WithContext(ctx).Model(&models.Property{})
    
    // 应用筛选条件
    if filter.DistrictID != nil {
        query = query.Where("district_id = ?", *filter.DistrictID)
    }
    if filter.MinPrice != nil {
        query = query.Where("price >= ?", *filter.MinPrice)
    }
    
    // 分页和排序
    offset := (filter.Page - 1) * filter.PageSize
    query = query.Offset(offset).Limit(filter.PageSize)
    
    // 预加载关联
    query = query.Preload("District").Preload("Estate")
    
    if err := query.Find(&properties).Error; err != nil {
        return nil, err
    }
    
    return properties, nil
}

func (r *PropertyRepo) FindByID(ctx context.Context, id uint) (*models.Property, error) {
    var property models.Property
    if err := r.db.WithContext(ctx).Preload("District").First(&property, id).Error; err != nil {
        return nil, err
    }
    return &property, nil
}

func (r *PropertyRepo) Create(ctx context.Context, property *models.Property) error {
    return r.db.WithContext(ctx).Create(property).Error
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
