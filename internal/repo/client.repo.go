package repo

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"errors"

	"gorm.io/gorm"
)

type IClientRepository interface {
	FindByID(id uint) (*model.Client, error)
	FindAll() ([]model.Client, error)
	GetListClient(req dto.ClientListRequestDto) ([]model.Client, int64, error)
	Create(client *model.Client) (*model.Client, error)
	Update(client *model.Client) (*model.Client, error)
	Delete(id uint) error
}

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository() IClientRepository {
	return &ClientRepository{db: global.Postgres}
}

func (cr *ClientRepository) FindByID(id uint) (*model.Client, error) {
	var client model.Client
	if err := cr.db.First(&client, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &client, nil
}

func (cr *ClientRepository) FindAll() ([]model.Client, error) {
	var clients []model.Client
	if err := cr.db.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (cr *ClientRepository) GetListClient(req dto.ClientListRequestDto) ([]model.Client, int64, error) {
	var clients []model.Client
	var total int64

	// Build query with filters
	query := cr.db.Model(&model.Client{})

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

	if err := query.Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

func (cr *ClientRepository) Create(client *model.Client) (*model.Client, error) {
	if err := cr.db.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (cr *ClientRepository) Update(client *model.Client) (*model.Client, error) {
	if err := cr.db.Save(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (cr *ClientRepository) Delete(id uint) error {
	if err := cr.db.Delete(&model.Client{}, id).Error; err != nil {
		return err
	}
	return nil
}

