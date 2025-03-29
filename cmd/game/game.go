package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"roulette/internal/config"
	"roulette/internal/database"
	gameRepo "roulette/internal/game/repo"
	gameService "roulette/internal/game/service"
)

var configPath = "config/local.yaml"

func runGame() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}

	repo := gameRepo.NewRepo(db)
	service := gameService.NewService(repo)

	for {
		fmt.Println("checking...: ", time.Now())
		nextCheckTime := checkRound(context.Background(), service)

		delay := time.Until(nextCheckTime)
		if delay < 0 {
			delay = 0
		}

		select {
		case <-time.After(delay):
		}
	}
}

func checkRound(ctx context.Context, service gameService.Service) time.Time {
	roundWithPlayers, err := service.GetCurrentRoundWithPlayers(ctx)
	if err != nil {
		if gameService.IsRoundNotFound(err) {
			errAdd := service.AddRound(ctx)
			if errAdd != nil {
				log.Fatalf("failed to add round: %v", err)
			}

			roundWithPlayers, err = service.GetCurrentRoundWithPlayers(ctx)
			if err != nil {
				log.Fatalf("failed to get new round: %v", err)
			}
		} else {
			return time.Now().Add(5 * time.Second)
		}
	}
	if roundWithPlayers.IsFinished {
		_, err = service.GetWinner(ctx, roundWithPlayers.ID)
		if err != nil {
			if time.Now().Unix() >= roundWithPlayers.StartedAt.Add(RoundDuration).Unix() {
				if err = service.AddWinner(ctx, roundWithPlayers); err != nil {
					return time.Now().Add(1 * time.Second)
				}
			}
		}

		errAdd := service.AddRound(ctx)
		if errAdd != nil {
			log.Fatalf("failed to add round: %v", err)
		}

		roundWithPlayers, err = service.GetCurrentRoundWithPlayers(ctx)
		if err != nil {
			log.Fatalf("failed to get new round: %v", err)
		}
	}

	fmt.Println(*roundWithPlayers.Round)

	if len(roundWithPlayers.UniquePlayers) < MinPlayers {
		return time.Now().Add(3 * time.Second)
	}

	if roundWithPlayers.StartedAt == nil {
		if err = service.StartRound(ctx, roundWithPlayers.ID); err != nil {
			return time.Now().Add(1 * time.Second)
		}
		roundWithPlayers, _ = service.GetCurrentRoundWithPlayers(ctx)
		return roundWithPlayers.StartedAt.Add(RoundDuration)
	}

	// maybe not necessary
	return roundWithPlayers.StartedAt.Add(RoundDuration)
}
