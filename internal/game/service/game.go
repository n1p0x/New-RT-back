package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"time"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/game/model"
)

func (s *service) GetCurrentRound(ctx context.Context) (*model.Round, error) {
	// TODO: cash
	round, err := s.repo.GetCurrentRound(ctx)
	if err != nil {
		return nil, err
	}
	return round, nil
}

func (s *service) GetCurrentRoundWithPlayers(ctx context.Context) (*model.RoundWithPlayers, error) {
	// TODO: cash
	curRound, err := s.repo.GetCurrentRound(ctx)
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, ErrRoundNotFound
		}
		return nil, err
	}

	players, err := s.repo.GetPlayers(ctx, curRound.ID)
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			stats := &model.RoundStats{
				TotalGifts:   0,
				TotalBet:     0,
				TotalTickets: 0,
			}
			return &model.RoundWithPlayers{
				Round:      curRound,
				RoundStats: stats,
			}, nil
		}
		return nil, err
	}

	uniquePlayers, _ := s.repo.GetUniquePlayers(ctx, curRound.ID)

	stats, err := s.repo.GetRoundStats(ctx, curRound.ID)

	round := &model.RoundWithPlayers{
		Round:         curRound,
		RoundStats:    stats,
		Players:       players,
		UniquePlayers: uniquePlayers,
	}

	return round, nil
}

func (s *service) GetRoundWithPlayers(ctx context.Context, roundID uint) (*model.RoundWithPlayers, error) {
	// TODO: cash
	round, err := s.repo.GetRound(ctx, roundID)
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, ErrRoundNotFound
		}
		return nil, err
	}

	players, err := s.repo.GetPlayers(ctx, round.ID)
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			stats := &model.RoundStats{
				TotalGifts:   0,
				TotalBet:     0,
				TotalTickets: 0,
			}
			return &model.RoundWithPlayers{
				Round:      round,
				RoundStats: stats,
			}, nil
		}
		return nil, err
	}

	uniquePlayers, _ := s.repo.GetUniquePlayers(ctx, round.ID)

	stats, err := s.repo.GetRoundStats(ctx, round.ID)

	roundWithPlayers := &model.RoundWithPlayers{
		Round:         round,
		RoundStats:    stats,
		Players:       players,
		UniquePlayers: uniquePlayers,
	}

	return roundWithPlayers, nil
}

func (s *service) GetWinner(ctx context.Context, roundID uint) (*model.Winner, error) {
	winner, err := s.repo.GetWinner(ctx, roundID)
	if err != nil {
		return nil, err
	}

	return winner, nil
}

func (s *service) AddRound(ctx context.Context) error {
	round := s.generateRound()

	roundDB := &dbModels.RoundDB{
		RoundNumber: round.RoundNumber,
		Secret:      round.Secret,
		Hash:        round.Hash,
	}
	err := s.repo.AddRound(ctx, roundDB)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddUserNft(ctx context.Context, userNftID uint) error {
	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		round, err := s.repo.GetCurrentRound(ctx)
		if err != nil {
			return ErrRoundNotFound
		}

		userNft, err := s.repo.GetUserNft(ctx, userNftID)
		if err != nil {
			return err
		}

		roundNft := &dbModels.RoundNftDB{
			RoundID:   round.ID,
			UserNftID: userNftID,
			Bet:       userNft.Floor,
		}
		err = s.repo.AddUserRoundNft(ctx, roundNft)
		if err != nil {
			return err
		}

		roundTicket := &dbModels.RoundTicketDB{
			RoundID: round.ID,
			UserID:  userNft.UserID,
			Tickets: int(userNft.Floor / 100_000_000),
		}
		err = s.repo.AddUserRoundTicket(ctx, roundTicket)
		if err != nil {
			return err
		}

		if time.Now().Unix() >= round.StartedAt.Add(3*time.Minute).Unix() {
			return ErrRoundFinished
		}

		return nil
	})
	if errTx != nil {
		return fmt.Errorf("failed to add user: %v", errTx)
	}

	return nil
}

