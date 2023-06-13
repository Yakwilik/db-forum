package repository

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
)

type User interface {
	CreateUser(user models.User) (createdUser models.User, userErr *forum_errors.UserError)
	GetExistingUsers(user models.User) (users models.Users, err error)
	GetUser(nickname string) (user models.User, err *forum_errors.UserError)
	UpdateUser(user models.User) (updatedUser models.User, err *forum_errors.UserError)
}
