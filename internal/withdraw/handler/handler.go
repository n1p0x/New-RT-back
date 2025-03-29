package handler

import (
	"github.com/gin-gonic/gin"

	userService "roulette/internal/user/service"
	"roulette/internal/withdraw/service"
)

type Handler struct {
	service     service.Service
	userService userService.Service
}

func NewHandler(service service.Service, userService userService.Service) *Handler {
	return &Handler{service: service, userService: userService}
}

func Router(h *Handler, r *gin.Engine) {
	router := r.Group("withdraw")
	{
		router.POST("/ton", h.addTon)
		router.POST("/nft", h.addNft)
		router.POST("/gift", h.addGift)
	}
}
