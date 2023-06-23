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
	"strings"
	"time"
)

func (h *Handler) CreatePosts(writer http.ResponseWriter, request *http.Request) {

	start := time.Now()

	defer func() {
		fmt.Printf("CreatePosts Function execution took %s\n", time.Since(start))
	}()
	posts := models.Posts{}

	err := json.NewDecoder(request.Body).Decode(&posts)

	if err != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{
			"message": err.Error(),
		})
	}
	vars := mux.Vars(request)
	slugOrId, _ := vars[SlugOrId]
	newPosts, threadErr := h.services.CreatePosts(posts, slugOrId)
	if threadErr != nil {
		switch threadErr.Code {
		case forum_errors.CantFindPost:
			fallthrough
		case forum_errors.ConflictingData:
			utils.JSONResponse(writer, http.StatusConflict, utils.InterfaceMap{
				"message": "Parent post was created in another thread",
			})
		case forum_errors.CantFindUser:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": "Can't find author by nickname",
			})
		case forum_errors.CantFindThread:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": "Can't find thread by id",
			})
		default:
			utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{
				"message": threadErr.Error(),
			})
		}

		return
	}
	utils.JSONResponse(writer, http.StatusCreated, newPosts)

}

func (h *Handler) Vote(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()

	defer func() {
		fmt.Printf("Vote Function execution took %s\n", time.Since(start))
	}()
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
		case forum_errors.CantFindUser:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": "can`t find user:" + slugOrId,
			})
		default:
			utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{
				"message": threadErr.Error(),
			})
		}
		return
	}
	utils.JSONResponse(writer, http.StatusOK, thread)
}

type updateThreadData struct {
	SlugOrId string
	models.ThreadUpdate
}

func getUpdateThreadData(request *http.Request) (data updateThreadData) {
	vars := mux.Vars(request)
	data.SlugOrId, _ = vars[SlugOrId]

	_ = json.NewDecoder(request.Body).Decode(&data)

	return data
}

func (h *Handler) UpdateThread(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()

	defer func() {
		fmt.Printf("UpdateThread Function execution took %s\n", time.Since(start))
	}()
	updateData := getUpdateThreadData(request)

	thread, threadErr := h.services.UpdateThread(updateData.SlugOrId, updateData.ThreadUpdate)
	if threadErr != nil {
		switch threadErr.Code {
		default:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": fmt.Sprintf("this error: %s", threadErr.Error()),
			})
		}
		return
	}

	utils.JSONResponse(writer, http.StatusOK, thread)

}

type getPostParams struct {
	Id            int
	RelatedUser   bool
	RelatedThread bool
	RelatedForum  bool
}

func getGetPostParams(request *http.Request) (params getPostParams) {
	vars := mux.Vars(request)
	idString, _ := vars[Id]

	params.Id, _ = strconv.Atoi(idString)
	related := request.FormValue("related")
	relatedSlice := strings.Split(related, ",")
	for _, related := range relatedSlice {
		switch related {
		case "forum":
			params.RelatedForum = true
		case "user":
			params.RelatedUser = true
		case "thread":
			params.RelatedThread = true
		}
	}

	return params
}
func (h *Handler) GetPost(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()

	defer func() {
		fmt.Printf("GetPost Function execution took %s\n", time.Since(start))
	}()
	params := getGetPostParams(request)

	fullPost, threadErr := h.services.GetPost(params.Id, params.RelatedUser, params.RelatedThread, params.RelatedForum)

	if threadErr != nil {
		switch threadErr.Code {
		default:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{"message": "can't find" + threadErr.Error()})
		}
		return
	}

	utils.JSONResponse(writer, http.StatusOK, fullPost)
}

type updatePostParams struct {
	Id int
	models.PostUpdate
}

func getUpdatePostParams(request *http.Request) (params updatePostParams) {
	vars := mux.Vars(request)
	idString, _ := vars[Id]

	params.Id, _ = strconv.Atoi(idString)

	_ = json.NewDecoder(request.Body).Decode(&params)

	return params
}

func (h *Handler) UpdatePost(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()

	defer func() {
		fmt.Printf("UpdatePost Function execution took %s\n", time.Since(start))
	}()
	params := getUpdatePostParams(request)

	post, threadErr := h.services.UpdatePost(params.Id, params.PostUpdate)

	if threadErr != nil {
		switch threadErr.Code {
		default:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{"message": "can't find" + threadErr.Error()})
		}
		return
	}
	utils.JSONResponse(writer, http.StatusOK, post)
}

func (h *Handler) GetThread(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()

	defer func() {
		fmt.Printf("GetThread Function execution took %s\n", time.Since(start))
	}()
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
	start := time.Now()

	params := getGetPostsRequestQueryParams(request)

	defer func() {
		fmt.Printf("GetPosts Function with sort = %s execution took %s\n", params.Sort, time.Since(start))
	}()
	posts, threadErr := h.services.GetPosts(params.SlugOrId, params.Limit, params.Sort, params.Desc, params.Since)
	if threadErr != nil {
		switch threadErr.Code {
		case forum_errors.CantFindThread:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": fmt.Sprintf("can't find thread with slug_or_id: %s", params.SlugOrId),
			})
		default:
			utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{
				"message": threadErr,
			})
		}
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
