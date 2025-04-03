package repo

import (
	"context"
	"math/big"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/gift/model"
)

func (r *repo) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	// TODO
	return database.RunInTx(ctx, r.db, f)
}

func (r *repo) GetCollectionByName(ctx context.Context, name string) (*model.Collection, error) {
	db := database.FromContext(ctx, r.db)

	var collection *model.Collection
	err := db.WithContext(ctx).
		Model(&dbModels.CollectionDB{}).
		Where("name = ?", name).
		First(&collection).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return collection, nil
}

func (r *repo) GetCollections(ctx context.Context) ([]*model.Collection, error) {
	db := database.FromContext(ctx, r.db)

	var collections []*model.Collection
	err := db.WithContext(ctx).
		Model(&dbModels.CollectionDB{}).
		Find(&collections).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return collections, nil
}

func (r *repo) GetUserNft(ctx context.Context, userNftID uint) (*model.UserNft, error) {
	db := database.FromContext(ctx, r.db)

	var userNft *model.UserNft
	err := db.WithContext(ctx).
		Raw(`
			SELECT un.id, un.user_id, un.nft_id, n.address
			FROM users_nfts un
				LEFT OUTER JOIN nfts n ON un.nft_id = n.id
			WHERE un.id = ?
		`, userNftID).
		Scan(&userNft).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}
	if userNft == nil {
		return nil, database.ErrNotFound
	}

	return userNft, nil
}

func (r *repo) GetUserGift(ctx context.Context, userGiftID uint) (*model.UserGift, error) {
	db := database.FromContext(ctx, r.db)

	var userGift *model.UserGift
	err := db.WithContext(ctx).
		Model(&dbModels.UserGiftDB{}).
		Where("id = ?", userGiftID).
		First(&userGift).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return userGift, nil
}

func (r *repo) GetUserGifts(ctx context.Context, userID uint, roundID uint) (*model.UserGifts, error) {
	db := database.FromContext(ctx, r.db)

	var nfts []*model.Gift
	_ = db.WithContext(ctx).
		Raw(`
			SELECT un.id AS id, n.name AS name, n.collectible_id AS collectible_id, n.lottie_url AS lottie_url, c.floor AS floor, rn.round_id = $2 AS is_bet
			FROM users u
				LEFT OUTER JOIN users_nfts un ON u.id = un.user_id
				LEFT OUTER JOIN nfts n ON n.id = un.nft_id
				LEFT OUTER JOIN collections c ON n.collection_id = c.id
				LEFT OUTER JOIN rounds_nfts rn ON un.id = rn.user_nft_id
			WHERE u.id = $1::BIGINT AND un.id IS NOT NULL
		`, userID, roundID).
		Scan(&nfts)

	//WHERE u.id = $1::BIGINT AND un.id IS NOT NULL AND (rn.round_id = $2 OR rn.round_id = $2 IS NULL)

	var gifts []*model.Gift
	_ = db.WithContext(ctx).
		Raw(`
			SELECT ug.id AS id, g.name AS name, g.collectible_id AS collectible_id, g.lottie_url AS lottie_url, c.floor AS floor, rg.round_id = $2 AS is_bet
			FROM users u
				LEFT OUTER JOIN users_gifts ug ON u.id = ug.user_id
				LEFT OUTER JOIN gifts g ON g.id = ug.gift_id
				LEFT OUTER JOIN collections c ON g.collection_id = c.id
				LEFT OUTER JOIN rounds_gifts rg ON ug.id = rg.user_gift_id
			WHERE u.id = $1::BIGINT AND ug.id IS NOT NULL
		`, userID, roundID).
		Scan(&gifts)

	//WHERE u.id = $1::BIGINT AND ug.id IS NOT NULL AND (rg.round_id = $2 OR rg.round_id = $2 IS NULL)

	if nfts == nil && gifts == nil {
		return nil, database.ErrNotFound
	}

	userGifts := &model.UserGifts{
		Nfts:  nfts,
		Gifts: gifts,
	}

	return userGifts, nil
}

func (r *repo) GetWinnerFee(ctx context.Context, userID uint) (int64, error) {
	db := database.FromContext(ctx, r.db)

	var fee int64
	err := db.WithContext(ctx).
		Raw(`
			SELECT COALESCE(SUM(fee), 0) AS fee
			FROM rounds_winners rw
			WHERE rw.user_id = $1 AND rw.is_paid IS NULL
			GROUP BY rw.user_id
		`, userID).
		Scan(&fee).Error
	if err != nil {
		return 0, nil
	}

	return fee, nil
}

func (r *repo) GetCurrentRoundID(ctx context.Context) (uint, error) {
	db := database.FromContext(ctx, r.db)

	var roundID uint
	err := db.WithContext(ctx).
		Raw(`
			SELECT id
			FROM rounds
			ORDER BY created_at DESC
			LIMIT 1
		`).
		First(&roundID).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return 0, database.ErrNotFound
		}
		return 0, err
	}

	return roundID, nil
}

func (r *repo) AddNft(ctx context.Context, nft *dbModels.NftDB) (uint, error) {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&dbModels.NftDB{}).
		Create(&nft).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return 0, database.ErrKeyConflict
		}
		return 0, err
	}

	return nft.ID, nil
}

func (r *repo) AddGift(ctx context.Context, gift *dbModels.GiftDB) (int64, error) {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).Create(&gift).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return 0, database.ErrKeyConflict
		}
		return 0, err
	}

	return gift.ID, nil
}

func (r *repo) AddUserNft(ctx context.Context, userID uint, nftID uint) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Select("user_id", "nft_id").
		Create(&dbModels.UserNftDB{UserID: userID, NftID: nftID}).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) AddUserGift(ctx context.Context, userID int64, giftID int64) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Select("user_id", "gift_id").
		Create(&dbModels.UserGiftDB{UserID: userID, GiftID: giftID}).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) UpdateCollectionFloor(ctx context.Context, collectionName string, floor *big.Int) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Model(&dbModels.CollectionDB{}).
		Where("name = ?", collectionName).
		UpdateColumn("floor", floor)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}
