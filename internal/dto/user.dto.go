package dto

type UserRequestDto struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

type UserUpdateRequestDto struct {
	Username string `json:"username" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Role     string `json:"role" binding:"omitempty"`
}

type UserResponseDto struct {
	Id       uint   `json:"id" binding:"required"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
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
