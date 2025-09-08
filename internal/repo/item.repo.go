package repo

import (
	"base_go_be/global"
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"errors"

	"gorm.io/gorm"
)

type IItemRepository interface {
	FindByID(id uint) (*model.Item, error)
	FindAll() ([]model.Item, error)
	GetListItem(req dto.ItemListRequestDto) ([]model.Item, int64, error)
	Create(item *model.Item) (*model.Item, error)
	Update(item *model.Item) (*model.Item, error)
	Delete(id uint) error
}

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository() IItemRepository {
	return &ItemRepository{db: global.Postgres}
}

func (ir *ItemRepository) FindByID(id uint) (*model.Item, error) {
	var item model.Item
	if err := ir.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (ir *ItemRepository) FindAll() ([]model.Item, error) {
	var items []model.Item
	if err := ir.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (ir *ItemRepository) GetListItem(req dto.ItemListRequestDto) ([]model.Item, int64, error) {
	var items []model.Item
	var total int64

	// Build query with filters
	query := ir.db.Model(&model.Item{})

	// Apply filters
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
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

	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (ir *ItemRepository) Create(item *model.Item) (*model.Item, error) {
	if err := ir.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (ir *ItemRepository) Update(item *model.Item) (*model.Item, error) {
	if err := ir.db.Save(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (ir *ItemRepository) Delete(id uint) error {
	if err := ir.db.Delete(&model.Item{}, id).Error; err != nil {
		return err
	}
	return nil
}
