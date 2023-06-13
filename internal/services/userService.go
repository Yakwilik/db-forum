package services

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/repository"
)

type UserService struct {
	repo *repository.Repository
}

func (u *UserService) CreateUser(user models.User) (createdUser models.User, userError *forum_errors.UserError) {
	return u.repo.CreateUser(user)

}

func (u *UserService) GetExistingUsers(user models.User) (users models.Users, err error) {
	return u.repo.GetExistingUsers(user)
}

func (u *UserService) GetUser(nickname string) (user models.User, err *forum_errors.UserError) {
	return u.repo.GetUser(nickname)
}

func (u *UserService) UpdateUser(user models.User) (updatedUser models.User, userErr *forum_errors.UserError) {
	updatedUser, userErr = u.GetUser(user.Nickname)
	if userErr != nil {
		return updatedUser, userErr
	}
	if user.About != "" {
		updatedUser.About = user.About
	}
	if user.Email != "" {
		updatedUser.Email = user.Email
	}
	if user.Fullname != "" {
		updatedUser.Fullname = user.Fullname
	}
	return u.repo.UpdateUser(updatedUser)
}

func NewUserService(repository *repository.Repository) *UserService {
	return &UserService{repo: repository}
}
