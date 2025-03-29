package service

//go:generate moq --out=mocks/mock_service.go --pkg=mocks . Service

import (
	"context"
	"fmt"
	"hash/crc32"
	"strconv"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/user/model"
	"roulette/internal/user/repo"
)

type Service interface {
	GetUser(ctx context.Context, userID uint) (*model.User, error)

	GetUserProfile(ctx context.Context, userID uint) (*model.UserProfile, error)

	GetUserByMemo(ctx context.Context, memo string) (*model.User, error)

	AddUser(ctx context.Context, userID uint, name, photoUrl, startParam *string) error

	UpdateUser(ctx context.Context, userID uint, name *string, photoUrl *string, balance *int) error

	generateMemo(userID uint) string
}

type service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) Service {
	return &service{repo: repo}
}

func (s *service) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *service) GetUserProfile(ctx context.Context, userID uint) (*model.UserProfile, error) {
	userProfile, err := s.repo.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (s *service) GetUserByMemo(ctx context.Context, memo string) (*model.User, error) {
	user, err := s.repo.GetUserByMemo(ctx, memo)
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *service) AddUser(ctx context.Context, userID uint, name, photoUrl, startParam *string) error {
	memo := s.generateMemo(userID)
	user := &dbModels.UserDB{
		ID:       userID,
		Name:     name,
		PhotoUrl: photoUrl,
		Memo:     memo,
	}

	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.repo.AddUser(ctx, user); err != nil {
			if database.IsKeyConflictErr(err) {
				return fmt.Errorf("user %d exists", userID)
			}
			return err
		}

		if startParam != nil {
			referrerID, err := strconv.ParseUint(*startParam, 10, 64)
			if err != nil {
				return nil
			}

			if _, err = s.repo.GetUser(ctx, uint(referrerID)); err != nil {
				return nil
			}

			ref := &dbModels.ReferralDB{
				ReferrerID: uint(referrerID),
				RefID:      userID,
			}
			if err = s.repo.AddReferral(ctx, ref); err != nil {
				return err
			}
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, userID uint, name *string, photoUrl *string, balance *int) error {
	user := &dbModels.UserDB{
		Name:     name,
		PhotoUrl: photoUrl,
		Balance:  *balance,
	}

	if err := s.repo.UpdateUser(ctx, userID, user); err != nil {
		return err
	}

	return nil
}

func (s *service) generateMemo(userID uint) string {
	userIdStr := fmt.Sprintf("%v", userID)
	crc := crc32.ChecksumIEEE([]byte(userIdStr))
	memo := fmt.Sprintf("%08x", crc)

	return memo
}
