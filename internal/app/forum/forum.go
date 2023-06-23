package forum

import (
	"github.com/db-forum.git/internal/app/server"
	"github.com/db-forum.git/internal/handler"
	"github.com/db-forum.git/internal/repository/postgres"
	"github.com/db-forum.git/internal/repository/postgres/forum_repo"
	"github.com/db-forum.git/internal/repository/postgres/thread_repo"
	"github.com/db-forum.git/internal/repository/postgres/user_repo"
	"github.com/db-forum.git/internal/services"
	"github.com/db-forum.git/pkg/repository"
	"log"
)

type Forum struct {
	serv *server.Server
}

func NewForum() (forum *Forum, err error) {
	dbConfig := postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "mysecretpassword",
		DBName:   "postgres",
	}
	dbConn, err := postgres.NewPostgresDB(dbConfig)
	if err != nil {
		log.Println(err)
		return forum, err
	}
	userRepo, _ := user_repo.NewUserRepo(dbConn)
	forumRepo, _ := forum_repo.NewForumRepo(dbConn)
	threadRepo, _ := thread_repo.NewThreadRepo(dbConn)
	repo := &repository.Repository{User: userRepo, Forum: forumRepo, Thread: threadRepo}
	userService := services.NewUserService(repo)
	forumService := services.NewForumService(repo)
	threadService := services.NewThreadService(repo)
	service := services.New(userService, forumService, threadService)
	h, err := handler.NewHandler(service, "/api")

	if err != nil {
		return forum, err
	}
	serv := server.NewServer(server.Config{
		Host:           "127.0.0.1",
		Port:           5000,
		MaxHeaderBytes: 5000,
	}, h.InitRoutes())

	return &Forum{serv: serv}, err
}

func (f *Forum) StartApp() error {
	log.Println("server starts")
	return f.serv.ListenAndServe()
}
