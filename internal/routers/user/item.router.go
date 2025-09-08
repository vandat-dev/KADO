package user

import (
	"base_go_be/internal/middlewares"
	"base_go_be/internal/wire"

	"github.com/gin-gonic/gin"
)

type ItemRouter struct{}

func (ir *ItemRouter) InitItemRouter(Router *gin.RouterGroup) {
	// WIRE go - get item controller with dependency injection
	itemController, _ := wire.InitItemRouterHandler()

	// private router - authentication required
	itemRouterPrivate := Router.Group("/item")
	itemRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		itemRouterPrivate.GET("/detail/:id", itemController.GetItemByID)
		itemRouterPrivate.GET("/list", itemController.GetListItem)
		itemRouterPrivate.POST("/create", itemController.CreateItem)
		itemRouterPrivate.PUT("/update/:id", itemController.UpdateItem)
		itemRouterPrivate.DELETE("/delete/:id", itemController.DeleteItem)
	}
}
