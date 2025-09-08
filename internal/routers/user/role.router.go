package user

import (
	"base_go_be/internal/middlewares"
	"base_go_be/internal/wire"

	"github.com/gin-gonic/gin"
)

type RoleRouter struct{}

func (rr *RoleRouter) InitRoleRouter(Router *gin.RouterGroup) {
	// WIRE go - get role controller with dependency injection
	roleController, _ := wire.InitRoleRouterHandler()

	// private router - authentication required
	roleRouterPrivate := Router.Group("/role")
	roleRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		roleRouterPrivate.GET("/detail/:id", roleController.GetRoleByID)
		roleRouterPrivate.GET("/list", roleController.GetListRole)
		roleRouterPrivate.POST("/create", roleController.CreateRole)
		roleRouterPrivate.PUT("/update/:id", roleController.UpdateRole)
		roleRouterPrivate.DELETE("/delete/:id", roleController.DeleteRole)
	}
}
