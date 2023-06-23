package services

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/repository"
)

type ForumService struct {
	repo *repository.Repository
}

func (f *ForumService) Clear() (forumErr *forum_errors.ForumError) {
	return f.repo.Clear()
}

func (f *ForumService) GetServiceStatus() (status models.Status, forumErr *forum_errors.ForumError) {
	return f.repo.GetServiceStatus()
}

func (f *ForumService) GetForumThreads(slug string, limit int, since string, desk bool) (threads models.Threads, forumErr *forum_errors.ForumError) {
	_, forumErr = f.repo.GetForumInfo(slug)
	if forumErr != nil {
		return threads, forumErr
	}
	return f.repo.GetForumThreads(slug, limit, since, desk)
}

func (f *ForumService) CreateForum(forum models.Forum) (createdForum models.Forum, forumErr *forum_errors.ForumError) {
	return f.repo.CreateForum(forum)
}

func (f *ForumService) GetForumInfo(slug string) (createdForum models.Forum, forumErr *forum_errors.ForumError) {
	return f.repo.GetForumInfo(slug)
}

func (f *ForumService) CreateThread(slug string, thread models.Thread) (createdThread models.Thread, forumErr *forum_errors.ForumError) {
	if thread.Slug != "" {
		thread, threadErr := f.repo.GetThreadBySlugOrId(thread.Slug)
		if threadErr == nil {
			forumErr = &forum_errors.ForumError{
				Reason: threadErr,
				Code:   forum_errors.ThreadAlreadyExists,
			}
			return thread, forumErr
		}
	}
	return f.repo.CreateThread(slug, thread)
}

func (f *ForumService) GetForumUsers(slug string, limit int, since string, desc bool) (users models.Users, forumErr *forum_errors.ForumError) {
	_, forumErr = f.repo.GetForumInfo(slug)
	if forumErr != nil {
		return users, forumErr
	}
	return f.repo.GetForumUsers(slug, limit, since, desc)
}

func NewForumService(repository *repository.Repository) *ForumService {
	return &ForumService{repo: repository}
}
