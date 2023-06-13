package user_repo

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

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(dbConn *sql.DB) (repo *UserRepo, err error) {
	if dbConn == nil {
		return repo, errors.New("can't create repo without connection")
	}
	return &UserRepo{db: dbConn}, nil
}

func (u *UserRepo) CreateUser(user models.User) (createdUser models.User, userErr *forum_errors.UserError) {
	_, err := u.db.Exec(`insert into users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);`, user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case postgres.UNIQUE_VIOLATION:
				return createdUser, &forum_errors.UserError{
					Reason: pqErr,
					Code:   forum_errors.UserAlreadyExists,
				}
			}
		}
		return models.User{}, &forum_errors.UserError{
			Reason: err,
			Code:   forum_errors.Unknown,
		}
	}
	return user, nil
}

func (u *UserRepo) GetExistingUsers(user models.User) (users models.Users, err error) {
	query, err := u.db.Query(`select nickname, fullname, about, email from users where (nickname) = $1 or email = $2;`, user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}
	err = scan.Rows(&users, query)
	return users, err

}

func (u *UserRepo) GetUser(nickname string) (user models.User, userErr *forum_errors.UserError) {
	userQuery, err := u.db.Query(`select nickname, fullname, about, email from users where nickname = $1 limit 1;`, nickname)
	userErr = &forum_errors.UserError{}
	if err != nil {
		userErr.Code = forum_errors.Unknown
		userErr.Reason = err
		return user, userErr
	}
	err = scan.Row(&user, userQuery)
	if err != nil {
		userErr.Reason = err
		if err == sql.ErrNoRows {
			userErr.Code = forum_errors.CantFindUser
		}
		return user, userErr
	}
	return user, nil
}

func (u *UserRepo) UpdateUser(user models.User) (updatedUser models.User, userErr *forum_errors.UserError) {
	userErr = &forum_errors.UserError{Code: forum_errors.Unknown}
	userQuery, err := u.db.Query(`update users set 
	                 fullname = $1, 
	                 about = $2, 
	                 email = $3
	                 where nickname = $4 
	                                returning fullname, about, email, nickname`,
		user.Fullname, user.About, user.Email, user.Nickname)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Printf("pq Error: %v", pqErr)
			switch pqErr.Code.Name() {
			case postgres.UNIQUE_VIOLATION:
				userErr.Code = forum_errors.ConflictingData
				userErr.Reason = pqErr
				return updatedUser, userErr
			}
		}
		userErr.Reason = err
		return models.User{}, userErr
	}
	err = scan.Row(&updatedUser, userQuery)
	if err != nil {
		userErr.Reason = err
		return models.User{}, userErr
	}
	log.Printf("updatedUser: %v", updatedUser)
	return updatedUser, nil
}
