//go:build wireinject

package wire

import (
	"base_go_be/internal/controller"
	"base_go_be/internal/repo"
	"base_go_be/internal/service"

	"github.com/google/wire"
)

func InitJobRouterHandler() (*controller.JobController, error) {
	wire.Build(
		repo.NewJobRepository,
		service.NewJobService,
		controller.NewJobController,
	)
	return new(controller.JobController), nil
}

