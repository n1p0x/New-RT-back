package handler

import (
	"github.com/gin-gonic/gin"

	"roulette/internal/user/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func Router(r *gin.Engine, h *Handler) {
	router := r.Group("user")
	{
		router.GET("/info/:user_id", h.getUser)
		router.GET("/profile/:user_id", h.getUserProfile)

		router.POST("/new", h.addUser)
	}
}
