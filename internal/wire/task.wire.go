//go:build wireinject
// +build wireinject

package wire

import (
	"base_go_be/internal/controller"
	"base_go_be/internal/repo"
	"base_go_be/internal/service"

	"github.com/google/wire"
)

func InitTaskRouterHandler() (*controller.TaskController, error) {
	wire.Build(
		repo.NewTaskRepository,
		service.NewTaskService,
		controller.NewTaskController,
	)

	return &controller.TaskController{}, nil
}

