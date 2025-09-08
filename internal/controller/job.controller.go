package controller

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/service"
	"base_go_be/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	jobService service.IJobService
}

func NewJobController(jobService service.IJobService) *JobController {
	return &JobController{
		jobService: jobService,
	}
}

// GetJobByID godoc
// @Summary Get job by ID
// @Description Get job details by ID
// @Tags jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Job ID"
// @Success 200 {object} response.Response{data=dto.JobResponseDto} "Job details"
// @Failure 422 {object} response.Response "Job not found"
// @Router /job/detail/{id} [get]
func (jc *JobController) GetJobByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := jc.jobService.GetJobByID(uint(id))
	response.HandleServiceResult(c, result)
}

// GetListJob godoc
// @Summary Get list of jobs
// @Description Get paginated list of jobs with optional filtering
// @Tags jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param skip query int false "Number of records to skip" default(0)
// @Param limit query int false "Number of records to return" default(10)
// @Param name query string false "Filter by job name"
// @Success 200 {object} response.Response{data=dto.JobListResponseDto} "List of jobs"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /job/list [get]
func (jc *JobController) GetListJob(c *gin.Context) {
	var req dto.JobListRequestDto
	if err := c.ShouldBindQuery(&req); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	result := jc.jobService.GetListJob(req)
	response.HandleServiceResult(c, result)
}

// CreateJob godoc
// @Summary Create a new job
// @Description Create a new job
// @Tags jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param job body dto.JobRequestDto true "Job data"
// @Success 200 {object} response.Response{data=dto.JobResponseDto} "Job created successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /job/create [post]
func (jc *JobController) CreateJob(c *gin.Context) {
	var jobRequest dto.JobRequestDto
	if err := c.ShouldBindJSON(&jobRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := jc.jobService.CreateJob(jobRequest)
	response.HandleServiceResult(c, result)
}

// UpdateJob godoc
// @Summary Update job
// @Description Update job by ID
// @Tags jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Job ID"
// @Param job body dto.JobUpdateRequestDto true "Job update data"
// @Success 200 {object} response.Response{data=dto.JobResponseDto} "Job updated successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 422 {object} response.Response "Job not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /job/update/{id} [put]
func (jc *JobController) UpdateJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	var updateRequest dto.JobUpdateRequestDto
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := jc.jobService.UpdateJob(uint(id), updateRequest)
	response.HandleServiceResult(c, result)
}

// DeleteJob godoc
// @Summary Delete job
// @Description Delete job by ID
// @Tags jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Job ID"
// @Success 200 {object} response.Response "Job deleted successfully"
// @Failure 422 {object} response.Response "Job not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /job/delete/{id} [delete]
func (jc *JobController) DeleteJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := jc.jobService.DeleteJob(uint(id))
	response.HandleServiceResult(c, result)
}
