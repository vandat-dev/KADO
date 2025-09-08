package dto

import "time"

type ClientRequestDto struct {
	Name string `json:"name" binding:"required"`
}

type ClientUpdateRequestDto struct {
	Name string `json:"name"`
}

type ClientResponseDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ClientListRequestDto for pagination and filtering
type ClientListRequestDto struct {
	Skip  int    `form:"skip" binding:"min=0"`
	Limit int    `form:"limit" binding:"min=0,max=100"`
	Name  string `form:"name"`
}

// ClientListResponseDto for paginated client list response
type ClientListResponseDto struct {
	Total int64               `json:"total"`
	Data  []ClientResponseDto `json:"data"`
}

