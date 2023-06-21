package repository

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
)

type Forum interface {
	CreateForum(user models.Forum) (createdForum models.Forum, forumErr *forum_errors.ForumError)
	GetForumInfo(slug string) (createdForum models.Forum, forumErr *forum_errors.ForumError)
	CreateThread(slug string, thread models.Thread) (createdThread models.Thread, forumErr *forum_errors.ForumError)
	GetForumUsers(slug string, limit int, since string, desk bool) (users models.Users, forumErr *forum_errors.ForumError)
	GetForumThreads(slug string, limit int, since string, desk bool) (threads models.Threads, forumErr *forum_errors.ForumError)
}
