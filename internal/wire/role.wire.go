//go:build wireinject

package wire

import (
	"base_go_be/internal/controller"
	"base_go_be/internal/repo"
	"base_go_be/internal/service"

	"github.com/google/wire"
)

func InitRoleRouterHandler() (*controller.RoleController, error) {
	wire.Build(
		repo.NewRoleRepository,
		service.NewRoleService,
		controller.NewRoleController,
	)
	return new(controller.RoleController), nil
}

