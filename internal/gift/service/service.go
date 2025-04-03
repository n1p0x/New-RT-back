package service

import (
	"context"

	"roulette/internal/gift/model"
	"roulette/internal/gift/repo"
	tgModel "roulette/internal/tg/model"
)

type Service interface {
	GetCollectionByName(ctx context.Context, name string) (*model.Collection, error)

	GetCollections(ctx context.Context) ([]*model.Collection, error)

	GetUserNft(ctx context.Context, userNftID uint) (*model.UserNft, error)

	GetUserGift(ctx context.Context, userGiftID uint) (*model.UserGift, error)

	GetUserGifts(ctx context.Context, userID uint) (*model.UserGifts, int64, error)

	AddUserNft(ctx context.Context, userID uint, name string, collectibleID uint, address string, lottieUrl string, collectionID uint) error

	AddUserGift(ctx context.Context, userID, giftID, msgID int64, name string, collectibleID, collectionID int, lottieUrl string) error

	UpdateCollectionsFloor(ctx context.Context, floors []*tgModel.CollectionFloor) error

	isGiftsAvailable(ctx context.Context, userID uint) (int64, error)
}

type service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) Service {
	return &service{repo: repo}
}
