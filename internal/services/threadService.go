package services

import (
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/repository"
	"time"
)

type ThreadService struct {
	repo *repository.Repository
}

func (t *ThreadService) GetPosts(slugOrId string, limit int, sort string, desc bool, since int) (posts models.Posts, threadErr *forum_errors.ThreadError) {
	thread, threadErr := t.repo.GetThreadBySlugOrId(slugOrId)

	if threadErr != nil {
		return posts, threadErr
	}
	switch sort {
	case "tree":
		return t.repo.GetPostsByThreadIdTree(thread.Id, limit, since, desc)
	case "parent_tree":
		return t.repo.GetPostsByThreadIdTreeParent(thread.Id, limit, since, desc)
	default:
		return t.repo.GetPostsByThreadIdFlat(thread.Id, limit, since, desc)
	}
}

func (t *ThreadService) GetThreadBySlugOrId(slugOrId string) (thread models.Thread, threadErr *forum_errors.ThreadError) {
	return t.repo.GetThreadBySlugOrId(slugOrId)
}

func (t *ThreadService) Vote(slugOrId string, vote models.Vote) (thread models.Thread, threadErr *forum_errors.ThreadError) {
	thread, threadErr = t.repo.GetThreadBySlugOrId(slugOrId)

	if threadErr != nil {
		return thread, threadErr
	}
	return t.repo.VoteForThread(thread.Id, vote)
}

func (t *ThreadService) CreatePosts(posts models.Posts, slugOrId string) (createdPosts models.Posts, forumErr *forum_errors.ThreadError) {
	thread, threadErr := t.repo.GetThreadBySlugOrId(slugOrId)
	if threadErr != nil {
		return createdPosts, threadErr
	}
	now := time.Now().Format(time.RFC3339)

	return t.repo.CreatePosts(posts, thread.Id, thread.Forum, now)
}

func NewThreadService(repository *repository.Repository) *ThreadService {
	return &ThreadService{repo: repository}
}
