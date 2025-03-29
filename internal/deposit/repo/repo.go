package repo

import (
	"context"
	"roulette/internal/database"

	"gorm.io/gorm"

	dbModels "roulette/internal/database/models"
	"roulette/internal/deposit/model"
)

type Repo interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error

	GetDepositTime(ctx context.Context, timeID dbModels.DepositTime) (*int64, error)

	GetNftDeposits(ctx context.Context) ([]*model.NftDeposit, error)

	AddTon(ctx context.Context, tonDeposit *dbModels.TonDepositDB) error

	AddNft(ctx context.Context, nftDeposit *dbModels.NftDepositDB) error

	AddGift(ctx context.Context, giftDeposit *dbModels.GiftDepositDB) error

	UpdateDepositTime(ctx context.Context, timeID uint) error

	UpdateUserBalance(ctx context.Context, userID uint, amount int) error

	UpdateNft(ctx context.Context, depositID uint, deposit *dbModels.NftDepositDB) error
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

func (r *repo) GetDepositTime(ctx context.Context, timeID dbModels.DepositTime) (*int64, error) {
	db := database.FromContext(ctx, r.db)

	var depositTime *dbModels.DepositTimeDB
	err := db.WithContext(ctx).First(&depositTime, "id = ?", timeID).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			err = db.WithContext(ctx).
				Raw(`
					INSERT INTO deposit_time (id) VALUES (?)
					ON CONFLICT (id) DO UPDATE SET start = CURRENT_TIMESTAMP
					RETURNING start
				`, timeID).
				Scan(&depositTime).Error
			if err != nil {
				return nil, err
			}

			unixTime := depositTime.Start.Unix()

			return &unixTime, nil
		}
		return nil, err
	}

	res := db.WithContext(ctx).
		Raw(`
			UPDATE deposit_time 
			SET start = CURRENT_TIMESTAMP 
			WHERE id = ?
		`, timeID).
		Save(&dbModels.DepositTimeDB{})
	if res.Error != nil {
		return nil, err
	}

	unixTime := depositTime.Start.Unix()

	return &unixTime, nil
}

func (r *repo) GetNftDeposits(ctx context.Context) ([]*model.NftDeposit, error) {
	db := database.FromContext(ctx, r.db)

	var deps []*model.NftDeposit
	err := db.WithContext(ctx).
		Raw(`
			SELECT *
			FROM nft_deposits
			WHERE is_confirmed IS NULL OR 
   				  created_at BETWEEN CURRENT_TIMESTAMP - INTERVAL '10 minutes' AND CURRENT_TIMESTAMP
		`).Scan(&deps).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return deps, nil
}

func (r *repo) AddTon(ctx context.Context, tonDeposit *dbModels.TonDepositDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&dbModels.TonDepositDB{}).
		Select("user_id", "amount", "payload", "msg_hash").
		Save(&tonDeposit).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		if database.IsFKeyConflictError(err) {
			return database.ErrFKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) AddNft(ctx context.Context, nftDeposit *dbModels.NftDepositDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&dbModels.NftDepositDB{}).
		Select("user_id", "sender", "nft_address", "trace_id").
		Create(&nftDeposit).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) AddGift(ctx context.Context, giftDeposit *dbModels.GiftDepositDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&dbModels.GiftDepositDB{}).
		Create(&giftDeposit).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}

	return nil
}

func (r *repo) UpdateDepositTime(ctx context.Context, timeID uint) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Model(&dbModels.DepositTimeDB{}).
		Where("id = ?", timeID).
		Update("id", timeID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
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

	balance := user.Balance + amount
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

func (r *repo) UpdateNft(ctx context.Context, depositID uint, deposit *dbModels.NftDepositDB) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Where("id = ?", depositID).
		Updates(deposit)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}
