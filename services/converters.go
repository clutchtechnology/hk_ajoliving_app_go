package services

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// convertPropertyToListItemResponse 通用的房产转换函数（直接返回，预加载了关联数据）
func convertPropertyToListItemResponse(property *models.Property) *models.Property {
	if property == nil {
		return nil
	}
	
	// 直接返回 property，GORM 已经预加载了所有关联数据
	return property
}
