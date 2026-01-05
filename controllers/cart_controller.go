package controllers

// CartHandler Methods:
// 0. NewCartHandler(cartService *service.CartService) -> 注入 CartService
// 1. GetCart(c *gin.Context) -> 获取购物车
// 2. AddToCart(c *gin.Context) -> 添加到购物车
// 3. UpdateCartItem(c *gin.Context) -> 更新购物车项
// 4. RemoveFromCart(c *gin.Context) -> 移除购物车项
// 5. ClearCart(c *gin.Context) -> 清空购物车

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
)

// CartHandlerInterface 购物车处理器接口
type CartHandlerInterface interface {
	GetCart(c *gin.Context)        // 1. 获取购物车
	AddToCart(c *gin.Context)      // 2. 添加到购物车
	UpdateCartItem(c *gin.Context) // 3. 更新购物车项
	RemoveFromCart(c *gin.Context) // 4. 移除购物车项
	ClearCart(c *gin.Context)      // 5. 清空购物车
}

// CartHandler 购物车处理器
type CartHandler struct {
	cartService *service.CartService
}

// 0. NewCartHandler 注入 CartService
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// 1. GetCart 获取购物车
// GetCart godoc
// @Summary      获取购物车
// @Description  获取当前用户的购物车内容
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=response.CartResponse}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID.(uint))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, cart)
}

// 2. AddToCart 添加到购物车
// AddToCart godoc
// @Summary      添加到购物车
// @Description  添加家具商品到购物车
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      request.AddToCartRequest  true  "购物车项信息"
// @Success      201   {object}  response.Response{data=response.AddToCartResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /api/v1/cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req request.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	result, err := h.cartService.AddToCart(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		if err == service.ErrFurnitureNotFound {
			response.NotFound(c, "Furniture not found")
			return
		}
		if err == service.ErrFurnitureNotAvailable {
			response.BadRequest(c, "Furniture is not available")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// 3. UpdateCartItem 更新购物车项
// UpdateCartItem godoc
// @Summary      更新购物车项
// @Description  更新购物车中的商品数量
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                              true  "购物车项ID"
// @Param        body  body      request.UpdateCartItemRequest  true  "更新信息"
// @Success      200   {object}  response.Response{data=response.UpdateCartItemResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /api/v1/cart/items/{id} [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cart item ID")
		return
	}

	var req request.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	result, err := h.cartService.UpdateCartItem(c.Request.Context(), userID.(uint), uint(id), &req)
	if err != nil {
		if err == service.ErrCartItemNotFound {
			response.NotFound(c, "Cart item not found")
			return
		}
		if err == service.ErrNotCartItemOwner {
			response.Forbidden(c, "You are not the owner of this cart item")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// 4. RemoveFromCart 移除购物车项
// RemoveFromCart godoc
// @Summary      移除购物车项
// @Description  从购物车中移除指定商品
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "购物车项ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/cart/items/{id} [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid cart item ID")
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	err = h.cartService.RemoveFromCart(c.Request.Context(), userID.(uint), uint(id))
	if err != nil {
		if err == service.ErrCartItemNotFound {
			response.NotFound(c, "Cart item not found")
			return
		}
		if err == service.ErrNotCartItemOwner {
			response.Forbidden(c, "You are not the owner of this cart item")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Cart item removed successfully"})
}

// 5. ClearCart 清空购物车
// ClearCart godoc
// @Summary      清空购物车
// @Description  清空当前用户的购物车
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	err := h.cartService.ClearCart(c.Request.Context(), userID.(uint))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Cart cleared successfully"})
}
