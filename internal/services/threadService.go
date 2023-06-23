package services

import (
	"fmt"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/repository"
	"time"
)

type ThreadService struct {
	repo *repository.Repository
}

func (t *ThreadService) UpdatePost(id int, update models.PostUpdate) (post models.Post, threadErr *forum_errors.ThreadError) {
	if update.Message == "" {
		return t.repo.GetPostById(id)
	}
	return t.repo.UpdatePost(id, update)
}

func (t *ThreadService) GetPost(id int, relatedUser, relatedThread, relatedForum bool) (post models.PostFull, threadErr *forum_errors.ThreadError) {
	onlyPost, threadErr := t.repo.GetPostById(id)

	if threadErr != nil {
		return post, threadErr
	}
	post.Post = onlyPost

	if relatedUser {
		author, threadErr := t.repo.GetUser(onlyPost.Author)
		if threadErr != nil {
			return post, (*forum_errors.ThreadError)(threadErr)
		}
		post.User = &author
	}
	if relatedThread {
		thread, threadErr := t.repo.GetThreadBySlugOrId(fmt.Sprintf("%d", onlyPost.Thread))
		if threadErr != nil {
			return post, threadErr
		}
		post.Thread = &thread
	}

	if relatedForum {
		forum, threadErr := t.repo.GetForumInfo(onlyPost.Forum)
		if threadErr != nil {
			return post, (*forum_errors.ThreadError)(threadErr)
		}
		post.Forum = &forum
	}
	return post, nil
}

func (t *ThreadService) UpdateThread(slugOrId string, update models.ThreadUpdate) (thread models.Thread, threadErr *forum_errors.ThreadError) {
	thread, threadErr = t.GetThreadBySlugOrId(slugOrId)

	if threadErr != nil {
		return thread, threadErr
	}
	if update.Message != "" {
		thread.Message = update.Message
	}
	if update.Title != "" {
		thread.Title = update.Title
	}

	return t.repo.UpdateThread(thread.Id, thread.ThreadUpdate)
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
	if len(posts) == 0 {
		return make(models.Posts, 0), nil
	}

	if posts[0].Parent != 0 {
		parent, threadErr := t.GetPost(int(posts[0].Parent), false, false, false)
		if threadErr != nil {
			return nil, threadErr
		}

		if int64(parent.Post.Thread) != thread.Id {
			return nil, &forum_errors.ThreadError{
				Reason: fmt.Errorf("conflict"),
				Code:   forum_errors.ConflictingData,
			}
		}
	}
	now := time.Now().Format(time.RFC3339)

	return t.repo.CreatePosts(posts, thread.Id, thread.Forum, now)
}

func NewThreadService(repository *repository.Repository) *ThreadService {
	return &ThreadService{repo: repository}
}
