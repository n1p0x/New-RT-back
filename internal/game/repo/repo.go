package repo

import (
	"context"

	"gorm.io/gorm"

	dbModels "roulette/internal/database/models"
	"roulette/internal/game/model"
)

type Repo interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error

	GetCurrentRound(ctx context.Context) (*model.Round, error)

	GetRound(ctx context.Context, roundID uint) (*model.Round, error)

	GetRoundStats(ctx context.Context, roundID uint) (*model.RoundStats, error)

	GetPlayers(ctx context.Context, roundID uint) ([]*model.Player, error)

	GetUniquePlayers(ctx context.Context, roundID uint) ([]*model.UniquePlayer, error)

	GetRoundNftIDs(ctx context.Context, roundID uint) ([]uint, error)

	GetRoundGiftIDs(ctx context.Context, roundID uint) ([]uint, error)

	GetUserNft(ctx context.Context, userNftID uint) (*model.Gift, error)

	GetUserGift(ctx context.Context, userGiftID uint) (*model.Gift, error)

	GetWinner(ctx context.Context, roundID uint) (*model.Winner, error)

	GetWinnerFee(ctx context.Context, userID uint) (int64, error)

	GetUserBalance(ctx context.Context, userID uint) (int64, error)

	GetReferrer(ctx context.Context, refID uint) (*model.Referrer, error)

	AddRound(ctx context.Context, round *dbModels.RoundDB) error

	AddUserRoundNft(ctx context.Context, roundNft *dbModels.RoundNftDB) error

	AddUserRoundGift(ctx context.Context, roundGift *dbModels.RoundGiftDB) error

	AddUserRoundTicket(ctx context.Context, roundTicket *dbModels.RoundTicketDB) error

	AddWinner(ctx context.Context, winner *dbModels.RoundWinnerDB) error

	AddReferralFee(ctx context.Context, fee *dbModels.ReferralFeeDB) error

	UpdateRoundStart(ctx context.Context, roundID uint) error

	UpdateUserNftOwner(ctx context.Context, ownerID uint, userNftID ...uint) error

	UpdateUserGiftOwner(ctx context.Context, ownerID uint, userGiftID ...uint) error

	UpdateWinner(ctx context.Context, winnerID uint, winner *dbModels.RoundWinnerDB) error

	UpdateFee(ctx context.Context, userID uint) error

	UpdateUserBalance(ctx context.Context, userID uint, balance int64) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}
