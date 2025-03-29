package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"roulette/internal/config"
	"roulette/internal/database"
	giftRepo "roulette/internal/gift/repo"
	giftService "roulette/internal/gift/service"
	tgService "roulette/internal/tg/service"
)

var configPath = "config/local.yaml"

func runFloor() {
	conf, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewDatabase(conf)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	repoNft := giftRepo.NewRepo(db)
	serviceNft := giftService.NewService(repoNft)

	service := tgService.NewService(conf)

	floors, err := service.GetFloors(ctx, 2422226195, -8093627540162659735)
	if err != nil {
		fmt.Println(err)
	}

	if err = serviceNft.UpdateCollectionsFloor(ctx, floors); err != nil {
		log.Printf("failed to update floors: %v", err)
	}
}
