package service

import (
	"context"
	"fmt"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/gift/model"
	tgModel "roulette/internal/tg/model"
)

func (s *service) GetCollections(ctx context.Context) ([]*model.Collection, error) {
	collections, err := s.repo.GetCollections(ctx)
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (s *service) GetUserNft(ctx context.Context, userNftID uint) (*model.UserNft, error) {
	userNft, err := s.repo.GetUserNft(ctx, userNftID)
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, fmt.Errorf("user nft %v not found", userNftID)
		}
		return nil, err
	}
	return userNft, nil
}

func (s *service) GetUserGift(ctx context.Context, userGiftID uint) (*model.UserGift, error) {
	userGift, err := s.repo.GetUserGift(ctx, userGiftID)
	if err != nil {
		return nil, err
	}
	return userGift, nil
}

func (s *service) GetUserGifts(ctx context.Context, userID uint) (*model.UserGifts, int64, error) {
	userGifts, err := s.repo.GetUserGifts(ctx, userID, 1)
	if err != nil {
		return nil, 0, err
	}

	fee, err := s.isGiftsAvailable(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return userGifts, fee, nil
}

func (s *service) AddUserNft(ctx context.Context, userID uint, name string, collectibleID uint, address string, lottieUrl string, collectionID uint) error {
	nft := &dbModels.NftDB{
		Name:          name,
		CollectibleID: collectibleID,
		Address:       address,
		LottieUrl:     lottieUrl,
		CollectionID:  collectionID,
	}
	nftID, err := s.repo.AddNft(ctx, nft)
	if err != nil {
		return err
	}

	if err = s.repo.AddUserNft(ctx, userID, nftID); err != nil {
		return err
	}

	return nil
}

func (s *service) AddUserGift(ctx context.Context, userID uint, giftID uint, title string, collectibleID uint, lottieUrl string) error {
	gift := &dbModels.GiftDB{
		ID:            giftID,
		Title:         title,
		CollectibleID: collectibleID,
		LottieUrl:     lottieUrl,
	}

	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		newGiftID, err := s.repo.AddGift(ctx, gift)
		if err != nil {
			return err
		}

		if err = s.repo.AddUserGift(ctx, userID, newGiftID); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (s *service) UpdateCollectionsFloor(ctx context.Context, floors []*tgModel.CollectionFloor) error {
	var errs []error
	for _, floor := range floors {
		if err := s.repo.UpdateCollectionFloor(ctx, floor.Name, floor.Floor); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("errs: %v", errs)
}

func (s *service) isGiftsAvailable(ctx context.Context, userID uint) (int64, error) {
	fee, err := s.repo.GetWinnerFee(ctx, userID)
	if err != nil {
		return 0, err
	}
	return fee, nil
}
