package handler

import (
	"encoding/json"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/utils"
	"net/http"
)

func (h *Handler) CreateForum(writer http.ResponseWriter, request *http.Request) {
	newForum := models.Forum{}
	err := json.NewDecoder(request.Body).Decode(&newForum)

	if err != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": err})
	}
	createdForum, forumErr := h.services.CreateForum(newForum)
	if forumErr != nil {
		switch forumErr.Code {
		case forum_errors.ForumAlreadyExists:
			h.sendExistingForum(writer, newForum.Slug)
		case forum_errors.CantFindUser:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": "can`t find such user",
			})
		}
		return
	}
	utils.JSONResponse(writer, http.StatusCreated, createdForum)
}
func (h *Handler) GetForum(writer http.ResponseWriter, request *http.Request) {

}
func (h *Handler) CreateThreadInForum(writer http.ResponseWriter, request *http.Request) {

}
func (h *Handler) GetUserOfForum(writer http.ResponseWriter, request *http.Request) {

}
func (h *Handler) GetThreadsOfForum(writer http.ResponseWriter, request *http.Request) {

}

func (h *Handler) sendExistingForum(writer http.ResponseWriter, slug string) {
	forum, forumErr := h.services.GetForumInfo(slug)
	if forumErr != nil {
		switch forumErr.Code {

		}

	}
	utils.JSONResponse(writer, http.StatusConflict, forum)
}
