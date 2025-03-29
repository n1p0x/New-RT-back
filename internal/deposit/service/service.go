package service

import (
	"context"
	"fmt"
	"log"
	"roulette/internal/utils"
	"strconv"
	"time"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/deposit/repo"
	giftService "roulette/internal/gift/service"
	tonService "roulette/internal/ton/service"
	userService "roulette/internal/user/service"
)

type Service interface {
	CheckTonDeposit(ctx context.Context) error

	CheckNftDeposit(ctx context.Context) error

	AddNft(ctx context.Context, userID uint, sender string, nftAddress string) error

	addTon(ctx context.Context, userID uint, userBalance int, amount int, msgHash string, payload *string) error

	addGift(ctx context.Context, userID uint, giftID uint, msgID string, createdAt time.Time, title string, collectibleID uint, lottieUrl string) error
}

type service struct {
	repo        repo.Repo
	tonService  tonService.Service
	userService userService.Service
	giftService giftService.Service
}

func NewService(repo repo.Repo, tonService tonService.Service, userService userService.Service, giftService giftService.Service) Service {
	return &service{
		repo:        repo,
		tonService:  tonService,
		userService: userService,
		giftService: giftService,
	}
}

func (s *service) CheckTonDeposit(ctx context.Context) error {
	start, err := s.repo.GetDepositTime(ctx, dbModels.TonTime)
	if err != nil {
		return fmt.Errorf("failed to get deposit time: %v", err)
	}

	msgs, err := s.tonService.GetTonTransfers(ctx, start)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}

	for _, msg := range msgs {
		if msg.MsgContent.Decoded.Comment == nil {
			continue
		}

		if len(*msg.MsgContent.Decoded.Comment) != 8 {
			continue
		}

		memo := *msg.MsgContent.Decoded.Comment
		user, _ := s.userService.GetUserByMemo(ctx, memo)
		if user == nil {
			continue
		}

		amount, _ := strconv.Atoi(msg.Value)
		if err = s.addTon(ctx, user.ID, user.Balance, amount, msg.Hash, nil); err != nil {
			log.Printf("failed to add ton deposit for user: %d", user.ID)
			continue
		}
	}

	return nil
}

func (s *service) CheckNftDeposit(ctx context.Context) error {
	deps, errDep := s.repo.GetNftDeposits(ctx)
	if errDep != nil {
		return errDep
	}

	for _, dep := range deps {
		transfer, errNft := s.tonService.GetNftTransfer(ctx, dep.NftAddress)

		if errNft != nil {
			continue
		}
		if transfer == nil {
			continue
		}

		transferSender, err := utils.GetAddress(transfer.Sender)
		if err != nil {
			continue
		}
		depSender, err := utils.GetAddress(dep.Sender)
		if err != nil {
			continue
		}
		transferSenderRaw, depSenderRaw := transferSender.StringRaw(), depSender.StringRaw()
		if transferSenderRaw != depSenderRaw {
			continue
		}

		nft, err := s.tonService.GetNft(ctx, dep.NftAddress)
		if err != nil {
			return fmt.Errorf("failed to get nft: %v", dep.NftAddress)
		}

		errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
			isConfirmed := true
			deposit := &dbModels.NftDepositDB{
				TraceID:     &transfer.TraceID,
				IsConfirmed: &isConfirmed,
			}
			if err = s.repo.UpdateNft(ctx, dep.ID, deposit); err != nil {
				if database.IsRecordNotFoundErr(err) {
					return fmt.Errorf("nft deposit %d not found", dep.ID)
				}
				return err
			}

			err = s.giftService.AddUserNft(ctx, dep.UserID, nft.Name, uint(nft.CollectibleID), nft.Address, nft.LottieUrl, 1)
			if err != nil {
				if database.IsKeyConflictErr(err) {
					return fmt.Errorf("nft %s exists", nft.Address)
				}
				return err
			}

			return nil
		})
		if errTx != nil {
			continue
		}
	}

	return nil
}

func (s *service) AddNft(ctx context.Context, userID uint, sender string, nftAddress string) error {
	if _, err := utils.GetAddress(sender); err != nil {
		return fmt.Errorf("invalid sender address: %s", sender)
	}
	if _, err := utils.GetAddress(nftAddress); err != nil {
		return fmt.Errorf("invalid nft address: %s", nftAddress)
	}

	nftDeposit := &dbModels.NftDepositDB{
		UserID:     userID,
		Sender:     sender,
		NftAddress: nftAddress,
	}
	if err := s.repo.AddNft(ctx, nftDeposit); err != nil {
		if database.IsKeyConflictErr(err) {
			return fmt.Errorf("nft deposit %s exists")
		}
		return err
	}

	return nil
}

func (s *service) addTon(ctx context.Context, userID uint, userBalance int, amount int, msgHash string, payload *string) error {
	tonDeposit := &dbModels.TonDepositDB{
		UserID:  userID,
		Amount:  amount,
		MsgHash: msgHash,
		Payload: payload,
	}

	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.repo.AddTon(ctx, tonDeposit); err != nil {
			if database.IsKeyConflictErr(err) {
				return fmt.Errorf("deposit %s is exists", msgHash)
			}
			if database.IsFKeyConflictError(err) {
				return fmt.Errorf("user %s not found", userID)
			}
			return err
		}

		newBalance := userBalance + amount
		if err := s.userService.UpdateUser(ctx, userID, nil, nil, &newBalance); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return fmt.Errorf("failed to execute ton tx: %s", errTx)
	}

	return nil
}

func (s *service) addGift(ctx context.Context, userID uint, giftID uint, msgID string, createdAt time.Time, title string, collectibleID uint, lottieUrl string) error {
	giftDeposit := &dbModels.GiftDepositDB{
		UserID:    userID,
		GiftID:    giftID,
		MsgID:     msgID,
		CreatedAt: createdAt,
	}

	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.repo.AddGift(ctx, giftDeposit); err != nil {
			if database.IsKeyConflictErr(err) {
				return fmt.Errorf("gift deposit %s is exists", msgID)
			}
			return err
		}

		if err := s.giftService.AddUserGift(ctx, userID, giftID, title, collectibleID, lottieUrl); err != nil {
			if database.IsKeyConflictErr(err) {
				return fmt.Errorf("gift %s exists", giftID)
			}
			return err
		}

		return nil
	})
	if errTx != nil {
		return fmt.Errorf("failed to execute add gift tx: %s", errTx)
	}

	return nil
}
