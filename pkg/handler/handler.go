package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type EndpointRegistrator interface {
	RegisterEndpoints(router *mux.Router)
}

type UserHandler interface {
	CreateUser(writer http.ResponseWriter, request *http.Request)
	GetUser(writer http.ResponseWriter, request *http.Request)
	UpdateUser(writer http.ResponseWriter, request *http.Request)
	EndpointRegistrator
}

type ThreadHandler interface {
	CreateThread(writer http.ResponseWriter, request *http.Request)
	GetThread(writer http.ResponseWriter, request *http.Request)
	UpdateThread(writer http.ResponseWriter, request *http.Request)
	GetPosts(writer http.ResponseWriter, request *http.Request)
	VoteThread(writer http.ResponseWriter, request *http.Request)
	EndpointRegistrator
}

type ServiceHandler interface {
	ClearService(writer http.ResponseWriter, request *http.Request)
	GetStatus(writer http.ResponseWriter, request *http.Request)
	EndpointRegistrator
}

type PostHandler interface {
	GetById(writer http.ResponseWriter, request *http.Request)
	UpdateById(writer http.ResponseWriter, request *http.Request)
	EndpointRegistrator
}

type ForumHandler interface {
	CreateForum(writer http.ResponseWriter, request *http.Request)
	GetForum(writer http.ResponseWriter, request *http.Request)
	CreateThreadInForum(writer http.ResponseWriter, request *http.Request)
	GetUserOfForum(writer http.ResponseWriter, request *http.Request)
	GetThreadsOfForum(writer http.ResponseWriter, request *http.Request)
	EndpointRegistrator
}

type Handler struct {
	UserHandler
}
