package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// TimeoutMiddleware attaches deadline to gin.Request.Context
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				c.AbortWithStatus(http.StatusGatewayTimeout)
			}
			cancel()
		}()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// AuthMiddleware validates tg init data
func AuthMiddleware(token string, mode string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if mode == gin.DebugMode {
			return
		}

		auth := strings.Split(ctx.GetHeader("authorization"), " ")
		if len(auth) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"detail": "Unauthorized",
			})
			return
		}

		authType, authData := auth[0], auth[1]

		if authType == "Tg" {
			if err := initdata.Validate(authData, token, time.Hour); err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
					"detail": "Invalid init data",
				})
			}

			initData, err := initdata.Parse(authData)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
					"detail": "Something went wrong",
				})
			}

			ctx.Request = ctx.Request.WithContext(
				withInitData(ctx.Request.Context(), initData),
			)
		}
	}
}
