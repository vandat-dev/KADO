package user

type UsersRouterGroup struct {
	UsersRouter
	TaskRouter
	ClientRouter
	JobRouter
	RoleRouter
	ItemRouter
}
