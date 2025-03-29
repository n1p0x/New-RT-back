package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"roulette/internal/database"
	"roulette/internal/middleware/handler"
	"roulette/internal/ton/service"
	"roulette/internal/utils"
)

func (h *Handler) getCollections(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		collections, err := h.service.GetCollections(c.Request.Context())
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, NewCollectionsResponse(collections))
	})
}

func (h *Handler) getNfts(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestUri struct {
			Address string `uri:"address"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		if _, err := utils.GetAddress(uri.Address); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		collections, err := h.service.GetCollections(c.Request.Context())
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		nfts, err := h.tonService.GetWalletNfts(c.Request.Context(), uri.Address, collections)
		if err != nil {
			if service.IsNftsNotFound(err) {
				return handler.NewErrorResponse(http.StatusNotFound, fmt.Sprintf("wallet %s nfts not found", uri.Address))
			}
			return handler.NewErrorResponse(http.StatusBadRequest, err.Error())
		}

		return handler.NewSuccessResponse(http.StatusOK, NewNftsResponse(nfts))
	})
}

func (h *Handler) getUserGifts(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		type RequestUri struct {
			UserID uint `uri:"user_id"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			return handler.NewUnprocessableErrorResponse(err)
		}

		userGifts, fee, err := h.service.GetUserGifts(c.Request.Context(), uri.UserID)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, fmt.Sprintf("user %d gifts not found", uri.UserID))
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, NewUserGiftsResponse(userGifts, fee == 0, fee))
	})
}
