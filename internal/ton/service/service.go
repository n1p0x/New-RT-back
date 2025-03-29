package service

import (
	"context"

	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"roulette/internal/config"
	giftModel "roulette/internal/gift/model"
	"roulette/internal/models"
	"roulette/internal/ton/model"
)

type Service interface {
	GetTonTransfers(ctx context.Context, start *int64) ([]*model.Message, error)

	GetNftTransfer(ctx context.Context, itemAddress string) (*model.NftTransfer, error)

	GetNftTransfers(ctx context.Context, start *int64) ([]*model.NftTransfer, error)

	GetWalletNfts(ctx context.Context, wallet string, collection []*giftModel.Collection) ([]*models.Nft, error)

	GetNft(ctx context.Context, address string) (*models.Nft, error)

	SendTon(ctx context.Context, dst string, amount uint) (string, error)

	SendNft(ctx context.Context, dst string, nftAddress string) (string, error)

	getWallet(ctx context.Context) (*wallet.Wallet, ton.APIClientWrapped, error)
}

type service struct {
	IsTestnet              bool
	TonCenterApiKey        string
	TonCenterApiKeyTestnet string
	AdminWallet            string
	Mnemonic               string
}

func NewService(cfg *config.Config) Service {
	return &service{
		IsTestnet:              cfg.TonConfig.IsTestnet,
		TonCenterApiKey:        cfg.TonConfig.TonCenterApiKey,
		TonCenterApiKeyTestnet: cfg.TonConfig.TonCenterApiKeyTestnet,
		AdminWallet:            cfg.TonConfig.AdminWallet,
		Mnemonic:               cfg.TonConfig.Mnemonic,
	}
}
