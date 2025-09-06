package repo

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(email string) *model.User
	GetUserByID(id uint) *model.User
	GetListUser(req dto.UserListRequestDto) ([]*model.User, int64, error)
	CreateUser(user *model.User) (uint, error)
	UpdateUser(id uint, user *model.User) (*model.User, error)
}

func NewUserRepository() IUserRepository {
	return &userRepository{db: global.Postgres}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) GetUserByEmail(email string) *model.User {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func (r *userRepository) GetUserByID(id uint) *model.User {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil
	}
	return &user
}

func (r *userRepository) GetListUser(req dto.UserListRequestDto) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{})

	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+req.Email+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(req.Skip).Limit(req.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) CreateUser(user *model.User) (uint, error) {
	err := r.db.Create(user).Error
	return user.ID, err
}

func (r *userRepository) UpdateUser(id uint, user *model.User) (*model.User, error) {
	// run update
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return nil, err
	}

	var updatedUser model.User
	if err := r.db.First(&updatedUser, id).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
