package main

import (
	"context"
	"log"
	"sync"
	"time"

	"roulette/internal/config"
	"roulette/internal/database"
	depositRepo "roulette/internal/deposit/repo"
	depositService "roulette/internal/deposit/service"
	giftRepo "roulette/internal/gift/repo"
	giftService "roulette/internal/gift/service"
	tonService "roulette/internal/ton/service"
	userRepo "roulette/internal/user/repo"
	userService "roulette/internal/user/service"
)

const (
	configPath = "config/prod.yaml"
	envPath    = ".env"
)

func main() {
	runDeposit()
}

func runDeposit() {
	conf, err := config.Load(configPath, envPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewDatabase(conf)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

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
