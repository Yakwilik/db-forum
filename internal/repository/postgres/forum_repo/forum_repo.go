package forum_repo

import (
	"database/sql"
	"errors"
	"github.com/blockloop/scan"
	"github.com/db-forum.git/internal/repository/postgres"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/lib/pq"
	"log"
)

type ForumRepo struct {
	db *sql.DB
}

func (f *ForumRepo) CreateForum(forum models.Forum) (createdForum models.Forum, forumErr *forum_errors.ForumError) {
	forumErr = &forum_errors.ForumError{Code: forum_errors.Unknown}
	forumQuery, err := f.db.Query(`insert into  forums (title, "user", slug) values ($1, (SELECT nickname FROM users WHERE nickname = $2), $3) returning *;`, forum.Title, forum.User, forum.Slug)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			forumErr.Reason = pqErr
			log.Printf("pq error: %v", pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case postgres.NOT_NULL_VIOLATION:
				fallthrough
			case postgres.FOREIGN_KEY_VIOLATION:
				forumErr.Code = forum_errors.CantFindUser
				return createdForum, forumErr
			case postgres.UNIQUE_VIOLATION:
				forumErr.Code = forum_errors.ForumAlreadyExists
				return createdForum, forumErr
			default:
				return createdForum, forumErr
			}
		}
	}
	err = scan.Row(&forum, forumQuery)
	if err != nil {
		return models.Forum{}, forumErr
	}
	return forum, nil
}

func (f *ForumRepo) GetForumInfo(slug string) (createdForum models.Forum, forumErr *forum_errors.ForumError) {
	forumErr = &forum_errors.ForumError{Code: forum_errors.Unknown}
	forumQuery, err := f.db.Query(`select slug, title, "user" from forums where slug = $1`, slug)
	if err != nil {
		forumErr.Reason = err
		return models.Forum{}, forumErr
	}
	err = scan.Row(&createdForum, forumQuery)
	if err != nil {
		forumErr.Reason = err
		if err == sql.ErrNoRows {
			forumErr.Code = forum_errors.CantFindForum
		}
		return createdForum, forumErr
	}
	return createdForum, nil
}

func (f *ForumRepo) CreateThread(slug string, thread models.Thread) (createdThread models.Thread, forumErr *forum_errors.ForumError) {
	//TODO implement me
	panic("implement me")
}

func (f *ForumRepo) GetForumUsers(slug string, limit int, since string, ask bool) (users models.Users, forumErr *forum_errors.ForumError) {
	//TODO implement me
	panic("implement me")
}

func NewForumRepo(dbConn *sql.DB) (repo *ForumRepo, err error) {
	if dbConn == nil {
		return repo, errors.New("can't create repo without connection")
	}
	return &ForumRepo{db: dbConn}, nil
}
