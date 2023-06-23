package thread_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/blockloop/scan"
	"github.com/db-forum.git/internal/repository/postgres"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/lib/pq"
	"strconv"
)

type ThreadRepo struct {
	db *sql.DB
}

func (t *ThreadRepo) UpdatePost(id int, update models.PostUpdate) (post models.Post, threadErr *forum_errors.ThreadError) {
	q, err := t.db.Query("update posts set message = $1 where id = $2 returning id, author, message, is_edited, forum, thread_id as thread, created", update.Message, id)

	if err != nil {
		return post, &forum_errors.ThreadError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}

	err = scan.Row(&post, q)

	if err != nil {
		return post, &forum_errors.ThreadError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}
	return post, nil
}

func (t *ThreadRepo) GetPostById(postId int) (post models.Post, threadErr *forum_errors.ThreadError) {
	threadErr = &forum_errors.ThreadError{
		Code: forum_errors.Unknown,
	}
	q, err := t.db.Query("select "+
		"id, parent, author, message, is_edited, forum, thread_id as thread, created, path "+
		"from posts where id = $1 limit 1;", postId)
	if err != nil {
		threadErr.Reason = err
		return post, threadErr
	}

	err = scan.Row(&post, q)
	if err != nil {
		threadErr.Reason = err
		if err == sql.ErrNoRows {
			threadErr.Code = forum_errors.CantFindPost
		}
		return post, threadErr
	}

	return post, nil

}

func (t *ThreadRepo) UpdateThread(threadID int64, update models.ThreadUpdate) (thread models.Thread, threadErr *forum_errors.ThreadError) {
	query, err := t.db.Query(`update threads set title = $1, message = $2 where id = $3 returning *`, update.Title, update.Message, threadID)
	if err != nil {
		return thread, &forum_errors.ThreadError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}

	err = scan.Row(&thread, query)

	if err != nil {
		return thread, &forum_errors.ThreadError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}
	return thread, nil
}

func (t *ThreadRepo) GetPostsByThreadIdTreeParent(threadId int64, limit int, since int, desc bool) (posts models.Posts, threadErr *forum_errors.ThreadError) {
	posts = make(models.Posts, 0)
	values := []interface{}{threadId}
	sinceQuery := ""
	if since >= 0 {
		sinceQuery = "and path[1] "
		if desc {
			sinceQuery += " <  "
		} else {
			sinceQuery += " > "
		}
		sinceQuery += "(SELECT path[1] FROM posts WHERE id = $2) "

		values = append(values, since)
	}

	innerOrderByQuery := "order by id"
	outerOrderByQuery := "order by "
	if desc {
		innerOrderByQuery += " DESC "
		outerOrderByQuery += "path[1] DESC, "
	} else {
		innerOrderByQuery += " ASC "
	}
	outerOrderByQuery += "path ASC, id ASC"
	if limit > 0 {
		innerOrderByQuery += fmt.Sprintf("LIMIT %d )", limit)
	}

	query := fmt.Sprintf(`select id, parent, author, forum, message, is_edited, created, thread_id as thread
								from posts 
								where path[1] IN (SELECT id from posts where thread_id = $1 and parent = 0 %s %s %s;`, sinceQuery, innerOrderByQuery, outerOrderByQuery)

	rows, err := t.db.Query(query, values...)

	threadErr = &forum_errors.ThreadError{Code: forum_errors.Unknown}
	if err != nil {
		threadErr.Reason = err
		return posts, threadErr
	}
	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Forum, &post.Message, &post.IsEdited, &post.Created, &post.Thread)
		if err != nil {
			threadErr.Reason = err
			return posts, threadErr
		}
		posts = append(posts, post)

	}
	//err = scan.Rows(&posts, rows)
	//if err != nil {
	//	threadErr.Reason = err
	//	return posts, threadErr
	//}
	return posts, nil
}

func (t *ThreadRepo) GetPostsByThreadIdTree(threadId int64, limit int, since int, desc bool) (posts models.Posts, threadErr *forum_errors.ThreadError) {
	posts = make(models.Posts, 0)
	values := []interface{}{threadId}
	sinceQuery := ""
	if since >= 0 {
		sinceQuery = "and path "
		if desc {
			sinceQuery += "< (SELECT path from posts where id = $2) "
		} else {
			sinceQuery += "> (SELECT path from posts where id = $2) "
		}
		values = append(values, since)
	}

	orderByQuery := "order by path"
	if desc {
		orderByQuery += " DESC, id "
	} else {
		orderByQuery += ", id ASC "
	}
	if limit > 0 {
		orderByQuery += fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`select id, parent, author, forum, message, is_edited, created, thread_id as thread
								from posts where thread_id = $1 %s %s;`, sinceQuery, orderByQuery)

	rows, err := t.db.Query(query, values...)

	threadErr = &forum_errors.ThreadError{Code: forum_errors.Unknown}
	if err != nil {
		threadErr.Reason = err
		return posts, threadErr
	}
	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Forum, &post.Message, &post.IsEdited, &post.Created, &post.Thread)
		if err != nil {
			threadErr.Reason = err
			return posts, threadErr
		}
		posts = append(posts, post)

	}
	//err = scan.Rows(&posts, rows)
	//if err != nil {
	//	threadErr.Reason = err
	//	return posts, threadErr
	//}
	return posts, nil

}

