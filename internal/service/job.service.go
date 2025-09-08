package service

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/pkg/response"
)

type IJobService interface {
	GetJobByID(id uint) *response.ServiceResult
	GetListJob(req dto.JobListRequestDto) *response.ServiceResult
	CreateJob(jobDto dto.JobRequestDto) *response.ServiceResult
	UpdateJob(id uint, updateDto dto.JobUpdateRequestDto) *response.ServiceResult
	DeleteJob(id uint) *response.ServiceResult
}

type jobService struct {
	jobRepo repo.IJobRepository
}

func NewJobService(jobRepo repo.IJobRepository) IJobService {
	return &jobService{jobRepo: jobRepo}
}

func (js *jobService) GetJobByID(id uint) *response.ServiceResult {
	result, err := js.jobRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeJobNotFound)
	}

	jobResponse := dto.JobResponseDto{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}
	return response.NewServiceResult(&jobResponse)
}

func (js *jobService) GetListJob(req dto.JobListRequestDto) *response.ServiceResult {
	jobs, total, err := js.jobRepo.GetListJob(req)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	//var jobResponses []dto.JobResponseDto
	//for _, job := range jobs {
	//	jobResponse := dto.JobResponseDto{
	//		ID:        job.ID,
	//		Name:      job.Name,
	//		CreatedAt: job.CreatedAt,
	//		UpdatedAt: job.UpdatedAt,
	//	}
	//	jobResponses = append(jobResponses, jobResponse)
	//}

	result := map[string]interface{}{
		"total": total,
		"data":  jobs,
	}

	return response.NewServiceResult(&result)
}

func (js *jobService) CreateJob(jobDto dto.JobRequestDto) *response.ServiceResult {
	job := &model.Job{
		Name: jobDto.Name,
	}

	createdJob, err := js.jobRepo.Create(job)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	jobResponse := dto.JobResponseDto{
		ID:        createdJob.ID,
		Name:      createdJob.Name,
		CreatedAt: createdJob.CreatedAt,
		UpdatedAt: createdJob.UpdatedAt,
	}

	return response.NewServiceResult(&jobResponse)
}

func (js *jobService) UpdateJob(id uint, updateDto dto.JobUpdateRequestDto) *response.ServiceResult {
	existingJob, err := js.jobRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeJobNotFound)
	}

	if updateDto.Name != "" {
		existingJob.Name = updateDto.Name
	}

	updatedJob, err := js.jobRepo.Update(existingJob)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	jobResponse := dto.JobResponseDto{
		ID:        updatedJob.ID,
		Name:      updatedJob.Name,
		CreatedAt: updatedJob.CreatedAt,
		UpdatedAt: updatedJob.UpdatedAt,
	}

	return response.NewServiceResult(&jobResponse)
}

func (js *jobService) DeleteJob(id uint) *response.ServiceResult {
	_, err := js.jobRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeJobNotFound)
	}

	err = js.jobRepo.Delete(id)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	return response.NewServiceResult(nil)
}
