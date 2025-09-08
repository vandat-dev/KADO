package service

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/pkg/response"
)

type IClientService interface {
	GetClientByID(id uint) *response.ServiceResult
	GetListClient(req dto.ClientListRequestDto) *response.ServiceResult
	CreateClient(clientDto dto.ClientRequestDto) *response.ServiceResult
	UpdateClient(id uint, updateDto dto.ClientUpdateRequestDto) *response.ServiceResult
	DeleteClient(id uint) *response.ServiceResult
}

type clientService struct {
	clientRepo repo.IClientRepository
}

func NewClientService(clientRepo repo.IClientRepository) IClientService {
	return &clientService{clientRepo: clientRepo}
}

func (cs *clientService) GetClientByID(id uint) *response.ServiceResult {
	result, err := cs.clientRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeClientNotFound)
	}

	clientResponse := dto.ClientResponseDto{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}
	return response.NewServiceResult(&clientResponse)
}

func (cs *clientService) GetListClient(req dto.ClientListRequestDto) *response.ServiceResult {
	clients, total, err := cs.clientRepo.GetListClient(req)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	//var clientResponses []dto.ClientResponseDto
	//for _, client := range clients {
	//	clientResponse := dto.ClientResponseDto{
	//		ID:        client.ID,
	//		Name:      client.Name,
	//		CreatedAt: client.CreatedAt,
	//		UpdatedAt: client.UpdatedAt,
	//	}
	//	clientResponses = append(clientResponses, clientResponse)
	//}

	result := map[string]interface{}{
		"total": total,
		"data":  clients,
	}
	return response.NewServiceResult(&result)
}

func (cs *clientService) CreateClient(clientDto dto.ClientRequestDto) *response.ServiceResult {
	client := &model.Client{
		Name: clientDto.Name,
	}

	createdClient, err := cs.clientRepo.Create(client)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	clientResponse := dto.ClientResponseDto{
		ID:        createdClient.ID,
		Name:      createdClient.Name,
		CreatedAt: createdClient.CreatedAt,
		UpdatedAt: createdClient.UpdatedAt,
	}

	return response.NewServiceResult(&clientResponse)
}

func (cs *clientService) UpdateClient(id uint, updateDto dto.ClientUpdateRequestDto) *response.ServiceResult {
	existingClient, err := cs.clientRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeClientNotFound)
	}

	if updateDto.Name != "" {
		existingClient.Name = updateDto.Name
	}

	updatedClient, err := cs.clientRepo.Update(existingClient)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	clientResponse := dto.ClientResponseDto{
		ID:        updatedClient.ID,
		Name:      updatedClient.Name,
		CreatedAt: updatedClient.CreatedAt,
		UpdatedAt: updatedClient.UpdatedAt,
	}

	return response.NewServiceResult(&clientResponse)
}

func (cs *clientService) DeleteClient(id uint) *response.ServiceResult {
	_, err := cs.clientRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeClientNotFound)
	}

	err = cs.clientRepo.Delete(id)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	return response.NewServiceResult(nil)
}
