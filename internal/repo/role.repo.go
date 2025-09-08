package repo

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"errors"

	"gorm.io/gorm"
)

type IRoleRepository interface {
	FindByID(id uint) (*model.Role, error)
	FindAll() ([]model.Role, error)
	GetListRole(req dto.RoleListRequestDto) ([]model.Role, int64, error)
	Create(role *model.Role) (*model.Role, error)
	Update(role *model.Role) (*model.Role, error)
	Delete(id uint) error
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() IRoleRepository {
	return &RoleRepository{db: global.Postgres}
}

func (rr *RoleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := rr.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (rr *RoleRepository) FindAll() ([]model.Role, error) {
	var roles []model.Role
	if err := rr.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (rr *RoleRepository) GetListRole(req dto.RoleListRequestDto) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	// Build query with filters
	query := rr.db.Model(&model.Role{})

	// Apply filters
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	} else {
		query = query.Limit(10)
	}

	query = query.Offset(req.Skip).Order("created_at DESC")

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (rr *RoleRepository) Create(role *model.Role) (*model.Role, error) {
	if err := rr.db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (rr *RoleRepository) Update(role *model.Role) (*model.Role, error) {
	if err := rr.db.Save(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (rr *RoleRepository) Delete(id uint) error {
	if err := rr.db.Delete(&model.Role{}, id).Error; err != nil {
		return err
	}
	return nil
}
