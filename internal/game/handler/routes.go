package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"roulette/internal/database"
	"roulette/internal/game/service"
	"roulette/internal/middleware/handler"
)

func (h *Handler) getCurrentRound(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		round, err := h.service.GetCurrentRoundWithPlayers(c.Request.Context())
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, NewRoundResponse(round))
	})
}

func (h *Handler) getRound(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestUri struct {
			RoundID uint `uri:"round_id"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		round, err := h.service.GetRoundWithPlayers(c.Request.Context(), uri.RoundID)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, NewRoundResponse(round))
	})
}

func (h *Handler) getWinner(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestUri struct {
			RoundID uint `uri:"round_id"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		winner, err := h.service.GetWinner(c.Request.Context(), uri.RoundID)
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, NewWinnerResponse(winner))
	})
}

func (h *Handler) addUserNft(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestBody struct {
			UserNftID uint `json:"userNftId"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		if err := h.service.AddUserNft(c.Request.Context(), body.UserNftID); err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			if service.IsRoundNotFound(err) {
				return handler.NewInternalErrorResponse(err)
			}
			if service.IsRoundFinished(err) {
				return handler.NewInternalErrorResponse(err)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusCreated, nil)
	})
}

func (h *Handler) addUserGift(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestBody struct {
			UserGiftID uint `json:"userGiftId"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		if err := h.service.AddUserGift(c.Request.Context(), body.UserGiftID); err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, err.Error())
			}
			if service.IsRoundNotFound(err) {
				return handler.NewInternalErrorResponse(err)
			}
			if service.IsRoundFinished(err) {
				return handler.NewInternalErrorResponse(err)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusCreated, nil)
	})
}

func (h *Handler) updateFee(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestUri struct {
			UserID uint `uri:"user_id"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		if err := h.service.UpdateFee(c.Request.Context(), uri.UserID); err != nil {
			if service.IsNotEnoughBalance(err) {
				return handler.NewErrorResponse(http.StatusBadRequest, fmt.Sprintf("user %d has insufficient balance", uri.UserID))
			}
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, fmt.Sprintf("user %d not found", uri.UserID))
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusNoContent, nil)
	})
}
