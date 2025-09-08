package dto

import "time"

type ItemRequestDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ItemUpdateRequestDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ItemResponseDto struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ItemListRequestDto for pagination and filtering
type ItemListRequestDto struct {
	Skip     int    `form:"skip" binding:"min=0"`
	Limit    int    `form:"limit" binding:"min=0,max=100"`
	Name     string `form:"name"`
	Category string `form:"category"`
}

// ItemListResponseDto for paginated item list response
type ItemListResponseDto struct {
	Total int64             `json:"total"`
	Data  []ItemResponseDto `json:"data"`
}
