package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"roulette/internal/database"
	"roulette/internal/middleware/handler"
)

func (h *Handler) addTon(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestBody struct {
			UserID      uint   `json:"userId"`
			Destination string `json:"destination"`
			Amount      int    `json:"amount"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		err := h.service.AddTon(c.Request.Context(), body.UserID, body.Destination, uint(body.Amount))
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusCreated, "")
	})
}

func (h *Handler) addNft(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestBody struct {
			UserNftID   uint   `json:"userNftId"`
			Destination string `json:"destination"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		err := h.service.AddNft(c.Request.Context(), body.UserNftID, body.Destination)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusCreated, nil)
	})
}

func (h *Handler) addGift(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestBody struct {
			UserGiftID uint `json:"userGiftId"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		err := h.service.AddGift(c.Request.Context(), body.UserGiftID)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusCreated, "")
	})
}
