package handler

import (
	"encoding/json"
	"fmt"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) CreatePosts(writer http.ResponseWriter, request *http.Request) {
	posts := models.Posts{}

	err := json.NewDecoder(request.Body).Decode(&posts)

	if err != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{
			"message": err.Error(),
		})
	}
	vars := mux.Vars(request)
	slugOrId, _ := vars[SlugOrId]
	if len(posts) == 0 {
		utils.JSONResponse(writer, http.StatusCreated, make(models.Posts, 0))
		return
	}
	newPosts, threadErr := h.services.CreatePosts(posts, slugOrId)
	if threadErr != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{
			"message": threadErr.Error(),
		})
		return
	}
	utils.JSONResponse(writer, http.StatusCreated, newPosts)

}

func (h *Handler) Vote(writer http.ResponseWriter, request *http.Request) {
	vote := models.Vote{}
	err := json.NewDecoder(request.Body).Decode(&vote)

	if err != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{
			"message": err.Error(),
		})
	}

	vars := mux.Vars(request)
	slugOrId, _ := vars[SlugOrId]

	thread, threadErr := h.services.Vote(slugOrId, vote)
	if threadErr != nil {
		switch threadErr.Code {
		case forum_errors.CantFindThread:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": "can`t find thread with slug_or_id:" + slugOrId,
			})
		}
		return
	}
	utils.JSONResponse(writer, http.StatusOK, thread)
}

func (h *Handler) GetThread(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slugOrId, _ := vars[SlugOrId]

	thread, threadErr := h.services.GetThreadBySlugOrId(slugOrId)
	if threadErr != nil {
		switch threadErr.Code {
		case forum_errors.CantFindThread:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": fmt.Sprintf("can't find thread with slug_or_id: %s", slugOrId),
			})
		}
		return
	}

	utils.JSONResponse(writer, http.StatusOK, thread)
}

func (h *Handler) GetPosts(writer http.ResponseWriter, request *http.Request) {
	params := getGetPostsRequestQueryParams(request)
	posts, threadErr := h.services.GetPosts(params.SlugOrId, params.Limit, params.Sort, params.Desc, params.Since)
	if threadErr != nil {
		fmt.Printf("ERROR : %+v\n", threadErr)
		utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{
			"message": threadErr,
		})
		return
	}

	utils.JSONResponse(writer, http.StatusOK, posts)
}

type GetPostsQueryParams struct {
	SlugOrId string
	Limit    int
	Since    int
	Desc     bool
	Sort     string
}

func getGetPostsRequestQueryParams(request *http.Request) (params GetPostsQueryParams) {
	vars := mux.Vars(request)
	slug, _ := vars[SlugOrId]
	params.SlugOrId = slug

	limit := request.FormValue("limit")
	if limit != "" {
		params.Limit, _ = strconv.Atoi(limit)
	} else {
		params.Limit = 100
	}

	since := request.FormValue("since")
	if since != "" {
		params.Since, _ = strconv.Atoi(since)
	} else {
		params.Since = -1
	}
	isDesc := request.FormValue("desc")
	if isDesc == "true" {
		params.Desc = true
	}
	params.Sort = request.FormValue("sort")
	return params
}
