package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
)

type Repo interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error

	AddTon(ctx context.Context, tonWithdraw *dbModels.TonWithdrawDB) error

	AddNft(ctx context.Context, nftWithdraw *dbModels.NftWithdrawDB) error

	AddGift(ctx context.Context, giftWithdraw *dbModels.GiftWithdrawDB) error

	UpdateUserBalance(ctx context.Context, userID uint, amount int) error

	DeleteNft(ctx context.Context, nftID uint) error

	DeleteGift(ctx context.Context, giftID uint) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	// TODO
	return database.RunInTx(ctx, r.db, f)
}

func (r *repo) AddTon(ctx context.Context, tonWithdraw *dbModels.TonWithdrawDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&dbModels.TonWithdrawDB{}).
		Create(&tonWithdraw).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) AddNft(ctx context.Context, nftWithdraw *dbModels.NftWithdrawDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&dbModels.NftWithdrawDB{}).
		Create(&nftWithdraw).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) AddGift(ctx context.Context, giftWithdraw *dbModels.GiftWithdrawDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&dbModels.GiftWithdrawDB{}).
		Create(&giftWithdraw).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) UpdateUserBalance(ctx context.Context, userID uint, amount int) error {
	db := database.FromContext(ctx, r.db)

	var user *dbModels.UserDB
	err := db.WithContext(ctx).
		Model(&dbModels.UserDB{}).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return database.ErrNotFound
		}
		return err
	}

	if user.Balance < amount {
		return fmt.Errorf("user %d has not enough funds", userID)
	}

	balance := user.Balance - amount
	res := db.WithContext(ctx).
		Model(&dbModels.UserDB{}).
		Where("id = ?", userID).
		UpdateColumn("balance", balance)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

func (r *repo) DeleteNft(ctx context.Context, nftID uint) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Model(&dbModels.NftDB{}).
		Delete(&dbModels.NftDB{}, "id = ?", nftID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

func (r *repo) DeleteGift(ctx context.Context, giftID uint) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Model(&dbModels.GiftDB{}).
		Where("id = $1", giftID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}
