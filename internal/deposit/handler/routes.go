package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"roulette/internal/middleware/handler"
)

func (h *Handler) addNftDeposit(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestBody struct {
			UserID     uint   `json:"userId"`
			Sender     string `json:"sender"`
			NftAddress string `json:"address"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		if err := h.service.AddNft(c.Request.Context(), body.UserID, body.Sender, body.NftAddress); err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusCreated, nil)
	})
}
