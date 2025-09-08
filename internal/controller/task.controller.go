package controller

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/service"
	"base_go_be/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService service.ITaskService
}

func NewTaskController(taskService service.ITaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Get task details by ID
// @Tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Success 200 {object} response.Response{data=dto.TaskResponseDto}
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Task not found"
// @Failure 422 {object} response.Response "Invalid task ID"
// @Router /task/detail/{id} [get]
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	idParam := c.Param("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 0)
	id := uint(idUint64)
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}

	result := tc.taskService.GetTaskByID(id)
	response.HandleServiceResult(c, result)
}

// GetListTask godoc
// @Summary Get list of tasks with pagination and filtering
// @Description Get paginated list of tasks with filtering options. Only admin users can access this endpoint without user filter.
// @Tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(10)
// @Param client query string false "Client filter"
// @Param job query string false "Job filter"
// @Param status query string false "Status filter" Enums(OPEN, IN_PROGRESS, PENDING, COMPLETED)
// @Param user_id query int false "User ID filter"
// @Success 200 {object} response.Response{data=dto.TaskListResponseDto} "Paginated list of tasks"
// @Failure 400 {object} response.Response "Invalid query parameters"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Access denied: Only admin can view all tasks"
// @Router /task/list [get]
func (tc *TaskController) GetListTask(c *gin.Context) {
	var req dto.TaskListRequestDto

	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Error("Failed to bind query parameters: " + err.Error())
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}

	role, _ := c.Get("system_role")
	result := tc.taskService.GetListTask(req, role.(string))
	response.HandleServiceResult(c, result)
}

// GetMyTasks godoc
// @Summary Get current user's tasks
// @Description Get tasks belonging to the authenticated user
// @Tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(10)
// @Param client query string false "Client filter"
// @Param job query string false "Job filter"
// @Param status query string false "Status filter" Enums(OPEN, IN_PROGRESS, PENDING, COMPLETED)
// @Failure 401 {object} response.Response "Unauthorized"
// @Success 200 {object} response.Response{data=dto.MyTaskRequestDto} "Paginated list of my tasks"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /task/my-tasks [get]
func (tc *TaskController) GetMyTasks(c *gin.Context) {

	var req dto.MyTaskRequestDto

	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Error("Failed to bind query parameters: " + err.Error())
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}

	// Get user ID from JWT token context
	userID, _ := c.Get("user_id")

	result := tc.taskService.GetTasksByUserID(req, userID.(uint))
	response.HandleServiceResult(c, result)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with the provided information
// @Tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param task body dto.CreateTaskDto true "Create Task Request"
// @Success 200 {object} response.Response{data=map[string]uint}
// @Failure 400 {object} response.Response "Invalid request payload"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 405 {object} response.Response "Method not allowed"
// @Router /task/create [post]
func (tc *TaskController) CreateTask(c *gin.Context) {
	taskRequest := dto.CreateTaskDto{}

	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		response.ErrorResponse(c, 400, "Invalid request payload")
		return
	}

	// Get user ID from JWT token context
	userID, _ := c.Get("user_id")

	result := tc.taskService.CreateTask(&taskRequest, userID.(uint))
	response.HandleServiceResult(c, result)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update task information by ID
// @Tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Param task body dto.UpdateTaskDto true "Update Task Request"
// @Success 200 {object} response.Response{data=dto.TaskResponseDto}
// @Failure 400 {object} response.Response "Invalid request payload"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Task not found"
// @Failure 422 {object} response.Response "Invalid task ID"
// @Router /task/update/{id} [put]
func (tc *TaskController) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	userID, _ := c.Get("user_id")
	idUint64, err := strconv.ParseUint(idParam, 10, 0)
	id := uint(idUint64)
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}

	taskRequest := dto.UpdateTaskDto{}
	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		response.ErrorResponse(c, 400, "Invalid request payload")
		return
	}

	result := tc.taskService.UpdateTask(id, &taskRequest, userID.(uint))
	response.HandleServiceResult(c, result)
}

// UpdateProgressTask godoc
// @Summary Update task progress
// @Description Update task progress (work time and status) by ID
// @Tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Param task body dto.TaskProcessDto true "Update Task Progress Request"
// @Success 200 {object} response.Response{data=dto.TaskResponseDto}
// @Failure 400 {object} response.Response "Invalid request payload"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Task not found"
// @Failure 422 {object} response.Response "Invalid task ID"
// @Router /task/update_progress/{id} [put]
func (tc *TaskController) UpdateProgressTask(c *gin.Context) {
	idParam := c.Param("id")
	userID, _ := c.Get("user_id")
	idUint64, err := strconv.ParseUint(idParam, 10, 0)
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}
	id := uint(idUint64)

	var taskRequest dto.TaskProcessDto
	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		response.ErrorResponse(c, 400, "Invalid request payload")
		return
	}
	result := tc.taskService.UpdateProgressTask(id, &taskRequest, userID.(uint))
	response.HandleServiceResult(c, result)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete task by ID
// @Tags task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Success 200 {object} response.Response{data=string}
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Task not found"
// @Failure 422 {object} response.Response "Invalid task ID"
// @Router /task/delete/{id} [delete]
func (tc *TaskController) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	userID, _ := c.Get("user_id")
	idUint64, err := strconv.ParseUint(idParam, 10, 0)
	id := uint(idUint64)
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}

	result := tc.taskService.DeleteTask(id, userID.(uint))
	response.HandleServiceResult(c, result)
}
