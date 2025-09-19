package user

import (
	"base_go_be/internal/middlewares"
	"base_go_be/internal/wire"

	"github.com/gin-gonic/gin"
)

type TaskRouter struct{}

func (tr *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) {
	// WIRE go - get task controller with dependency injection
	taskController, _ := wire.InitTaskRouterHandler()

	// private router - authentication required
	taskRouterPrivate := Router.Group("/task")
	taskRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		taskRouterPrivate.GET("/detail/:id", taskController.GetTaskByID)
		taskRouterPrivate.GET("/list", taskController.GetListTask)
		taskRouterPrivate.GET("/my_tasks", taskController.GetMyTasks)
		taskRouterPrivate.POST("/create", taskController.CreateTask)
		taskRouterPrivate.PUT("/update/:id", taskController.UpdateTask)
		taskRouterPrivate.PUT("/update_progress/:id", taskController.UpdateProgressTask)
		taskRouterPrivate.DELETE("/delete/:id", taskController.DeleteTask)
	}
}
