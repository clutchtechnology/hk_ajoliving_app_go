package controllers

// BaseHandler Methods:
// 0. NewBaseHandler() -> 创建基础处理器
// 1. HealthCheck(c *gin.Context) -> 健康检查
// 2. Version(c *gin.Context) -> 版本信息

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
)

// BaseHandlerInterface 基础处理器接口
type BaseHandlerInterface interface {
	HealthCheck(c *gin.Context) // 1. 健康检查
	Version(c *gin.Context)     // 2. 版本信息
}

// BaseHandler 基础路由处理器
type BaseHandler struct{}

// 0. NewBaseHandler 创建基础处理器
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// 1. HealthCheck 健康检查
// HealthCheck godoc
// @Summary      健康检查
// @Description  检查服务是否正常运行
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=map[string]interface{}}
// @Router       /api/v1/health [get]
func (h *BaseHandler) HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "AJO Living API",
	})
}

// 2. Version 版本信息
// Version godoc
// @Summary      版本信息
// @Description  获取 API 版本信息
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=map[string]interface{}}
// @Router       /api/v1/version [get]
func (h *BaseHandler) Version(c *gin.Context) {
	response.Success(c, gin.H{
		"version":     "1.0.0",
		"api_version": "v1",
		"build_time":  "2025-12-18",
		"go_version":  "1.21+",
	})
}
