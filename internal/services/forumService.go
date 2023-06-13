package services

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/repository"
)

type ForumService struct {
	repo *repository.Repository
}

func (f *ForumService) CreateForum(forum models.Forum) (createdForum models.Forum, forumErr *forum_errors.ForumError) {
	return f.repo.CreateForum(forum)
}

func (f *ForumService) GetForumInfo(slug string) (createdForum models.Forum, forumErr *forum_errors.ForumError) {
	return f.repo.GetForumInfo(slug)
}

func (f *ForumService) CreateThread(slug string, thread models.Thread) (createdThread models.Thread, forumErr *forum_errors.ForumError) {
	//TODO implement me
	panic("implement me")
}

func (f *ForumService) GetForumUsers(slug string, limit int, since string, ask bool) (users models.Users, forumErr *forum_errors.ForumError) {
	//TODO implement me
	panic("implement me")
}

func NewForumService(repository *repository.Repository) *ForumService {
	return &ForumService{repo: repository}
}
