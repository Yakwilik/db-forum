package handler

import (
	"encoding/json"
	"fmt"
	"github.com/db-forum.git/pkg/forum_errors"
	"github.com/db-forum.git/pkg/models"
	"github.com/db-forum.git/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) GetUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	nickname, ok := vars[Nickname]
	if !ok {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": "no user nickname provided"})
	}
	user, err := h.services.GetUser(nickname)

	if err != nil {
		if err.Code == forum_errors.CantFindUser {
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": fmt.Sprintf("Can't find user by nickname: %s", nickname),
			})
			return
		}
		utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{
			"message": err.Reason,
		})
		return
	}

	utils.JSONResponse(writer, http.StatusOK, user)
}

func (h *Handler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	nickname, ok := vars[Nickname]
	if !ok {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": "no user nickname provided"})
	}
	updatableUser := models.User{
		Nickname: nickname,
	}
	err := json.NewDecoder(request.Body).Decode(&updatableUser)
	if err != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": err})
		return
	}
	updatedUser, userErr := h.services.UpdateUser(updatableUser)
	if userErr != nil {
		switch userErr.Code {
		case forum_errors.CantFindUser:
			utils.JSONResponse(writer, http.StatusNotFound, utils.InterfaceMap{
				"message": fmt.Sprintf("Can't find user by nickname: %s", nickname),
			})
		case forum_errors.ConflictingData:
			utils.JSONResponse(writer, http.StatusConflict, utils.InterfaceMap{
				"message": fmt.Sprintf("conflicting data: %s", nickname),
			})

		default:
			utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": userErr})

		}
		return

	}
	utils.JSONResponse(writer, http.StatusOK, updatedUser)
}

func (h *Handler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	nickname, ok := vars[Nickname]
	if !ok {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": "no user nickname provided"})
	}
	newUser := models.User{
		Nickname: nickname,
	}

	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		utils.JSONResponse(writer, http.StatusBadRequest, utils.InterfaceMap{"message": err})
		return
	}
	user, userErr := h.services.CreateUser(newUser)
	if userErr != nil {
		switch userErr.Code {
		case forum_errors.UserAlreadyExists:
			h.sendExistingUsers(newUser, writer)
			return
		}
		utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{
			"message": userErr.Reason,
		})
	}

	utils.JSONResponse(writer, http.StatusCreated, user)
}

func (h *Handler) sendExistingUsers(user models.User, writer http.ResponseWriter) {
	users, err := h.services.GetExistingUsers(user)
	if err != nil {
		utils.JSONResponse(writer, http.StatusInternalServerError, utils.InterfaceMap{
			"message": fmt.Sprintf("can`t get existing users with nickname %s or email %s becaust: %+v", user.Nickname, user.Email, err),
		})
		return
	}
	utils.JSONResponse(writer, http.StatusConflict, users)
}
