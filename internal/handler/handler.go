package handler

import (
	"errors"
	"fmt"
	"github.com/db-forum.git/internal/middleware"
	"github.com/db-forum.git/pkg/services"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	Nickname        = "NICKNAME"
	Slug            = "SLUG"
	SlugOrId        = "SLUG_OR_ID"
	ApplicationJSON = "application/json"
)

type Handler struct {
	services *services.Services
	mux      *mux.Router
}

func NewHandler(services *services.Services, apiRoot string) (*Handler, error) {
	if services == nil {
		return nil, errors.New("nil services")
	}
	return &Handler{
		services: services,
		mux:      mux.NewRouter().PathPrefix(apiRoot).Subrouter(),
	}, nil
}

func (h *Handler) InitRoutes() *mux.Router {
	h.mux.Use(func(handler http.Handler) http.Handler {
		return middleware.ContentTypeMiddleware(ApplicationJSON, handler)
	})
	userRouter := h.mux.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc(fmt.Sprintf("/{%s}/create", Nickname), h.CreateUser).Methods("POST")
	userRouter.HandleFunc(fmt.Sprintf("/{%s}/profile", Nickname), h.UpdateUser).Methods("POST")
	userRouter.HandleFunc(fmt.Sprintf("/{%s}/profile", Nickname), h.GetUser).Methods("GET")

	forumRouter := h.mux.PathPrefix("/forum").Subrouter()
	forumRouter.HandleFunc("/create", h.CreateForum).Methods("POST")
	forumRouter.HandleFunc(fmt.Sprintf("/{%s}/details", Slug), h.GetForum).Methods("GET")
	forumRouter.HandleFunc(fmt.Sprintf("/{%s}/create", Slug), h.CreateThreadInForum).Methods("POST")
	forumRouter.HandleFunc(fmt.Sprintf("/{%s}/threads", Slug), h.GetThreadsOfForum).Methods("GET")
	forumRouter.HandleFunc(fmt.Sprintf("/{%s}/users", Slug), h.GetUsersOfForum).Methods("GET")

	threadRouter := h.mux.PathPrefix("/thread").Subrouter()
	threadRouter.HandleFunc(fmt.Sprintf("/{%s}/create", SlugOrId), h.CreatePosts).Methods("POST")
	threadRouter.HandleFunc(fmt.Sprintf("/{%s}/vote", SlugOrId), h.Vote).Methods("POST")
	threadRouter.HandleFunc(fmt.Sprintf("/{%s}/details", SlugOrId), h.GetThread).Methods("GET")
	threadRouter.HandleFunc(fmt.Sprintf("/{%s}/posts", SlugOrId), h.GetPosts).Methods("GET")

	return h.mux
}
