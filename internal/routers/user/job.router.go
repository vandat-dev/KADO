package user

import (
	"base_go_be/internal/middlewares"
	"base_go_be/internal/wire"

	"github.com/gin-gonic/gin"
)

type JobRouter struct{}

func (jr *JobRouter) InitJobRouter(Router *gin.RouterGroup) {
	// WIRE go - get job controller with dependency injection
	jobController, _ := wire.InitJobRouterHandler()

	// private router - authentication required
	jobRouterPrivate := Router.Group("/job")
	jobRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		jobRouterPrivate.GET("/detail/:id", jobController.GetJobByID)
		jobRouterPrivate.GET("/list", jobController.GetListJob)
		jobRouterPrivate.POST("/create", jobController.CreateJob)
		jobRouterPrivate.PUT("/update/:id", jobController.UpdateJob)
		jobRouterPrivate.DELETE("/delete/:id", jobController.DeleteJob)
	}
}
