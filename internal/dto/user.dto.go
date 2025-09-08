package dto

import "time"

type UserRequestDto struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	FullName    string `json:"full_name"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	Address     string `json:"address"`
	SystemRole  string `json:"system_role" binding:"required,oneof=ADMIN USER SUPER_ADMIN"`
}

type UserUpdateRequestDto struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	Address     string `json:"address"`
	SystemRole  string `json:"system_role" binding:"omitempty,oneof=ADMIN USER SUPER_ADMIN"`
	IsActive    *bool  `json:"is_active"`
}

type UserResponseDto struct {
	Id          uint      `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	SystemRole  string    `json:"system_role"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserResponseBaseDto for basic user information in responses
type UserResponseBaseDto struct {
	Id         uint   `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	SystemRole string `json:"system_role"`
}

// UserListRequestDto for pagination and filtering
type UserListRequestDto struct {
	Skip  int    `form:"skip" binding:"min=0"`
	Limit int    `form:"limit" binding:"min=0,max=100"`
	Email string `form:"email"`
}

// UserListResponseDto for paginated user list response
type UserListResponseDto struct {
	Total int64             `json:"total"`
	Data  []UserResponseDto `json:"data"`
}
