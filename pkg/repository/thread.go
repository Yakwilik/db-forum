package repository

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
)

type Thread interface {
	GetThreadBySlugOrId(slugOrId string) (thread models.Thread, threadErr *forum_errors.ThreadError)
	CreatePosts(posts models.Posts, threadID int64, forumSlug string, createdTime string) (createdPosts models.Posts, threadErr *forum_errors.ThreadError)
	GetPostsByThreadIdFlat(threadId int64, limit int, since int, desc bool) (posts models.Posts, threadErr *forum_errors.ThreadError)
	GetPostsByThreadIdTree(threadId int64, limit int, since int, desc bool) (posts models.Posts, threadErr *forum_errors.ThreadError)
	GetPostById(postId int) (post models.Post, threadErr *forum_errors.ThreadError)
	UpdatePost(id int, update models.PostUpdate) (post models.Post, threadErr *forum_errors.ThreadError)
	GetPostsByThreadIdTreeParent(threadId int64, limit int, since int, desc bool) (posts models.Posts, threadErr *forum_errors.ThreadError)
	UpdateThread(threadID int64, update models.ThreadUpdate) (thread models.Thread, threadErr *forum_errors.ThreadError)
	VoteForThread(threadID int64, vote models.Vote) (thread models.Thread, threadErr *forum_errors.ThreadError)
}
