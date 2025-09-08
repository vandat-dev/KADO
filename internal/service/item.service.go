package service

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/pkg/response"
)

type IItemService interface {
	GetItemByID(id uint) *response.ServiceResult
	GetListItem(req dto.ItemListRequestDto) *response.ServiceResult
	CreateItem(itemDto dto.ItemRequestDto) *response.ServiceResult
	UpdateItem(id uint, updateDto dto.ItemUpdateRequestDto) *response.ServiceResult
	DeleteItem(id uint) *response.ServiceResult
}

type itemService struct {
	itemRepo repo.IItemRepository
}

func NewItemService(itemRepo repo.IItemRepository) IItemService {
	return &itemService{itemRepo: itemRepo}
}

func (is *itemService) GetItemByID(id uint) *response.ServiceResult {
	result, err := is.itemRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeItemNotFound)
	}

	itemResponse := dto.ItemResponseDto{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Category:    result.Category,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}
	return response.NewServiceResult(&itemResponse)
}

func (is *itemService) GetListItem(req dto.ItemListRequestDto) *response.ServiceResult {
	items, total, err := is.itemRepo.GetListItem(req)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	//var itemResponses []dto.ItemResponseDto
	//for _, item := range items {
	//	itemResponse := dto.ItemResponseDto{
	//		ID:          item.ID,
	//		Name:        item.Name,
	//		Description: item.Description,
	//		Category:    item.Category,
	//		CreatedAt:   item.CreatedAt,
	//		UpdatedAt:   item.UpdatedAt,
	//	}
	//	itemResponses = append(itemResponses, itemResponse)
	//}
	//
	//result := dto.ItemListResponseDto{
	//	Total: total,
	//	Data:  itemResponses,
	//}
	result := map[string]interface{}{
		"total": total,
		"data":  items,
	}

	return response.NewServiceResult(&result)
}

func (is *itemService) CreateItem(itemDto dto.ItemRequestDto) *response.ServiceResult {
	item := &model.Item{
		Name:        itemDto.Name,
		Description: itemDto.Description,
		Category:    itemDto.Category,
	}

	createdItem, err := is.itemRepo.Create(item)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	itemResponse := dto.ItemResponseDto{
		ID:          createdItem.ID,
		Name:        createdItem.Name,
		Description: createdItem.Description,
		Category:    createdItem.Category,
		CreatedAt:   createdItem.CreatedAt,
		UpdatedAt:   createdItem.UpdatedAt,
	}

	return response.NewServiceResult(&itemResponse)
}

func (is *itemService) UpdateItem(id uint, updateDto dto.ItemUpdateRequestDto) *response.ServiceResult {
	existingItem, err := is.itemRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeItemNotFound)
	}

	if updateDto.Name != "" {
		existingItem.Name = updateDto.Name
	}
	if updateDto.Description != "" {
		existingItem.Description = updateDto.Description
	}
	if updateDto.Category != "" {
		existingItem.Category = updateDto.Category
	}

	updatedItem, err := is.itemRepo.Update(existingItem)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	itemResponse := dto.ItemResponseDto{
		ID:          updatedItem.ID,
		Name:        updatedItem.Name,
		Description: updatedItem.Description,
		Category:    updatedItem.Category,
		CreatedAt:   updatedItem.CreatedAt,
		UpdatedAt:   updatedItem.UpdatedAt,
	}

	return response.NewServiceResult(&itemResponse)
}

func (is *itemService) DeleteItem(id uint) *response.ServiceResult {
	_, err := is.itemRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeItemNotFound)
	}

	err = is.itemRepo.Delete(id)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	return response.NewServiceResult(nil)
}
