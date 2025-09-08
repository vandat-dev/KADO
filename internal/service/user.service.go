package service

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/pkg/jwt"
	"base_go_be/pkg/response"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUserByID(id uint) *response.ServiceResult
	GetListUser(req dto.UserListRequestDto, userRole string) *response.ServiceResult
	CreateUser(userDto dto.UserRequestDto) *response.ServiceResult
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
		Id:          result.ID,
		Email:       result.Email,
		Username:    result.Username,
		FullName:    result.FullName,
		PhoneNumber: result.PhoneNumber,
		Gender:      result.Gender,
		Address:     result.Address,
		SystemRole:  result.SystemRole,
		IsActive:    result.IsActive,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
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
	//userDTOs := make([]dto.UserResponseDto, 0, len(users))
	//for _, u := range users {
	//	userDTOs = append(userDTOs, dto.UserResponseDto{
	//		Id:          u.ID,
	//		Email:       u.Email,
	//		Username:    u.Username,
	//		FullName:    u.FullName,
	//		PhoneNumber: u.PhoneNumber,
	//		Gender:      u.Gender,
	//		Address:     u.Address,
	//		SystemRole:  u.SystemRole,
	//		IsActive:    u.IsActive,
	//		CreatedAt:   u.CreatedAt,
	//		UpdatedAt:   u.UpdatedAt,
	//	})
	//}
	//
	//result := &dto.UserListResponseDto{
	//	Data:  userDTOs,
	//	Total: total,
	//}
	result := map[string]interface{}{
		"total": total,
		"data":  users,
	}
	return response.NewServiceResult(result)
}

func (us *userService) CreateUser(userDto dto.UserRequestDto) *response.ServiceResult {

	existingUser := us.userRepo.GetUserByEmail(userDto.Email)
	if existingUser != nil {
		return response.NewServiceErrorWithCode(409, response.ErrCodeUserHasExists)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	user := &model.User{
		Email:       userDto.Email,
		Username:    userDto.Username,
		FullName:    userDto.FullName,
		Password:    string(hashedPassword),
		PhoneNumber: userDto.PhoneNumber,
		Gender:      userDto.Gender,
		Address:     userDto.Address,
		SystemRole:  userDto.SystemRole,
		IsActive:    true,
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
	if updateDto.FullName != "" {
		updateUser.FullName = updateDto.FullName
	}
	if updateDto.PhoneNumber != "" {
		updateUser.PhoneNumber = updateDto.PhoneNumber
	}
	if updateDto.Gender != "" {
		updateUser.Gender = updateDto.Gender
	}
	if updateDto.Address != "" {
		updateUser.Address = updateDto.Address
	}
	if updateDto.SystemRole != "" {
		updateUser.SystemRole = updateDto.SystemRole
	}
	if updateDto.IsActive != nil {
		updateUser.IsActive = *updateDto.IsActive
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
		Id:          updatedUser.ID,
		Email:       updatedUser.Email,
		Username:    updatedUser.Username,
		FullName:    updatedUser.FullName,
		PhoneNumber: updatedUser.PhoneNumber,
		Gender:      updatedUser.Gender,
		Address:     updatedUser.Address,
		SystemRole:  updatedUser.SystemRole,
		IsActive:    updatedUser.IsActive,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
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
	token, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.TokenExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate refresh token with longer expiry
	refreshToken, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.RefreshExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	authResponse := &dto.AuthResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
	}

	return response.NewServiceResult(authResponse)
}

func (us *userService) Register(registerDto dto.RegisterRequestDto) *response.ServiceResult {
	userDto := dto.UserRequestDto{
		Email:      registerDto.Email,
		Username:   registerDto.Username,
		Password:   registerDto.Password,
		SystemRole: registerDto.Role,
	}
	createResult := us.CreateUser(userDto)
	if createResult.Error != nil {
		return createResult // Return CreateUser
	}

	userID := createResult.Data.(uint)
	user := us.userRepo.GetUserByID(userID)
	if user == nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.TokenExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate refresh token with longer expiry
	refreshToken, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.RefreshExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	authResponse := &dto.AuthResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
	}

	return response.NewServiceResult(authResponse)
}
