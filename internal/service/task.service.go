package service

import (
	"base_go_be/global"
	"base_go_be/internal/constants"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/pkg/response"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Helper function to convert model.Task to dto.TaskResponseDto
func (ts *TaskService) modelToResponseDto(task *model.Task) *dto.TaskResponseDto {
	return &dto.TaskResponseDto{
		ID:              task.ID,
		UserID:          task.UserID,
		UserInformation: task.UserInformation,
		Client:          task.Client,
		Job:             task.Job,
		Item:            task.Item,
		Role:            task.Role,
		Note:            task.Note,
		OT:              task.OT,
		Volume:          task.Volume,
		StartedAt:       task.StartedAt,
		EndedAt:         task.EndedAt,
		WorkTime:        task.WorkTime,
		Status:          task.Status,
		Delivery:        task.Delivery,
		CreatedAt:       task.CreatedAt,
		CustomCreatedAt: task.CustomCreatedAt,
		UpdatedAt:       task.UpdatedAt,
	}
}

type ITaskService interface {
	GetTaskByID(id uint) *response.ServiceResult
	GetListTask(req dto.TaskListRequestDto, userRole string) *response.ServiceResult
	GetTasksByUserID(req dto.MyTaskRequestDto, userID uint) *response.ServiceResult
	CreateTask(taskRequest *dto.CreateTaskDto, userID uint) *response.ServiceResult
	UpdateTask(id uint, taskRequest *dto.UpdateTaskDto, userID uint) *response.ServiceResult
	UpdateProgressTask(id uint, taskRequest *dto.TaskProcessDto, userID uint) *response.ServiceResult
	DeleteTask(id uint, userID uint) *response.ServiceResult
}

type TaskService struct {
	taskRepo repo.ITaskRepository
}

func NewTaskService(taskRepo repo.ITaskRepository) ITaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (ts *TaskService) GetTaskByID(id uint) *response.ServiceResult {
	task, err := ts.taskRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(404, response.ErrCodeTaskNotFound)
	}

	fmt.Print(task)
	//taskResponse := ts.modelToResponseDto(task)
	return response.NewServiceResult(task)
}

