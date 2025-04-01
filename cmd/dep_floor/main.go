package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"roulette/internal/config"
	"roulette/internal/database"
	depositRepo "roulette/internal/deposit/repo"
	depositService "roulette/internal/deposit/service"
	giftRepo "roulette/internal/gift/repo"
	giftService "roulette/internal/gift/service"
	tgService "roulette/internal/tg/service"
	tonService "roulette/internal/ton/service"
	userRepo "roulette/internal/user/repo"
	userService "roulette/internal/user/service"
)

const (
	configPath = "config/prod.yaml"
	envPath    = ".env"
)

func main() {
	conf, err := config.Load(configPath, envPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewDatabase(conf)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		runDeposit(conf, db)
	}()

	go func() {
		defer wg.Done()
		runFloor(conf, db)
	}()

	wg.Wait()
}

func runDeposit(conf *config.Config, db *gorm.DB) {
	serviceTon := tonService.NewService(conf)

	repoUser := userRepo.NewRepo(db)
	serviceUser := userService.NewService(repoUser)

	repoGift := giftRepo.NewRepo(db)
	serviceGift := giftService.NewService(repoGift)

	repo := depositRepo.NewRepo(db)
	service := depositService.NewService(repo, serviceTon, serviceUser, serviceGift)

	errCh := make(chan error, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if errTon := service.CheckTonDeposit(ctx); errTon != nil {
			errCh <- errTon
		}
	}()

	go func() {
		defer wg.Done()
		if errNft := service.CheckNftDeposit(ctx); errNft != nil {
			errCh <- errNft
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		log.Fatalf("timeout: %v", ctx.Err())
	case err := <-errCh:
		if err != nil {
			log.Fatalf("failed to get deposits: %v", err)
		}
		log.Println("checking deposits completed")
	}
}

func runFloor(conf *config.Config, db *gorm.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

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
