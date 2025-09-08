package dto

import (
	"base_go_be/internal/model"
	"time"
)

type TaskRequestDto struct {
	UserID          uint                   `json:"user_id"`
	UserInformation *model.UserInformation `json:"user_information"`
	Client          string                 `json:"client" binding:"required"`
	Job             string                 `json:"job" binding:"required"`
	Item            string                 `json:"item"`
	Role            string                 `json:"role"`
	Note            string                 `json:"note"`
	OT              *int                   `json:"ot"`
	Volume          *int                   `json:"volume"`
	StartedAt       *time.Time             `json:"started_at"`
	EndedAt         *time.Time             `json:"ended_at"`
	WorkTime        int                    `json:"work_time"`
	Status          string                 `json:"status"`
	Delivery        bool                   `json:"delivery"`
}

type CreateTaskDto struct {
	UserInformation *model.UserInformation `json:"user_information,omitempty"`
	Client          string                 `json:"client,omitempty"`
	Job             string                 `json:"job,omitempty"`
	Item            string                 `json:"item,omitempty"`
	Role            string                 `json:"role,omitempty"`
	Note            string                 `json:"note,omitempty"`
	OT              int                    `json:"ot,omitempty"`
	Volume          int                    `json:"volume,omitempty"`
}

type UpdateTaskDto struct {
	UserInformation *model.UserInformation `json:"user_information"`
	Client          string                 `json:"client"`
	Job             string                 `json:"job"`
	Item            string                 `json:"item"`
	Role            string                 `json:"role"`
	Note            string                 `json:"note"`
	OT              *int                   `json:"ot"`
	Volume          *int                   `json:"volume"`
	StartedAt       *time.Time             `json:"started_at"`
	EndedAt         *time.Time             `json:"ended_at"`
	WorkTime        *int                   `json:"work_time"`
	Minute          *int                   `json:"minute"`
	Status          string                 `json:"status" binding:"omitempty,oneof=OPEN IN_PROGRESS PENDING COMPLETED"`
	CustomCreatedAt *time.Time             `json:"custom_created_at"`
}

type TaskResponseDto struct {
	ID              uint                   `json:"id"`
	UserID          uint                   `json:"user_id"`
	UserInformation *model.UserInformation `json:"user_information"`
	Client          string                 `json:"client"`
	Job             string                 `json:"job"`
	Item            string                 `json:"item"`
	Role            string                 `json:"role"`
	Note            string                 `json:"note"`
	OT              int                    `json:"ot"`
	Volume          int                    `json:"volume"`
	StartedAt       *time.Time             `json:"started_at"`
	EndedAt         *time.Time             `json:"ended_at"`
	WorkTime        int                    `json:"work_time"`
	Status          string                 `json:"status"`
	Delivery        bool                   `json:"delivery"`
	CreatedAt       time.Time              `json:"created_at"`
	CustomCreatedAt *time.Time             `json:"custom_created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// TaskListRequestDto for pagination and filtering
type TaskListRequestDto struct {
	Skip   int    `form:"skip" binding:"min=0"`
	Limit  int    `form:"limit" binding:"min=0,max=100"`
	Client string `form:"client"`
	Job    string `form:"job"`
	Status string `form:"status" binding:"omitempty,oneof=OPEN IN_PROGRESS PENDING COMPLETED"`
	UserID uint   `form:"user_id"`
}

type MyTaskRequestDto struct {
	Skip   int    `form:"skip" binding:"min=0"`
	Limit  int    `form:"limit" binding:"min=0,max=100"`
	Client string `form:"client"`
	Job    string `form:"job"`
	Status string `form:"status" binding:"omitempty,oneof=OPEN IN_PROGRESS PENDING COMPLETED"`
}

type TaskProcessDto struct {
	WorkTime *int   `json:"work_time"`
	Status   string `json:"status" binding:"omitempty,oneof=OPEN IN_PROGRESS PENDING COMPLETED"`
}

// TaskListResponseDto for paginated task list response
type TaskListResponseDto struct {
	Total int64             `json:"total"`
	Data  []TaskResponseDto `json:"data"`
}
