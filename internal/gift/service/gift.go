package service

import (
	"context"
	"fmt"
	"math/big"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/gift/model"
	tgModel "roulette/internal/tg/model"
)

func (s *service) GetCollectionByName(ctx context.Context, name string) (*model.Collection, error) {
	collection, err := s.repo.GetCollectionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

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
	roundID, _ := s.repo.GetCurrentRoundID(ctx)

	userGifts, err := s.repo.GetUserGifts(ctx, userID, roundID)
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

func (s *service) AddUserGift(ctx context.Context, userID, giftID, msgID int64, name string, collectibleID, collectionID int, lottieUrl string) error {
	gift := &dbModels.GiftDB{
		ID:            giftID,
		MsgID:         msgID,
		Name:          name,
		CollectibleID: collectibleID,
		LottieUrl:     lottieUrl,
		CollectionID:  collectionID,
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
	collections, err := s.repo.GetCollections(ctx)
	if err != nil {
		return fmt.Errorf("failed to get collections: %v", err)
	}

	floorMap := make(map[string]*big.Int)
	for _, floor := range floors {
		floorMap[floor.Name] = floor.Floor
	}

	defaultFloor := big.NewInt(550000000)
	for _, collection := range collections {
		if _, ok := floorMap[collection.Name]; !ok {
			collectionFloor := &tgModel.CollectionFloor{
				Name:  collection.Name,
				Floor: defaultFloor,
			}
			floors = append(floors, collectionFloor)
		}
	}

	fmt.Println(floors)

	var errs []error
	for _, floor := range floors {
		if err = s.repo.UpdateCollectionFloor(ctx, floor.Name, floor.Floor); err != nil {
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
