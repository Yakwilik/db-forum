package handler

import (
	"errors"
	"fmt"
	"github.com/db-forum.git/internal/middleware"
	"github.com/db-forum.git/pkg/services"
	"github.com/db-forum.git/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	Nickname        = "NICKNAME"
	Slug            = "SLUG"
	SlugOrId        = "SLUG_OR_ID"
	Id              = "ID"
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
	threadRouter.HandleFunc(fmt.Sprintf("/{%s}/details", SlugOrId), h.UpdateThread).Methods("POST")
	threadRouter.HandleFunc(fmt.Sprintf("/{%s}/posts", SlugOrId), h.GetPosts).Methods("GET")

	postRouter := h.mux.PathPrefix("/post").Subrouter()
	postRouter.HandleFunc(fmt.Sprintf("/{%s}/details", Id), h.GetPost).Methods("GET")
	postRouter.HandleFunc(fmt.Sprintf("/{%s}/details", Id), h.UpdatePost).Methods("POST")

	serviceRouter := h.mux.PathPrefix("/service").Subrouter()

	serviceRouter.HandleFunc("/status", h.GetServiceStatus).Methods("GET")
	serviceRouter.HandleFunc("/clear", h.Clear).Methods("POST")

	return h.mux
}
func (h *Handler) Clear(writer http.ResponseWriter, request *http.Request) {
	_ = h.services.Clear()

	writer.WriteHeader(http.StatusOK)
}
func (h *Handler) GetServiceStatus(writer http.ResponseWriter, request *http.Request) {
	status, err := h.services.GetServiceStatus()
	if err != nil {
		fmt.Printf("error: %+v", err)
	}

	utils.JSONResponse(writer, http.StatusOK, status)
}
