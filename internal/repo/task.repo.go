package repo

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"errors"

	"gorm.io/gorm"
)

type ITaskRepository interface {
	FindByID(id uint) (*model.Task, error)
	FindAll() ([]model.Task, error)
	GetListTask(req dto.TaskListRequestDto) ([]model.Task, int64, error)
	FindByUserID(req dto.MyTaskRequestDto, userID uint) ([]model.Task, int64, error)
	Create(task *model.Task) (*model.Task, error)
	Update(task *model.Task) (*model.Task, error)
	Delete(id uint) error
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository() ITaskRepository {
	return &TaskRepository{db: global.Postgres}
}

func (tr *TaskRepository) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	if err := tr.db.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (tr *TaskRepository) FindAll() ([]model.Task, error) {
	var tasks []model.Task
	if err := tr.db.Preload("User").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (tr *TaskRepository) GetListTask(req dto.TaskListRequestDto) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	// Build query with filters
	query := tr.db.Model(&model.Task{})

	// Apply filters
	if req.Client != "" {
		query = query.Where("client = ?", req.Client)
	}
	if req.Job != "" {
		query = query.Where("job = ?", req.Job)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Limit(req.Limit).Offset(req.Skip).Order("created_at DESC")

	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (tr *TaskRepository) FindByUserID(req dto.MyTaskRequestDto, userID uint) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := tr.db.Model(&model.Task{}).Where("user_id = ?", userID)

	if req.Client != "" {
		query = query.Where("client = ?", req.Client)
	}
	if req.Job != "" {
		query = query.Where("job = ?", req.Job)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Limit(req.Limit).Offset(req.Skip).Order("created_at DESC")

	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil

}

func (tr *TaskRepository) Create(task *model.Task) (*model.Task, error) {
	if err := tr.db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *TaskRepository) Update(task *model.Task) (*model.Task, error) {
	if err := tr.db.Model(task).Updates(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *TaskRepository) Delete(id uint) error {
	if err := tr.db.Delete(&model.Task{}, id).Error; err != nil {
		return err
	}
	return nil
}
