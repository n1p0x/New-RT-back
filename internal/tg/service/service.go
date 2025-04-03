package service

import (
	"context"

	"github.com/celestix/gotgproto"

	"roulette/internal/config"
	"roulette/internal/tg/model"
)

type Service interface {
	SendGift(ctx context.Context, userID int64, msgID int) error

	GetClient(ctx context.Context) (*gotgproto.Client, error)

	GetFloors(ctx context.Context, channelID int64) ([]*model.CollectionFloor, error)

	GetFloorsLow(ctx context.Context, channelID, accessHash int64) ([]*model.CollectionFloor, error)

	getChannelMessages(ctx context.Context, channelID int64) (string, error)
}

type service struct {
	ClientID    int
	ClientHash  string
	ClientPhone string
}

func NewService(cfg *config.Config) Service {
	return &service{
		ClientID:    cfg.TgConfig.ClientID,
		ClientHash:  cfg.TgConfig.ClientHash,
		ClientPhone: cfg.TgConfig.ClientPhone,
	}
}
