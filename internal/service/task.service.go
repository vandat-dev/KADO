package service

import (
	"base_go_be/global"
	"base_go_be/internal/constants"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/internal/until"
	"base_go_be/pkg/response"
	"bytes"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ITaskService interface {
	GetTaskByID(id uint) *response.ServiceResult
	GetListTask(req dto.TaskListRequestDto, userRole string) *response.ServiceResult
	GetStatisticTask(req dto.TaskStatisticRequestDto) *response.ServiceResult
	GetTasksByUserID(req dto.MyTaskRequestDto, userID uint) *response.ServiceResult
	CreateTask(taskRequest *dto.CreateTaskDto, userID uint) *response.ServiceResult
	ExportTasks(req dto.TaskStatisticRequestDto) *response.ServiceResult
	UpdateTask(id uint, taskRequest *dto.UpdateTaskDto, userID uint, userSystemRole string) *response.ServiceResult
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

	result := map[string]interface{}{
		"total": total,
		"data":  tasks,
	}
	return response.NewServiceResult(result)
}

func (ts *TaskService) GetStatisticTask(req dto.TaskStatisticRequestDto) *response.ServiceResult {
	tasks, err := ts.taskRepo.GetListStatisticTask(req)
	if err != nil {
		global.Logger.Error("Failed to get tasks from repository: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	taskResponse, errHandle := ts.ResponseDataTaskExport(req, tasks)
	if errHandle != nil {
		global.Logger.Error("Invalid export type: " + errHandle.Error())
		return response.NewServiceErrorWithCode(400, response.ErrCodeInvalidData)
	}

	return response.NewServiceResult(taskResponse)
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

func (ts *TaskService) ExportTasks(req dto.TaskStatisticRequestDto) *response.ServiceResult {

	tasks, err := ts.taskRepo.GetListStatisticTask(req)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}
	resp, err := ts.ResponseDataTaskExport(req, tasks)
	if err != nil {
		return response.NewServiceErrorWithCode(400, response.ErrCodeInvalidData)
	}

	// --- Handle Excel ---
	buf, err := ts.ProcessExportTasks(req, resp)
	if err != nil {
		return response.NewServiceErrorWithCode(400, response.ErrCodeTaskExportFailed)
	}

	data := map[string]interface{}{
		"filename": "timesheet.xlsx",
		"content":  buf.Bytes(),
	}
	return response.NewServiceResult(data)
}

func (ts *TaskService) ResponseDataTaskExport(req dto.TaskStatisticRequestDto, tasks []model.Task) ([]map[string]interface{}, error) {
	if req.StartDate == nil || req.EndDate == nil {
		return nil, fmt.Errorf("start_date and end_date are required")
	}
	switch req.TypeExport {
	case constants.TaskExportSummary:
		taskSummaries := ts.summaryTask(tasks)
		return taskSummaries, nil

	case constants.TaskExportProjectTotals:
		taskSummaries := ts.analysisTasksByType(tasks, req.StartDate, req.EndDate, constants.TaskExportProjectTotals)
		return taskSummaries, nil

	case constants.TaskExportTimeTotals:
		taskSummaries := ts.analysisTasksByType(tasks, req.StartDate, req.EndDate, constants.TaskExportTimeTotals)
		return taskSummaries, nil

	default:
		return nil, fmt.Errorf("invalid TypeExport value: %s", req.TypeExport)
	}

}

func (ts *TaskService) ProcessExportTasks(req dto.TaskStatisticRequestDto, resp []map[string]interface{}) (*bytes.Buffer, error) {
	switch req.TypeExport {
	case constants.TaskExportSummary:
		headers := until.GetHeaderSummary()
		buf, err := exportExcelSummary(resp, headers)
		return buf, err

	case constants.TaskExportProjectTotals:
		headers := until.GetDateRange(req.StartDate, req.EndDate)
		buf, err := exportExcelJob(resp, headers)
		return buf, err

	case constants.TaskExportTimeTotals:
		headers := until.GetDateRange(req.StartDate, req.EndDate)
		buf, err := exportExcelTimeTotal(resp, headers)
		return buf, err

	default:
		return nil, fmt.Errorf("failed export process: %s", req.TypeExport)
	}
}

func (ts *TaskService) UpdateTask(id uint, taskRequest *dto.UpdateTaskDto, userID uint, userSystemRole string) *response.ServiceResult {
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

	if userSystemRole != constants.Admin && userID != existingTask.UserID {
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
	if taskRequest.Hours != "" {
		existingTask.Hours = taskRequest.Hours
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
	if taskRequest.BreakTime != nil {
		existingTask.BreakTime = *taskRequest.BreakTime
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
