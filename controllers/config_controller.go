package controllers

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ConfigHandler 配置处理器
type ConfigHandler struct {
	*BaseHandler
	service service.ConfigService
}

// NewConfigHandler 创建配置处理器实例
func NewConfigHandler(baseHandler *BaseHandler, service service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		BaseHandler: baseHandler,
		service:     service,
	}
}

// GetConfig 获取系统配置
// @Summary 获取系统配置
// @Description 获取系统配置信息，包括系统信息、应用配置、功能开关、API配置、UI配置等
// @Tags 系统配置
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.ConfigResponse} "配置信息"
// @Failure 500 {object} models.Response "服务器错误"
// @Router /api/v1/config [get]
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	data, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		h.Logger.Error("获取系统配置失败", zap.Error(err))
		models.InternalError(c, "获取配置失败")
		return
	}

	models.Success(c, data)
}

// GetRegions 获取区域配置
// @Summary 获取区域配置
// @Description 获取香港三大区域（香港岛、九龙、新界）及其下属地区的配置信息
// @Tags 系统配置
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.RegionsResponse} "区域配置"
// @Failure 500 {object} models.Response "服务器错误"
// @Router /api/v1/config/regions [get]
func (h *ConfigHandler) GetRegions(c *gin.Context) {
	data, err := h.service.GetRegions(c.Request.Context())
	if err != nil {
		h.Logger.Error("获取区域配置失败", zap.Error(err))
		models.InternalError(c, "获取区域配置失败")
		return
	}

	models.Success(c, data)
}

// GetPropertyTypes 获取房产类型配置
// @Summary 获取房产类型配置
// @Description 获取房产类型、房源类型、状态等配置信息
// @Tags 系统配置
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.PropertyTypesResponse} "房产类型配置"
// @Failure 500 {object} models.Response "服务器错误"
// @Router /api/v1/config/property-types [get]
func (h *ConfigHandler) GetPropertyTypes(c *gin.Context) {
	data, err := h.service.GetPropertyTypes(c.Request.Context())
	if err != nil {
		h.Logger.Error("获取房产类型配置失败", zap.Error(err))
		models.InternalError(c, "获取房产类型配置失败")
		return
	}

	models.Success(c, data)
}
