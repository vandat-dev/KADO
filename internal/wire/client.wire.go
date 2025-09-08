//go:build wireinject

package wire

import (
	"base_go_be/internal/controller"
	"base_go_be/internal/repo"
	"base_go_be/internal/service"

	"github.com/google/wire"
)

func InitClientRouterHandler() (*controller.ClientController, error) {
	wire.Build(
		repo.NewClientRepository,
		service.NewClientService,
		controller.NewClientController,
	)
	return new(controller.ClientController), nil
}

