package handler

import (
	"github.com/gin-gonic/gin"

	"roulette/internal/gift/service"
	tonService "roulette/internal/ton/service"
)

type Handler struct {
	service    service.Service
	tonService tonService.Service
}

func NewHandler(service service.Service, tonService tonService.Service) *Handler {
	return &Handler{service: service, tonService: tonService}
}

func Router(h *Handler, r *gin.Engine) {
	router := r.Group("gift")
	{
		router.GET("/collections", h.getCollections)
		router.GET("/nft/:address", h.getNfts)
		router.GET("/user/:user_id", h.getUserGifts)
	}
}
