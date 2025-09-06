package service

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/pkg/config"
	"base_go_be/pkg/jwt"
	"base_go_be/pkg/response"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUserByID(id uint) *response.ServiceResult
	GetListUser(req dto.UserListRequestDto, userRole string) *response.ServiceResult
	CreateUser(email string, username string, password string, role string) *response.ServiceResult
	UpdateUser(id uint, updateDto dto.UserUpdateRequestDto) *response.ServiceResult
	Login(email string, password string) *response.ServiceResult
	Register(registerDto dto.RegisterRequestDto) *response.ServiceResult
}

type userService struct {
	userRepo repo.IUserRepository
}

func NewUserService(userRepo repo.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

func (us *userService) GetUserByID(id uint) *response.ServiceResult {
	result := us.userRepo.GetUserByID(id)
	if result == nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeUserNotFound)
	}
	userResponse := dto.UserResponseDto{
		Id:       result.ID,
		Email:    result.Email,
		Username: result.Username,
		Role:     result.Role,
	}
	return response.NewServiceResult(&userResponse)
}

func (us *userService) GetListUser(req dto.UserListRequestDto, userRole string) *response.ServiceResult {
	// Check authorization - only ADMIN can get user list
	if userRole != "ADMIN" {
		return response.NewServiceErrorWithCode(403, response.ErrCodeAccessDenied)
	}

	users, total, err := us.userRepo.GetListUser(req)
	if err != nil {
		global.Logger.Error("Failed to get users from repository: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Convert to DTOs
	userDTOs := make([]dto.UserResponseDto, 0, len(users))
	for _, u := range users {
		userDTOs = append(userDTOs, dto.UserResponseDto{
			Id:       u.ID,
			Email:    u.Email,
			Username: u.Username,
			Role:     u.Role,
		})
	}

	result := &dto.UserListResponseDto{
		Data:  userDTOs,
		Total: total,
	}
	return response.NewServiceResult(result)
}

func (us *userService) CreateUser(email string, username string, password string, role string) *response.ServiceResult {

	existingUser := us.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return response.NewServiceErrorWithCode(409, response.ErrCodeUserHasExists)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	user := &model.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	userID, err := us.userRepo.CreateUser(user)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}
	return response.NewServiceResult(userID)
}

func (us *userService) UpdateUser(id uint, updateDto dto.UserUpdateRequestDto) *response.ServiceResult {

	existingUser := us.userRepo.GetUserByID(id)
	if existingUser == nil {
		return response.NewServiceErrorWithCode(404, response.ErrCodeUserNotFound)
	}

	updateUser := &model.User{}

	if updateDto.Username != "" {
		updateUser.Username = updateDto.Username
	}
	if updateDto.Role != "" {
		updateUser.Role = updateDto.Role
	}

	if updateDto.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateDto.Password), bcrypt.DefaultCost)
		if err != nil {
			return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
		}
		updateUser.Password = string(hashedPassword)
	}

	updatedUser, err := us.userRepo.UpdateUser(id, updateUser)
	if err != nil {
		return response.NewServiceErrorWithCode(400, response.ErrCodeUserHasExists)
	}

	userResponse := dto.UserResponseDto{
		Id:       updatedUser.ID,
		Email:    updatedUser.Email,
		Username: updatedUser.Username,
		Role:     updatedUser.Role,
	}

	return response.NewServiceResult(&userResponse)
}

func (us *userService) Login(email string, password string) *response.ServiceResult {
	user := us.userRepo.GetUserByEmail(email)
	if user == nil {
		return response.NewServiceErrorWithCode(401, response.ErrCodeInvalidLogin)
	}

	// Compare password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return response.NewServiceErrorWithCode(401, response.ErrCodeInvalidLogin)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role, config.JWT.SecretKey, config.JWT.TokenExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate refresh token with longer expiry
	refreshToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role, config.JWT.SecretKey, config.JWT.RefreshExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	authResponse := &dto.AuthResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
		User: dto.UserResponseDto{
			Id:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
	}

	return response.NewServiceResult(authResponse)
}

func (us *userService) Register(registerDto dto.RegisterRequestDto) *response.ServiceResult {
	createResult := us.CreateUser(registerDto.Email, registerDto.Username, registerDto.Password, registerDto.Role)
	if createResult.Error != nil {
		return createResult // Return CreateUser
	}

	userID := createResult.Data.(uint)
	user := us.userRepo.GetUserByID(userID)
	if user == nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role, config.JWT.SecretKey, config.JWT.TokenExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate refresh token with longer expiry
	refreshToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role, config.JWT.SecretKey, config.JWT.RefreshExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	authResponse := &dto.AuthResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
		User: dto.UserResponseDto{
			Id:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
	}

	return response.NewServiceResult(authResponse)
}
