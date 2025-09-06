package routers

import (
	"base_go_be/internal/routers/user"
)

type RouterGroup struct {
	User user.UsersRouterGroup
}

var RouterGroupApp = new(RouterGroup)
