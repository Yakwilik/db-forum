package services

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
)

type Thread interface {
	CreatePosts(posts models.Posts, slugOrId string) (createdPosts models.Posts, threadErr *forum_errors.ThreadError)
	Vote(slugOrId string, vote models.Vote) (thread models.Thread, threadErr *forum_errors.ThreadError)
	GetThreadBySlugOrId(slugOrId string) (thread models.Thread, threadErr *forum_errors.ThreadError)
	GetPosts(slugOrId string, limit int, sort string, desk bool, since int) (posts models.Posts, threadErr *forum_errors.ThreadError)
}
