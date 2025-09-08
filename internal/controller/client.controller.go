package controller

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/service"
	"base_go_be/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientController struct {
	clientService service.IClientService
}

func NewClientController(clientService service.IClientService) *ClientController {
	return &ClientController{
		clientService: clientService,
	}
}

// GetClientByID godoc
// @Summary Get client by ID
// @Description Get client details by ID
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Client ID"
// @Success 200 {object} response.Response{data=dto.ClientResponseDto} "Client details"
// @Failure 422 {object} response.Response "Client not found"
// @Router /client/detail/{id} [get]
func (cc *ClientController) GetClientByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := cc.clientService.GetClientByID(uint(id))
	response.HandleServiceResult(c, result)
}

// GetListClient godoc
// @Summary Get list of clients
// @Description Get paginated list of clients with optional filtering
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param skip query int false "Number of records to skip" default(0)
// @Param limit query int false "Number of records to return" default(10)
// @Param name query string false "Filter by client name"
// @Success 200 {object} response.Response{data=dto.ClientListResponseDto} "List of clients"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /client/list [get]
func (cc *ClientController) GetListClient(c *gin.Context) {
	var req dto.ClientListRequestDto
	if err := c.ShouldBindQuery(&req); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	result := cc.clientService.GetListClient(req)
	response.HandleServiceResult(c, result)
}

// CreateClient godoc
// @Summary Create a new client
// @Description Create a new client
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param client body dto.ClientRequestDto true "Client data"
// @Success 200 {object} response.Response{data=dto.ClientResponseDto} "Client created successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /client/create [post]
func (cc *ClientController) CreateClient(c *gin.Context) {
	var clientRequest dto.ClientRequestDto
	if err := c.ShouldBindJSON(&clientRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := cc.clientService.CreateClient(clientRequest)
	response.HandleServiceResult(c, result)
}

// UpdateClient godoc
// @Summary Update client
// @Description Update client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Client ID"
// @Param client body dto.ClientUpdateRequestDto true "Client update data"
// @Success 200 {object} response.Response{data=dto.ClientResponseDto} "Client updated successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 422 {object} response.Response "Client not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /client/update/{id} [put]
func (cc *ClientController) UpdateClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	var updateRequest dto.ClientUpdateRequestDto
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := cc.clientService.UpdateClient(uint(id), updateRequest)
	response.HandleServiceResult(c, result)
}

// DeleteClient godoc
// @Summary Delete client
// @Description Delete client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Client ID"
// @Success 200 {object} response.Response "Client deleted successfully"
// @Failure 422 {object} response.Response "Client not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /client/delete/{id} [delete]
func (cc *ClientController) DeleteClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := cc.clientService.DeleteClient(uint(id))
	response.HandleServiceResult(c, result)
}
