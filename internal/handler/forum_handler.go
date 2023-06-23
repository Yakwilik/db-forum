package handler

import (
	"encoding/json"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	vars := mux.Vars(request)
	slug, ok := vars[Slug]
	if !ok {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": "no slug provided"})
	}
	forum, forumErr := h.services.GetForumInfo(slug)

	if forumErr != nil {
		switch forumErr.Code {
		case forum_errors.CantFindForum:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{"message": "can`t find forum"})
		default:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{"message": forumErr.Reason})
		}
		return
	}
	utils.JSONResponse(writer, http.StatusOK, forum)
}
func (h *Handler) CreateThreadInForum(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug, ok := vars[Slug]
	if !ok {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": "no slug provided"})
	}
	newThread := models.Thread{
		Forum: slug,
	}
	err := json.NewDecoder(request.Body).Decode(&newThread)

	if err != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": err})
	}
	newThread, forumErr := h.services.CreateThread(slug, newThread)
	if forumErr != nil {
		switch forumErr.Code {
		case forum_errors.CantFindForum:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{"message": "no forum with slug " + slug})
		case forum_errors.CantFindUser:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{"message": "no user with nickname " + newThread.Author})
		case forum_errors.ThreadAlreadyExists:
			utils.JSONResponse(writer, http.StatusConflict, newThread)
		default:
			utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{"message": forumErr.Reason})
		}
		return
	}
	utils.JSONResponse(writer, http.StatusCreated, newThread)
}
func (h *Handler) GetUsersOfForum(writer http.ResponseWriter, request *http.Request) {
	params := getRequestQueryParams(request)

	users, forumErr := h.services.GetForumUsers(params.Slug, params.Limit, params.Since, params.Desc)

	if forumErr != nil {
		switch forumErr.Code {
		case forum_errors.CantFindForum:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{"message": "can't find forum"})
		default:
			utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{"message": forumErr})
		}
		return
	}

	utils.JSONResponse(writer, http.StatusOK, users)
}
func (h *Handler) GetThreadsOfForum(writer http.ResponseWriter, request *http.Request) {
	params := getRequestQueryParams(request)

	threads, forumErr := h.services.GetForumThreads(params.Slug, params.Limit, params.Since, params.Desc)
	if forumErr != nil {
		utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
			"message": forumErr.Reason.Error(),
		})
		return
	}
	utils.JSONResponse(writer, http.StatusOK, threads)

}

func (h *Handler) sendExistingForum(writer http.ResponseWriter, slug string) {
	forum, forumErr := h.services.GetForumInfo(slug)
	if forumErr != nil {
		switch forumErr.Code {

		}

	}
	utils.JSONResponse(writer, http.StatusConflict, forum)
}

type QueryParams struct {
	Slug  string
	Limit int
	Since string
	Desc  bool
	Sort  string
}

func getRequestQueryParams(request *http.Request) (params QueryParams) {
	vars := mux.Vars(request)
	slug, _ := vars[Slug]
	params.Slug = slug

	limit := request.FormValue("limit")
	if limit != "" {
		params.Limit, _ = strconv.Atoi(limit)
	} else {
		params.Limit = 100
	}

	params.Since = request.FormValue("since")
	isDesc := request.FormValue("desc")
	if isDesc == "true" {
		params.Desc = true
	}
	params.Sort = request.FormValue("sort")
	return params
}
