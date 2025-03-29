package handler

import (
	"github.com/gin-gonic/gin"

	"roulette/internal/game/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func Router(h *Handler, r *gin.Engine) {
	router := r.Group("game")
	{
		router.GET("/round", h.getCurrentRound)
		//router.GET("/round/:round_id", h.getRound)
		router.GET("/winner/:round_id", h.getWinner)

		router.POST("/nft", h.addUserNft)
		router.POST("/gift", h.addUserGift)

		router.PUT("/fee/:user_id", h.updateFee)
	}
}
