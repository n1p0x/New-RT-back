package service

import (
	"context"

	"roulette/internal/game/model"
	"roulette/internal/game/repo"
)

type Service interface {
	GetCurrentRound(ctx context.Context) (*model.Round, error)

	GetCurrentRoundWithPlayers(ctx context.Context) (*model.RoundWithPlayers, error)

	GetRoundWithPlayers(ctx context.Context, roundID uint) (*model.RoundWithPlayers, error)

	GetWinner(ctx context.Context, roundID uint) (*model.Winner, error)

	AddRound(ctx context.Context) error

	AddUserNft(ctx context.Context, userNftID uint) error

	AddUserGift(ctx context.Context, userGiftID uint) error

	AddWinner(ctx context.Context, roundWithPlayers *model.RoundWithPlayers) error

	UpdateFee(ctx context.Context, userID uint) error

	StartRound(ctx context.Context, roundID uint) error

	getWinner(roundNumber string, totalTickets int, players []*model.Player) (uint, int)

	generateRound() *model.Round

	generateSecret(length int) string

	getHash(roundNumber string, secret string) string

	verifyHash(round *model.Round) bool
}

type service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) Service {
	return &service{repo: repo}
}
