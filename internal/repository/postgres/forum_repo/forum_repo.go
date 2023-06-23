package forum_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/blockloop/scan"
	"github.com/db-forum.git/internal/repository/postgres"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/lib/pq"
	"time"
)

type ForumRepo struct {
	db *sql.DB
}

func (f *ForumRepo) Clear() (forumErr *forum_errors.ForumError) {
	forumErr = &forum_errors.ForumError{
		Reason: nil,
		Code:   forum_errors.Unknown,
	}
	_, err := f.db.Exec(`TRUNCATE votes, posts, threads, forums, users CASCADE;`)
	if err != nil {
		forumErr.Reason = err
		return forumErr
	}

	return nil
}

func (f *ForumRepo) GetServiceStatus() (status models.Status, forumErr *forum_errors.ForumError) {
	q, err := f.db.Query("SELECT" +
		"(SELECT COUNT(*) FROM users) as user," +
		"(SELECT COUNT(*) FROM forums) as forum," +
		"(SELECT COUNT(*) FROM threads) as thread," +
		"(SELECT COUNT(*) FROM posts) as post;")

	if err != nil {
		return status, &forum_errors.ForumError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}
	err = scan.Row(&status, q)

	if err != nil {
		return status, &forum_errors.ForumError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}

	return status, nil
}

func (f *ForumRepo) GetForumThreads(slug string, limit int, since string, desc bool) (threads models.Threads, forumErr *forum_errors.ForumError) {
	threads = make(models.Threads, 0)
	values := []interface{}{slug}
	sinceQuery := ""
	if since != "" {
		sinceQuery = "and created "
		if desc {
			sinceQuery += "<= $2"
		} else {
			sinceQuery += ">= $2"
		}
		values = append(values, since)
	}

	orderByQuery := "order by created "
	if desc {
		orderByQuery += "DESC "
	}
	if limit > 0 {
		orderByQuery += fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`select id, title, author, forum, message, slug, created, votes
								from threads where forum = $1 %s %s;`, sinceQuery, orderByQuery)
	rows, err := f.db.Query(query, values...)
	forumErr = &forum_errors.ForumError{Code: forum_errors.Unknown}
	if err != nil {
		forumErr.Reason = err
		return threads, forumErr
	}
	err = scan.Rows(&threads, rows)
	if err != nil {
		forumErr.Reason = err
		return threads, forumErr
	}
	return threads, nil
}

func (f *ForumRepo) CreateForum(forum models.Forum) (createdForum models.Forum, forumErr *forum_errors.ForumError) {
	forumErr = &forum_errors.ForumError{Code: forum_errors.Unknown}
	forumQuery, err := f.db.Query(`insert into  forums (title, "user", slug) values ($1, (SELECT nickname FROM users WHERE nickname = $2), $3) returning *;`, forum.Title, forum.User, forum.Slug)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			forumErr.Reason = pqErr
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
	forumQuery, err := f.db.Query(`select slug, title, "user", threads, posts from forums where slug = $1`, slug)
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
	forumErr = &forum_errors.ForumError{Code: forum_errors.Unknown}
	if thread.Created == "" {
		thread.Created = time.Now().Format(time.RFC3339)
	}
	query, err := f.db.Query(`insert into threads (title, 
                     author, forum, message, 
                     slug, created) values 
				  	 	($1, 
				  	 	(SELECT nickname FROM users WHERE nickname = $2),
				  		(SELECT slug from forums where forums.slug = $3),
				  		$4, $5, $6) returning *;`,
		thread.Title, thread.Author, slug, thread.Message, thread.Slug, thread.Created)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			forumErr.Reason = pqErr
			switch pqErr.Code.Name() {
			case postgres.FOREIGN_KEY_VIOLATION:
				forumErr.Code = forum_errors.CantFindForum
				return createdThread, forumErr
			case postgres.NOT_NULL_VIOLATION:
				forumErr.Code = forum_errors.CantFindUser
				return createdThread, forumErr
			default:
				return createdThread, forumErr
			}
		}
	}
	err = scan.Row(&thread, query)
	if err != nil {
		forumErr.Reason = err
		return createdThread, forumErr
	}
	return thread, nil
}

func (f *ForumRepo) GetForumUsers(slug string, limit int, since string, desc bool) (users models.Users, forumErr *forum_errors.ForumError) {
	users = make(models.Users, 0)
	values := []interface{}{slug}
	sinceQuery := ""
	if since != "" {
		sinceQuery = " and nickname "
		if desc {
			sinceQuery += "< $2 "
		} else {
			sinceQuery += "> $2 "
		}
		values = append(values, since)
	}
	if desc {
		sinceQuery += " order by nickname DESC "
	} else {
		sinceQuery += " order by nickname ASC "
	}
	if limit > 0 {
		sinceQuery += fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf("select nickname, fullname, about, email from user_forums where forum = $1 %s;", sinceQuery)

	queryResult, err := f.db.Query(query, values...)

	if err != nil {
		return users, &forum_errors.ForumError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}

	err = scan.Rows(&users, queryResult)
	if err != nil {
		return users, &forum_errors.ForumError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}

	return users, nil

}

func NewForumRepo(dbConn *sql.DB) (repo *ForumRepo, err error) {
	if dbConn == nil {
		return repo, errors.New("can't create repo without connection")
	}
	return &ForumRepo{db: dbConn}, nil
}
