package dto

import "time"

type JobRequestDto struct {
	Name string `json:"name" binding:"required"`
}

type JobUpdateRequestDto struct {
	Name string `json:"name"`
}

type JobResponseDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JobListRequestDto for pagination and filtering
type JobListRequestDto struct {
	Skip  int    `form:"skip" binding:"min=0"`
	Limit int    `form:"limit" binding:"min=0,max=100"`
	Name  string `form:"name"`
}

// JobListResponseDto for paginated job list response
type JobListResponseDto struct {
	Total int64            `json:"total"`
	Data  []JobResponseDto `json:"data"`
}

