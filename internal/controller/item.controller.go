package controller

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/service"
	"base_go_be/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
	itemService service.IItemService
}

func NewItemController(itemService service.IItemService) *ItemController {
	return &ItemController{
		itemService: itemService,
	}
}

// GetItemByID godoc
// @Summary Get item by ID
// @Description Get item details by ID
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Item ID"
// @Success 200 {object} response.Response{data=dto.ItemResponseDto} "Item details"
// @Failure 422 {object} response.Response "Item not found"
// @Router /item/detail/{id} [get]
func (ic *ItemController) GetItemByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := ic.itemService.GetItemByID(uint(id))
	response.HandleServiceResult(c, result)
}

// GetListItem godoc
// @Summary Get list of items
// @Description Get paginated list of items with optional filtering
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param skip query int false "Number of records to skip" default(0)
// @Param limit query int false "Number of records to return" default(10)
// @Param name query string false "Filter by item name"
// @Param category query string false "Filter by item category"
// @Success 200 {object} response.Response{data=dto.ItemListResponseDto} "List of items"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /item/list [get]
func (ic *ItemController) GetListItem(c *gin.Context) {
	var req dto.ItemListRequestDto
	if err := c.ShouldBindQuery(&req); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	result := ic.itemService.GetListItem(req)
	response.HandleServiceResult(c, result)
}

// CreateItem godoc
// @Summary Create a new item
// @Description Create a new item
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param item body dto.ItemRequestDto true "Item data"
// @Success 200 {object} response.Response{data=dto.ItemResponseDto} "Item created successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /item/create [post]
func (ic *ItemController) CreateItem(c *gin.Context) {
	var itemRequest dto.ItemRequestDto
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := ic.itemService.CreateItem(itemRequest)
	response.HandleServiceResult(c, result)
}

// UpdateItem godoc
// @Summary Update item
// @Description Update item by ID
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Item ID"
// @Param item body dto.ItemUpdateRequestDto true "Item update data"
// @Success 200 {object} response.Response{data=dto.ItemResponseDto} "Item updated successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 422 {object} response.Response "Item not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /item/update/{id} [put]
func (ic *ItemController) UpdateItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	var updateRequest dto.ItemUpdateRequestDto
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := ic.itemService.UpdateItem(uint(id), updateRequest)
	response.HandleServiceResult(c, result)
}

// DeleteItem godoc
// @Summary Delete item
// @Description Delete item by ID
// @Tags items
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Item ID"
// @Success 200 {object} response.Response "Item deleted successfully"
// @Failure 422 {object} response.Response "Item not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /item/delete/{id} [delete]
func (ic *ItemController) DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := ic.itemService.DeleteItem(uint(id))
	response.HandleServiceResult(c, result)
}
