package service

import (
	"context"
	"time"

	dbModels "roulette/internal/database/models"
	giftService "roulette/internal/gift/service"
	tgService "roulette/internal/tg/service"
	tonService "roulette/internal/ton/service"
	userService "roulette/internal/user/service"
	"roulette/internal/withdraw/repo"
)

type Service interface {
	// AddTon adds new ton withdraw
	AddTon(ctx context.Context, userID uint, dst string, amount uint) error

	// AddNft adds new nft withdraw
	AddNft(ctx context.Context, userNftID uint, dst string) error

	// AddGift adds new gift withdraw
	AddGift(ctx context.Context, userGiftID uint) error
}

type service struct {
	repo        repo.Repo
	tonService  tonService.Service
	tgService   tgService.Service
	giftService giftService.Service
	userService userService.Service
}

func NewService(repo repo.Repo, tonService tonService.Service, tgService tgService.Service, giftService giftService.Service, userService userService.Service) Service {
	return &service{
		repo:        repo,
		tonService:  tonService,
		tgService:   tgService,
		giftService: giftService,
		userService: userService,
	}
}

func (s *service) AddTon(ctx context.Context, userID uint, dst string, amount uint) error {
	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdateUserBalance(ctx, userID, int(amount)+TonFee); err != nil {
			return err
		}

		_, err := s.tonService.SendTon(ctx, dst, amount)
		if err != nil {
			return err
		}

		tonWithdraw := &dbModels.TonWithdrawDB{
			UserID:      userID,
			Destination: dst,
			Amount:      amount,
			CreatedAt:   time.Time{},
		}
		if err = s.repo.AddTon(ctx, tonWithdraw); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (s *service) AddNft(ctx context.Context, userNftID uint, dst string) error {
	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		userNft, err := s.giftService.GetUserNft(ctx, userNftID)
		if err != nil {
			return err
		}

		if err = s.repo.UpdateUserBalance(ctx, userNft.UserID, NftFee); err != nil {
			return err
		}

		_, err = s.tonService.SendNft(ctx, dst, userNft.Address)
		if err != nil {
			return err
		}

		if err = s.repo.DeleteNft(ctx, userNft.NftID); err != nil {
			return err
		}

		nftWithdraw := &dbModels.NftWithdrawDB{
			UserID:      userNft.UserID,
			Destination: dst,
			NftAddress:  userNft.Address,
		}
		if err = s.repo.AddNft(ctx, nftWithdraw); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (s *service) AddGift(ctx context.Context, userGiftID uint) error {
	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		userGift, err := s.giftService.GetUserGift(ctx, userGiftID)
		if err != nil {
			return err
		}

		if err = s.repo.UpdateUserBalance(ctx, userGift.UserID, GiftFee); err != nil {
			return err
		}

		if err = s.tgService.SendGift(ctx, int64(userGift.UserID), userGift.MsgID); err != nil {
			return err
		}

		if err = s.repo.DeleteGift(ctx, userGift.GiftID); err != nil {
			return err
		}

		giftWithdraw := &dbModels.GiftWithdrawDB{
			UserID:    userGift.UserID,
			GiftID:    userGift.GiftID,
			CreatedAt: time.Time{},
		}
		if err = s.repo.AddGift(ctx, giftWithdraw); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}
