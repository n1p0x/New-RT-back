package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"roulette/internal/middleware/handler"
	"roulette/internal/user/service"
)

func (h *Handler) getUser(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestUri struct {
			UserID uint `uri:"user_id"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		user, err := h.service.GetUser(c.Request.Context(), uri.UserID)
		if err != nil {
			if service.IsUserNotFound(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusOK, NewUserResponse(user))
	})
}

func (h *Handler) getUserProfile(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestUri struct {
			UserID uint `uri:"user_id"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		userProfile, err := h.service.GetUserProfile(c.Request.Context(), uri.UserID)
		if err != nil {
			return handler.NewErrorResponse(http.StatusNotFound, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusOK, NewUserProfileResponse(userProfile))
	})
}

func (h *Handler) addUser(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestBody struct {
			UserID     uint    `json:"userId"`
			Name       *string `json:"name"`
			PhotoUrl   *string `json:"photoUrl"`
			StartParam *string `json:"startParam"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		err := h.service.AddUser(c.Request.Context(), body.UserID, body.Name, body.PhotoUrl, body.StartParam)
		if err != nil {
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusCreated, nil)
	})
}
