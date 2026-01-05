package services

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// convertPropertyToListItemResponse 通用的房产转换函数
func convertPropertyToListItemResponse(property *models.Property) *models.Property {
	if property == nil {
		return nil
	}

	resp := &models.Property{
		ID:            property.ID,
		PropertyNo:    property.PropertyNo,
		Title:         property.Title,
		Price:         property.Price,
		Area:          property.Area,
		Address:       property.Address,
		DistrictID:    property.DistrictID,
		Status:        string(property.Status),
		ViewCount:     property.ViewCount,
		FavoriteCount: property.FavoriteCount,
		CreatedAt:     property.CreatedAt,
		Bedrooms:      property.Bedrooms,
	}

	// 处理可选的指针类型字段
	if property.BuildingName != nil {
		resp.BuildingName = *property.BuildingName
	}
	
	if property.Bathrooms != nil {
		resp.Bathrooms = *property.Bathrooms
	}

	// 处理枚举类型
	resp.PropertyType = string(property.PropertyType)
	resp.ListingType = string(property.ListingType)

	// 处理关联数据 - District
	if property.District != nil {
		resp.District = &models.DistrictResponse{
			ID:         property.District.ID,
			NameZhHant: property.District.NameZhHant,
			Region:     string(property.District.Region),
		}
		if property.District.NameZhHans != nil {
			resp.District.NameZhHans = *property.District.NameZhHans
		}
		if property.District.NameEn != nil {
			resp.District.NameEn = *property.District.NameEn
		}
	}

	// 设置封面图片
	if len(property.Images) > 0 {
		for _, img := range property.Images {
			if img.IsCover {
				resp.CoverImage = img.URL
				break
			}
		}
		// 如果没有封面图，使用第一张图
		if resp.CoverImage == "" && len(property.Images) > 0 {
			resp.CoverImage = property.Images[0].URL
		}
	}

	return resp
}
