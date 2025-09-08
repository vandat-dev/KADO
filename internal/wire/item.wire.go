//go:build wireinject

package wire

import (
	"base_go_be/internal/controller"
	"base_go_be/internal/repo"
	"base_go_be/internal/service"

	"github.com/google/wire"
)

func InitItemRouterHandler() (*controller.ItemController, error) {
	wire.Build(
		repo.NewItemRepository,
		service.NewItemService,
		controller.NewItemController,
	)
	return new(controller.ItemController), nil
}

