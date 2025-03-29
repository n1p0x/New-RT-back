package handler

import (
	"github.com/gin-gonic/gin"

	"roulette/internal/deposit/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func Router(h *Handler, r *gin.Engine) {
	router := r.Group("deposit")
	{
		router.POST("/nft", h.addNftDeposit)
	}
}
