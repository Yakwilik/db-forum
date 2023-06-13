package services

type Services struct {
	User
	Forum
}

func New(userService User, forumService Forum) *Services {
	return &Services{User: userService, Forum: forumService}
}
