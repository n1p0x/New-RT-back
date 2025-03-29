package repo

//go:generate moq --out=mocks/mock_repo.go --pkg=mocks . Repo

import (
	"context"

	"gorm.io/gorm"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/user/model"
)

type Repo interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error

	GetUser(ctx context.Context, userID uint) (*model.User, error)

	GetUserProfile(ctx context.Context, userID uint) (*model.UserProfile, error)

	GetUserByMemo(ctx context.Context, memo string) (*model.User, error)

	AddUser(ctx context.Context, user *dbModels.UserDB) error

	AddReferral(ctx context.Context, ref *dbModels.ReferralDB) error

	UpdateUser(ctx context.Context, userID uint, user *dbModels.UserDB) error
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

func (r *repo) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	db := database.FromContext(ctx, r.db)

	var user *model.User
	err := db.WithContext(ctx).
		Model(&dbModels.UserDB{}).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repo) GetUserProfile(ctx context.Context, userID uint) (*model.UserProfile, error) {
	db := database.FromContext(ctx, r.db)

	var userProfile *model.UserProfile
	err := db.WithContext(ctx).
		Raw(`
			SELECT
				u.balance,
				(
					SELECT COALESCE(COUNT(r.ref_id), 0) AS refs
					FROM referrals r
					WHERE r.referrer_id = $1
				),
				(
					SELECT COALESCE(SUM(r.fee), 0) AS earned
					FROM referrals_fees r
					WHERE r.user_id = $1
				),
				(
					SELECT COUNT(*) FROM (
						SELECT DISTINCT rt.round_id
						FROM rounds_tickets rt
						WHERE rt.user_id = $1
					) AS games
				)
			FROM
				users u
			WHERE u.id = $1
		`, userID).
		Scan(&userProfile).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return userProfile, nil
}

func (r *repo) GetUserByMemo(ctx context.Context, memo string) (*model.User, error) {
	db := database.FromContext(ctx, r.db)

	var user *model.User
	err := db.WithContext(ctx).
		Model(&dbModels.UserDB{}).
		Where("memo = ?", memo).
		First(&user).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repo) AddUser(ctx context.Context, user *dbModels.UserDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Select("id", "username", "photo_url", "memo").
		Create(&user).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (r *repo) AddReferral(ctx context.Context, ref *dbModels.ReferralDB) error {
	db := database.FromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Select("referrer_id", "ref_id").
		Create(&ref).Error
	if err != nil {
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (r *repo) UpdateUser(ctx context.Context, userID uint, user *dbModels.UserDB) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Where("id = ?", userID).
		Updates(user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}
