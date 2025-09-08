package user

import (
	"base_go_be/internal/middlewares"
	"base_go_be/internal/wire"

	"github.com/gin-gonic/gin"
)

type ClientRouter struct{}

func (cr *ClientRouter) InitClientRouter(Router *gin.RouterGroup) {
	// WIRE go - get client controller with dependency injection
	clientController, _ := wire.InitClientRouterHandler()

	// private router - authentication required
	clientRouterPrivate := Router.Group("/client")
	clientRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		clientRouterPrivate.GET("/detail/:id", clientController.GetClientByID)
		clientRouterPrivate.GET("/list", clientController.GetListClient)
		clientRouterPrivate.POST("/create", clientController.CreateClient)
		clientRouterPrivate.PUT("/update/:id", clientController.UpdateClient)
		clientRouterPrivate.DELETE("/delete/:id", clientController.DeleteClient)
	}
}