func (t *ThreadRepo) GetPostsByThreadIdFlat(threadId int64, limit int, since int, desc bool) (posts models.Posts, threadErr *forum_errors.ThreadError) {
	posts = make(models.Posts, 0)
	values := []interface{}{threadId}
	sinceQuery := ""
	if since >= 0 {
		sinceQuery = "and id "
		if desc {
			sinceQuery += "< $2"
		} else {
			sinceQuery += "> $2"
		}
		values = append(values, since)
	}

	orderByQuery := "order by created "
	if desc {
		orderByQuery += "DESC, id DESC "
	} else {
		orderByQuery += ", id "
	}
	if limit > 0 {
		orderByQuery += fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`select id, parent, author, forum, message, is_edited, created, thread_id as thread
								from posts where thread_id = $1 %s %s;`, sinceQuery, orderByQuery)
	rows, err := t.db.Query(query, values...)
	threadErr = &forum_errors.ThreadError{Code: forum_errors.Unknown}
	if err != nil {
		threadErr.Reason = err
		return posts, threadErr
	}
	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Forum, &post.Message, &post.IsEdited, &post.Created, &post.Thread)
		if err != nil {
			threadErr.Reason = err
			return posts, threadErr
		}
		posts = append(posts, post)

	}
	//err = scan.Rows(&posts, rows)
	//if err != nil {
	//	threadErr.Reason = err
	//	return posts, threadErr
	//}
	return posts, nil
}

func (t *ThreadRepo) VoteForThread(threadID int64, vote models.Vote) (thread models.Thread, threadErr *forum_errors.ThreadError) {
	threadErr = &forum_errors.ThreadError{
		Code: forum_errors.Unknown,
	}
	_, err := t.db.Exec(`insert into votes (user_nickname, thread_id, voice) VALUES ($1, $2, $3) 
							ON CONFLICT (user_nickname, thread_id) DO UPDATE SET voice = EXCLUDED.voice;`,
		vote.Nickname, threadID, vote.Voice)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case postgres.FOREIGN_KEY_VIOLATION:
				threadErr.Code = forum_errors.CantFindUser
			}
		}
		threadErr.Reason = err
		return thread, threadErr
	}
	return t.GetThreadById(threadID)
}

func (t *ThreadRepo) GetThreadById(id int64) (thread models.Thread, forumError *forum_errors.ThreadError) {
	forumError = &forum_errors.ThreadError{Code: forum_errors.Unknown}
	queryResult, err := t.db.Query("select * from threads where id = $1", id)

	if err != nil {
		forumError.Reason = err
		return thread, forumError
	}

	err = scan.Row(&thread, queryResult)

	if err != nil {
		if err == sql.ErrNoRows {
			forumError.Code = forum_errors.CantFindThread
		}
		forumError.Reason = err
		return thread, forumError
	}
	return thread, nil

}

func (t *ThreadRepo) GetThreadBySlugOrId(slugOrId string) (thread models.Thread, forumError *forum_errors.ThreadError) {
	forumError = &forum_errors.ThreadError{Code: forum_errors.Unknown}
	id, err := strconv.ParseInt(slugOrId, 10, 64)

	if err == nil {
		return t.GetThreadById(id)
	}
	queryResult, err := t.db.Query("select * from threads where slug = $1", slugOrId)

	if err != nil {
		forumError.Reason = err
		return thread, forumError
	}

	err = scan.Row(&thread, queryResult)

	if err != nil {
		forumError.Reason = err
		if err == sql.ErrNoRows {
			forumError.Code = forum_errors.CantFindThread
		}
		return thread, forumError
	}
	return thread, nil

}

func (t *ThreadRepo) CreatePosts(posts models.Posts, threadID int64, forumSlug string, createdTime string) (createdPosts models.Posts, threadError *forum_errors.ThreadError) {
	threadError = &forum_errors.ThreadError{Code: forum_errors.Unknown}
	query := `insert into posts (parent, author, message, forum, created, thread_id) VALUES `

	values := []interface{}{}
	for i, post := range posts {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		values = append(values, post.Parent, post.Author, post.Message, forumSlug, createdTime, threadID)
	}
	query = query[:len(query)-1]

	query += " RETURNING id, parent, author, message, is_edited, forum, thread_id as thread, created;"

	queryResult, err := t.db.Query(query, values...)
	if err != nil {
		threadError.Reason = err
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case postgres.FOREIGN_KEY_VIOLATION:
				threadError.Code = forum_errors.CantFindUser
			}
		}
		return createdPosts, threadError
	}
	err = scan.Rows(&createdPosts, queryResult)

	if err != nil {
		threadError.Reason = err
		return createdPosts, threadError
	}
	return createdPosts, nil
}

func NewThreadRepo(dbConn *sql.DB) (repo *ThreadRepo, err error) {
	if dbConn == nil {
		return repo, errors.New("can't create repo without connection")
	}
	return &ThreadRepo{db: dbConn}, nil
}
