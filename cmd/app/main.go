package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"roulette/internal/config"
	"roulette/internal/database"
	depositHandler "roulette/internal/deposit/handler"
	depositRepo "roulette/internal/deposit/repo"
	depositService "roulette/internal/deposit/service"
	gameHandler "roulette/internal/game/handler"
	gameRepo "roulette/internal/game/repo"
	gameService "roulette/internal/game/service"
	giftHandler "roulette/internal/gift/handler"
	giftRepo "roulette/internal/gift/repo"
	giftService "roulette/internal/gift/service"
	"roulette/internal/middleware"
	tgService "roulette/internal/tg/service"
	tonService "roulette/internal/ton/service"
	userHandler "roulette/internal/user/handler"
	userRepo "roulette/internal/user/repo"
	userService "roulette/internal/user/service"
	withdrawHandler "roulette/internal/withdraw/handler"
	withdrawRepo "roulette/internal/withdraw/repo"
	withdrawService "roulette/internal/withdraw/service"
)

const (
	configPath = "config/prod.yaml"
	envPath    = ".env"
)

func main() {
	runApplication()
}

func runApplication() {
	cfg, err := config.Load(configPath, envPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app := fx.New(
		fx.Supply(cfg),
		fx.StopTimeout(cfg.ServerConfig.GracefulShutdown+time.Second),
		fx.Provide(
			database.NewDatabase,

			tonService.NewService,
			tgService.NewService,

			userRepo.NewRepo,
			userService.NewService,
			userHandler.NewHandler,

			giftRepo.NewRepo,
			giftService.NewService,
			giftHandler.NewHandler,

			depositRepo.NewRepo,
			depositService.NewService,
			depositHandler.NewHandler,

			withdrawRepo.NewRepo,
			withdrawService.NewService,
			withdrawHandler.NewHandler,

			gameRepo.NewRepo,
			gameService.NewService,
			gameHandler.NewHandler,

			newServer,
		),
		fx.Invoke(
			userHandler.Router,
			giftHandler.Router,
			depositHandler.Router,
			withdrawHandler.Router,
			gameHandler.Router,
			func(r *gin.Engine) {},
		),
	)

	app.Run()
}

func newServer(lc fx.Lifecycle, cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Mode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Origin},
		AllowMethods:     []string{"HEAD", "GET", "POST", "PUT"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(middleware.AuthMiddleware(cfg.TgConfig.BotToken, gin.DebugMode))
	r.Use(middleware.TimeoutMiddleware(cfg.ServerConfig.WriteTimeout))

	r.HEAD("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port),
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen and server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return r
}