func (ts *TaskService) GetListTask(req dto.TaskListRequestDto, userRole string) *response.ServiceResult {
	// Check authorization - only ADMIN can get full task list without user filter
	if userRole != "ADMIN" {
		return response.NewServiceErrorWithCode(403, response.ErrCodeAccessDenied)
	}

	tasks, total, err := ts.taskRepo.GetListTask(req)
	if err != nil {
		global.Logger.Error("Failed to get tasks from repository: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Convert to DTOs
	//taskResponse := make([]dto.TaskResponseDto, len(tasks))
	//for i, task := range tasks {
	//	taskResponse[i] = *ts.modelToResponseDto(&task)
	//}

	result := map[string]interface{}{
		"total": total,
		"data":  tasks,
	}

	//result := &dto.TaskListResponseDto{
	//	Data:  taskResponse,
	//	Total: total,
	//}
	return response.NewServiceResult(result)
}

func (ts *TaskService) GetTasksByUserID(req dto.MyTaskRequestDto, userID uint) *response.ServiceResult {
	tasks, total, err := ts.taskRepo.FindByUserID(req, userID)
	if err != nil {
		global.Logger.Error("Failed to get user tasks from repository: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	result := map[string]interface{}{
		"total": total,
		"data":  tasks,
	}

	return response.NewServiceResult(result)
}

func (ts *TaskService) CreateTask(taskRequest *dto.CreateTaskDto, userID uint) *response.ServiceResult {
	now := time.Now()
	task := &model.Task{
		UserID:          userID,
		UserInformation: taskRequest.UserInformation,
		Client:          taskRequest.Client,
		Job:             taskRequest.Job,
		Item:            taskRequest.Item,
		Role:            taskRequest.Role,
		Note:            taskRequest.Note,
		OT:              taskRequest.OT,
		Volume:          taskRequest.Volume,
		Status:          constants.TaskStatusOpen,
		CustomCreatedAt: &now,
	}

	createdTask, err := ts.taskRepo.Create(task)
	if err != nil {
		global.Logger.Error("Failed to create task: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	return response.NewServiceResult(createdTask.ID)
}

func (ts *TaskService) UpdateTask(id uint, taskRequest *dto.UpdateTaskDto, userID uint) *response.ServiceResult {
	// First find the existing task
	// Example request body:
	// {
	//   "ot": 234,
	//   "started_at": "2025-09-11T10:30:00Z",
	//   "ended_at": "2025-09-11T11:30:00Z"
	// }
	existingTask, err := ts.taskRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(404, response.ErrCodeTaskNotFound)
	}

	if userID != existingTask.UserID {
		global.Logger.Info("User does not have permission to interact with this task")
		return response.NewServiceErrorWithCode(403, response.ErrCodeTaskPermissionDenied)
	}

	if taskRequest.UserInformation != nil {
		existingTask.UserInformation = taskRequest.UserInformation
	}
	if taskRequest.Client != "" {
		existingTask.Client = taskRequest.Client
	}
	if taskRequest.Job != "" {
		existingTask.Job = taskRequest.Job
	}
	if taskRequest.Item != "" {
		existingTask.Item = taskRequest.Item
	}
	if taskRequest.Role != "" {
		existingTask.Role = taskRequest.Role
	}
	if taskRequest.Note != "" {
		existingTask.Note = taskRequest.Note
	}
	if taskRequest.Status != "" {
		existingTask.Status = taskRequest.Status
	}
	
	if taskRequest.OT != "" {
		existingTask.OT = taskRequest.OT
	}
	if taskRequest.Volume != nil {
		existingTask.Volume = *taskRequest.Volume
	}
	if taskRequest.StartedAt != nil {
		existingTask.StartedAt = taskRequest.StartedAt
	}
	if taskRequest.EndedAt != nil {
		existingTask.EndedAt = taskRequest.EndedAt
	}
	if taskRequest.WorkTime != nil {
		existingTask.WorkTime = *taskRequest.WorkTime
	}
	if taskRequest.Minute != nil {
		existingTask.Minute = *taskRequest.Minute
	}

	if taskRequest.CustomCreatedAt != nil {
		existingTask.CustomCreatedAt = taskRequest.CustomCreatedAt
	}

	updatedTask, err := ts.taskRepo.Update(existingTask)
	if err != nil {
		global.Logger.Error("Failed to update task: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	//taskResponse := ts.modelToResponseDto(updatedTask)
	return response.NewServiceResult(updatedTask)
}

func (ts *TaskService) UpdateProgressTask(id uint, taskRequest *dto.TaskProcessDto, userID uint) *response.ServiceResult {
	existingTask, err := ts.taskRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(404, response.ErrCodeTaskNotFound)
	}
	if userID != existingTask.UserID {
		global.Logger.Info("User does not have permission to interact with this task")
		return response.NewServiceErrorWithCode(403, response.ErrCodeTaskPermissionDenied)
	}

	// Update WorkTime if provided
	if taskRequest.WorkTime != nil {
		existingTask.WorkTime = *taskRequest.WorkTime
	}

	// Update Status if provided
	if taskRequest.Status != "" {
		existingTask.Status = taskRequest.Status
	}

	now := time.Now()
	switch existingTask.Status {
	case constants.TaskStatusInProgress:
		existingTask.StartedAt = &now

	case constants.TaskStatusCompleted:
		existingTask.EndedAt = &now
	}

	updatedTask, err := ts.taskRepo.Update(existingTask)
	if err != nil {
		global.Logger.Error("Failed to update task progress: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	//taskResponse := ts.modelToResponseDto(updatedTask)
	return response.NewServiceResult(updatedTask)
}

func (ts *TaskService) DeleteTask(id uint, userID uint) *response.ServiceResult {
	// Check if task exists
	existingTask, err := ts.taskRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(404, response.ErrCodeTaskNotFound)
	}

	if userID != existingTask.UserID {
		global.Logger.Info("User does not have permission to interact with this task",
			zap.Uint("user_id", userID),
			zap.Uint("task_user_id", existingTask.UserID))
		return response.NewServiceErrorWithCode(403, response.ErrCodeTaskPermissionDenied)
	}

	err = ts.taskRepo.Delete(id)
	if err != nil {
		global.Logger.Error("Failed to delete task: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	return response.NewServiceResult("Task deleted successfully")
}
