package service

import (
	"base_go_be/internal/dto"
	"base_go_be/internal/model"
	"base_go_be/internal/repo"
	"base_go_be/pkg/response"
)

type IRoleService interface {
	GetRoleByID(id uint) *response.ServiceResult
	GetListRole(req dto.RoleListRequestDto) *response.ServiceResult
	CreateRole(roleDto dto.RoleRequestDto) *response.ServiceResult
	UpdateRole(id uint, updateDto dto.RoleUpdateRequestDto) *response.ServiceResult
	DeleteRole(id uint) *response.ServiceResult
}

type roleService struct {
	roleRepo repo.IRoleRepository
}

func NewRoleService(roleRepo repo.IRoleRepository) IRoleService {
	return &roleService{roleRepo: roleRepo}
}

func (rs *roleService) GetRoleByID(id uint) *response.ServiceResult {
	result, err := rs.roleRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeRoleNotFound)
	}

	roleResponse := dto.RoleResponseDto{
		ID:        result.ID,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}
	return response.NewServiceResult(&roleResponse)
}

func (rs *roleService) GetListRole(req dto.RoleListRequestDto) *response.ServiceResult {
	roles, total, err := rs.roleRepo.GetListRole(req)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	//var roleResponses []dto.RoleResponseDto
	//for _, role := range roles {
	//	roleResponse := dto.RoleResponseDto{
	//		ID:        role.ID,
	//		Name:      role.Name,
	//		CreatedAt: role.CreatedAt,
	//		UpdatedAt: role.UpdatedAt,
	//	}
	//	roleResponses = append(roleResponses, roleResponse)
	//}

	result := map[string]interface{}{
		"total": total,
		"data":  roles,
	}

	return response.NewServiceResult(&result)
}

func (rs *roleService) CreateRole(roleDto dto.RoleRequestDto) *response.ServiceResult {
	role := &model.Role{
		Name: roleDto.Name,
	}

	createdRole, err := rs.roleRepo.Create(role)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	roleResponse := dto.RoleResponseDto{
		ID:        createdRole.ID,
		Name:      createdRole.Name,
		CreatedAt: createdRole.CreatedAt,
		UpdatedAt: createdRole.UpdatedAt,
	}

	return response.NewServiceResult(&roleResponse)
}

func (rs *roleService) UpdateRole(id uint, updateDto dto.RoleUpdateRequestDto) *response.ServiceResult {
	existingRole, err := rs.roleRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeRoleNotFound)
	}

	if updateDto.Name != "" {
		existingRole.Name = updateDto.Name
	}

	updatedRole, err := rs.roleRepo.Update(existingRole)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	roleResponse := dto.RoleResponseDto{
		ID:        updatedRole.ID,
		Name:      updatedRole.Name,
		CreatedAt: updatedRole.CreatedAt,
		UpdatedAt: updatedRole.UpdatedAt,
	}

	return response.NewServiceResult(&roleResponse)
}

func (rs *roleService) DeleteRole(id uint) *response.ServiceResult {
	_, err := rs.roleRepo.FindByID(id)
	if err != nil {
		return response.NewServiceErrorWithCode(422, response.ErrCodeRoleNotFound)
	}

	err = rs.roleRepo.Delete(id)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	return response.NewServiceResult(nil)
}
