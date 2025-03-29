package repo

import (
	"context"
	"math/big"

	"gorm.io/gorm"

	dbModels "roulette/internal/database/models"
	"roulette/internal/gift/model"
)

type Repo interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error

	GetCollectionByName(ctx context.Context, name string) (*model.Collection, error)

	GetCollections(ctx context.Context) ([]*model.Collection, error)

	GetUserNft(ctx context.Context, userNftID uint) (*model.UserNft, error)

	GetUserGift(ctx context.Context, userGiftID uint) (*model.UserGift, error)

	GetUserGifts(ctx context.Context, userID uint, roundID uint) (*model.UserGifts, error)

	GetWinnerFee(ctx context.Context, userID uint) (int64, error)

	AddNft(ctx context.Context, nft *dbModels.NftDB) (uint, error)

	AddGift(ctx context.Context, gift *dbModels.GiftDB) (uint, error)

	AddUserNft(ctx context.Context, userID uint, nftID uint) error

	AddUserGift(ctx context.Context, userID uint, giftID uint) error

	UpdateCollectionFloor(ctx context.Context, collectionName string, floor *big.Int) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}
