package services

type Services struct {
	User
	Forum
	Thread
}

func New(userService User, forumService Forum, threadService Thread) *Services {
	return &Services{User: userService, Forum: forumService, Thread: threadService}
}
