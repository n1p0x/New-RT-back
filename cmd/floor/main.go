package main

import (
	"context"
	"log"
	"time"

	"roulette/internal/config"
	"roulette/internal/database"
	giftRepo "roulette/internal/gift/repo"
	giftService "roulette/internal/gift/service"
	tgService "roulette/internal/tg/service"
)

const (
	configPath = "config/prod.yaml"
	envPath    = ".env"
)

func main() {
	runFloor()
}

func runFloor() {
	log.Println("starting...", time.Now().Format("2006-01-02 15:04:05"))

	cfg, err := config.Load(configPath, envPath)
	if err != nil {
		log.Fatalf("failed to load cfgig: %v", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	repoGift := giftRepo.NewRepo(db)
	serviceGift := giftService.NewService(repoGift)

	service := tgService.NewService(cfg)

	floors, err := service.GetFloorsLow(ctx, 2422226195, -6696550584382502344)
	if err != nil {
		log.Fatalf("failed to get floors: %v", err)
	}

	if err = serviceGift.UpdateCollectionsFloor(ctx, floors); err != nil {
		log.Fatalf("failed to update floors: %v", err)
	}

	log.Println("finished", time.Now().Format("2006-01-02 15:04:05"))
}
