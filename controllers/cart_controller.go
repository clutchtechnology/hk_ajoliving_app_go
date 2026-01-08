package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// CartController Methods:
// 0. NewCartController(service *services.CartService) -> 注入 CartService
// 1. GetCart(c *gin.Context) -> 获取购物车
// 2. AddToCart(c *gin.Context) -> 添加到购物车
// 3. UpdateCartItem(c *gin.Context) -> 更新购物车项
// 4. RemoveFromCart(c *gin.Context) -> 移除购物车项
// 5. ClearCart(c *gin.Context) -> 清空购物车

type CartController struct {
	service *services.CartService
}

// 0. NewCartController 构造函数
func NewCartController(service *services.CartService) *CartController {
	return &CartController{service: service}
}

// 1. GetCart 获取购物车
// @Summary 获取购物车
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} tools.Response{data=models.CartResponse}
// @Router /api/v1/cart [get]
func (ctrl *CartController) GetCart(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	cart, err := ctrl.service.GetCart(c.Request.Context(), userID.(uint))
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, cart)
}

// 2. AddToCart 添加到购物车
// @Summary 添加到购物车
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.AddToCartRequest true "添加到购物车请求"
// @Success 201 {object} tools.Response{data=models.CartItemResponse}
// @Router /api/v1/cart/items [post]
func (ctrl *CartController) AddToCart(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	item, err := ctrl.service.AddToCart(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "furniture not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, item)
}

// 3. UpdateCartItem 更新购物车项
// @Summary 更新购物车项
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "购物车项ID"
// @Param body body models.UpdateCartItemRequest true "更新购物车项请求"
// @Success 200 {object} tools.Response{data=models.CartItemResponse}
// @Router /api/v1/cart/items/{id} [put]
func (ctrl *CartController) UpdateCartItem(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	// 解析购物车项ID
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid cart item id")
		return
	}

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	item, err := ctrl.service.UpdateCartItem(c.Request.Context(), userID.(uint), uint(itemID), &req)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "cart item not found")
			return
		}
		if err == tools.ErrForbidden {
			tools.Forbidden(c, "not allowed to update this cart item")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, item)
}

// 4. RemoveFromCart 移除购物车项
// @Summary 移除购物车项
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Param id path int true "购物车项ID"
// @Success 200 {object} tools.Response
// @Router /api/v1/cart/items/{id} [delete]
func (ctrl *CartController) RemoveFromCart(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	// 解析购物车项ID
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid cart item id")
		return
	}

	err = ctrl.service.RemoveFromCart(c.Request.Context(), userID.(uint), uint(itemID))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "cart item not found")
			return
		}
		if err == tools.ErrForbidden {
			tools.Forbidden(c, "not allowed to remove this cart item")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "cart item removed successfully"})
}

// 5. ClearCart 清空购物车
// @Summary 清空购物车
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} tools.Response
// @Router /api/v1/cart [delete]
func (ctrl *CartController) ClearCart(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	err := ctrl.service.ClearCart(c.Request.Context(), userID.(uint))
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "cart cleared successfully"})
}
