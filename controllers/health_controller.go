package controllers

import (
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck 健康检查
func (ctrl *HealthController) HealthCheck(c *gin.Context) {
	tools.Success(c, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Version 版本信息
func (ctrl *HealthController) Version(c *gin.Context) {
	tools.Success(c, gin.H{
		"version":     "1.0.0",
		"app":         "AJO Living API",
		"environment": "development",
		"build_time":  "2026-01-06",
	})
}
