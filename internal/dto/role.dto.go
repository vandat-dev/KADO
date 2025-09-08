package dto

import "time"

type RoleRequestDto struct {
	Name string `json:"name" binding:"required"`
}

type RoleUpdateRequestDto struct {
	Name string `json:"name"`
}

type RoleResponseDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoleListRequestDto for pagination and filtering
type RoleListRequestDto struct {
	Skip  int    `form:"skip" binding:"min=0"`
	Limit int    `form:"limit" binding:"min=0,max=100"`
	Name  string `form:"name"`
}

// RoleListResponseDto for paginated role list response
type RoleListResponseDto struct {
	Total int64             `json:"total"`
	Data  []RoleResponseDto `json:"data"`
}
