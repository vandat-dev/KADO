package controller

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/service"
	"base_go_be/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService service.IRoleService
}

func NewRoleController(roleService service.IRoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// GetRoleByID godoc
// @Summary Get role by ID
// @Description Get role details by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response{data=dto.RoleResponseDto} "Role details"
// @Failure 422 {object} response.Response "Role not found"
// @Router /role/detail/{id} [get]
func (rc *RoleController) GetRoleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := rc.roleService.GetRoleByID(uint(id))
	response.HandleServiceResult(c, result)
}

// GetListRole godoc
// @Summary Get list of roles
// @Description Get paginated list of roles with optional filtering
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param skip query int false "Number of records to skip" default(0)
// @Param limit query int false "Number of records to return" default(10)
// @Param name query string false "Filter by role name"
// @Success 200 {object} response.Response{data=dto.RoleListResponseDto} "List of roles"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /role/list [get]
func (rc *RoleController) GetListRole(c *gin.Context) {
	var req dto.RoleListRequestDto
	if err := c.ShouldBindQuery(&req); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	result := rc.roleService.GetListRole(req)
	response.HandleServiceResult(c, result)
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param role body dto.RoleRequestDto true "Role data"
// @Success 200 {object} response.Response{data=dto.RoleResponseDto} "Role created successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /role/create [post]
func (rc *RoleController) CreateRole(c *gin.Context) {
	var roleRequest dto.RoleRequestDto
	if err := c.ShouldBindJSON(&roleRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := rc.roleService.CreateRole(roleRequest)
	response.HandleServiceResult(c, result)
}

// UpdateRole godoc
// @Summary Update role
// @Description Update role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Role ID"
// @Param role body dto.RoleUpdateRequestDto true "Role update data"
// @Success 200 {object} response.Response{data=dto.RoleResponseDto} "Role updated successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 422 {object} response.Response "Role not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /role/update/{id} [put]
func (rc *RoleController) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	var updateRequest dto.RoleUpdateRequestDto
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := rc.roleService.UpdateRole(uint(id), updateRequest)
	response.HandleServiceResult(c, result)
}

// DeleteRole godoc
// @Summary Delete role
// @Description Delete role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response "Role deleted successfully"
// @Failure 422 {object} response.Response "Role not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /role/delete/{id} [delete]
func (rc *RoleController) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := rc.roleService.DeleteRole(uint(id))
	response.HandleServiceResult(c, result)
}