func (s *service) AddUserGift(ctx context.Context, userGiftID uint) error {
	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		round, err := s.repo.GetCurrentRound(ctx)
		if err != nil {
			return ErrRoundNotFound
		}

		userGift, err := s.repo.GetUserGift(ctx, userGiftID)
		if err != nil {
			return err
		}

		roundGift := &dbModels.RoundGiftDB{
			RoundID:    round.ID,
			UserGiftID: userGiftID,
			Bet:        userGift.Floor,
		}
		err = s.repo.AddUserRoundGift(ctx, roundGift)
		if err != nil {
			return err
		}

		roundTicket := &dbModels.RoundTicketDB{
			RoundID: round.ID,
			UserID:  userGift.UserID,
			Tickets: int(userGift.Floor / 100_000_000),
		}
		err = s.repo.AddUserRoundTicket(ctx, roundTicket)
		if err != nil {
			return err
		}

		if time.Now().Unix() >= round.StartedAt.Add(3*time.Minute).Unix() {
			return ErrRoundFinished
		}

		return nil
	})
	if errTx != nil {
		return fmt.Errorf("failed to add user: %v", errTx)
	}

	return nil
}

func (s *service) AddWinner(ctx context.Context, roundWithPlayers *model.RoundWithPlayers) error {
	userID, ticket := s.getWinner(roundWithPlayers.RoundNumber, roundWithPlayers.TotalTickets, roundWithPlayers.Players)
	fee := roundWithPlayers.TotalBet / 100 * 5

	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		nfts, err := s.repo.GetRoundNftIDs(ctx, roundWithPlayers.ID)
		if err == nil {
			if err = s.repo.UpdateUserNftOwner(ctx, userID, nfts...); err != nil {
				return err
			}
		}

		gifts, err := s.repo.GetRoundGiftIDs(ctx, roundWithPlayers.ID)
		if err == nil {
			if err = s.repo.UpdateUserGiftOwner(ctx, userID, gifts...); err != nil {
				return err
			}
		}

		winnerRound := &dbModels.RoundWinnerDB{
			RoundID: roundWithPlayers.ID,
			Ticket:  ticket,
			UserID:  userID,
			Fee:     fee,
		}
		if err = s.repo.AddWinner(ctx, winnerRound); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (s *service) UpdateFee(ctx context.Context, userID uint) error {
	errTx := s.repo.RunInTx(ctx, func(ctx context.Context) error {
		balance, err := s.repo.GetUserBalance(ctx, userID)
		if err != nil {
			return err
		}
		fee, err := s.repo.GetWinnerFee(ctx, userID)
		if err != nil {
			return err
		}

		newBalance := balance - fee
		if newBalance < 0 {
			return ErrNotEnoughBalance
		}

		if err = s.repo.UpdateUserBalance(ctx, userID, newBalance); err != nil {
			return err
		}
		if err = s.repo.UpdateFee(ctx, userID); err != nil {
			return err
		}

		referrer, err := s.repo.GetReferrer(ctx, userID)
		if referrer == nil || err != nil {
			return nil
		}

		refFee := newBalance / 100 * 5
		if referrer.IsSpec == true {
			refFee = newBalance / 100 * 15
		}
		newReferrerBalance := referrer.Balance + refFee
		if err = s.repo.UpdateUserBalance(ctx, referrer.ID, newReferrerBalance); err != nil {
			return err
		}
		referralFee := &dbModels.ReferralFeeDB{
			UserID: referrer.ID,
			Amount: fee,
			Fee:    refFee,
		}
		if err = s.repo.AddReferralFee(ctx, referralFee); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (s *service) StartRound(ctx context.Context, roundID uint) error {
	if err := s.repo.UpdateRoundStart(ctx, roundID); err != nil {
		return err
	}
	return nil
}

func (s *service) getWinner(round string, totalTickets int, players []*model.Player) (uint, int) {
	roundNumber, _ := strconv.ParseFloat(round, 64)
	winningTicket := int(math.Floor(float64(totalTickets)*roundNumber)) + 1

	currentTicket := 0
	for _, player := range players {
		currentTicket += player.Tickets
		if currentTicket >= winningTicket {
			return player.UserID, winningTicket
		}
	}

	return 0, -1
}

func (s *service) generateRound() *model.Round {
	roundNumber := fmt.Sprintf("%.16f", rand.Float64())
	secret := s.generateSecret(10)
	hash := s.getHash(roundNumber, secret)

	return &model.Round{
		RoundNumber: roundNumber,
		Secret:      secret,
		Hash:        hash,
	}
}

func (s *service) generateSecret(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func (s *service) getHash(roundNumber string, secret string) string {
	data := fmt.Sprintf("%s%s", roundNumber, secret)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *service) verifyHash(round *model.Round) bool {
	data := fmt.Sprintf("%f%s", round.RoundNumber, round.Secret)
	hash := md5.Sum([]byte(data))
	calculatedHash := hex.EncodeToString(hash[:])
	return calculatedHash == round.Hash
}
