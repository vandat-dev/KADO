package repo

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"errors"

	"gorm.io/gorm"
)

type IJobRepository interface {
	FindByID(id uint) (*model.Job, error)
	FindAll() ([]model.Job, error)
	GetListJob(req dto.JobListRequestDto) ([]model.Job, int64, error)
	Create(job *model.Job) (*model.Job, error)
	Update(job *model.Job) (*model.Job, error)
	Delete(id uint) error
}

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository() IJobRepository {
	return &JobRepository{db: global.Postgres}
}

func (jr *JobRepository) FindByID(id uint) (*model.Job, error) {
	var job model.Job
	if err := jr.db.First(&job, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}
	return &job, nil
}

func (jr *JobRepository) FindAll() ([]model.Job, error) {
	var jobs []model.Job
	if err := jr.db.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (jr *JobRepository) GetListJob(req dto.JobListRequestDto) ([]model.Job, int64, error) {
	var jobs []model.Job
	var total int64

	// Build query with filters
	query := jr.db.Model(&model.Job{})

	// Apply filters
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	} else {
		query = query.Limit(10)
	}

	query = query.Offset(req.Skip).Order("created_at DESC")

	if err := query.Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (jr *JobRepository) Create(job *model.Job) (*model.Job, error) {
	if err := jr.db.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (jr *JobRepository) Update(job *model.Job) (*model.Job, error) {
	if err := jr.db.Save(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (jr *JobRepository) Delete(id uint) error {
	if err := jr.db.Delete(&model.Job{}, id).Error; err != nil {
		return err
	}
	return nil
}

