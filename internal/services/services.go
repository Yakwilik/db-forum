package services

import "github.com/db-forum.git/pkg/services"

func New(user services.User, forum services.Forum) *services.Services {
	return services.New(user, forum)
}
