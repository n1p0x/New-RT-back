package repo

import (
	"context"

	"roulette/internal/database"
	dbModels "roulette/internal/database/models"
	"roulette/internal/game/model"
)

func (r *repo) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	// TODO
	return database.RunInTx(ctx, r.db, f)
}

func (r *repo) GetCurrentRound(ctx context.Context) (*model.Round, error) {
	db := database.FromContext(ctx, r.db)

	var round *model.Round
	err := db.WithContext(ctx).
		Raw(`
			SELECT id, round AS round_number, secret, hash, created_at, started_at, CURRENT_TIMESTAMP > started_at + INTERVAL '3 minutes' AS is_finished
			FROM rounds
			ORDER BY created_at DESC
			LIMIT 1
		`).
		First(&round).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return round, nil
}

func (r *repo) GetRound(ctx context.Context, roundID uint) (*model.Round, error) {
	db := database.FromContext(ctx, r.db)

	var round *model.Round
	err := db.WithContext(ctx).
		Raw(`
			SELECT id, round AS round_number, secret, hash, created_at, started_at, CURRENT_TIMESTAMP > started_at + INTERVAL '3 minutes' AS is_finished
			FROM rounds
			WHERE id = $1
		`, roundID).
		First(&round).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return round, nil
}

func (r *repo) GetRoundStats(ctx context.Context, roundID uint) (*model.RoundStats, error) {
	db := database.FromContext(ctx, r.db)

	var stats *model.RoundStats
	err := db.WithContext(ctx).
		Raw(`
			SELECT 
                COALESCE(COUNT(bet), 0) AS total_gifts,
                COALESCE(SUM(bet), 0) AS total_bet,
                COALESCE((
                    SELECT SUM(tickets) 
                    FROM rounds_tickets 
                    WHERE round_id = $1
                ), 0) AS total_tickets
            FROM (
                SELECT bet FROM rounds_nfts WHERE round_id = $1
                UNION ALL
                SELECT bet FROM rounds_gifts WHERE round_id = $1
            ) AS bets
		`, roundID).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *repo) GetPlayers(ctx context.Context, roundID uint) ([]*model.Player, error) {
	db := database.FromContext(ctx, r.db)

	var players []*model.Player
	err := db.WithContext(ctx).
		Raw(`
			SELECT rt.id AS id, rt.user_id AS user_id, rt.round_id AS round_id, rt.tickets AS tickets, rt.created_at AS created_at
			FROM rounds_tickets rt
			WHERE round_id = $1
			ORDER BY tickets DESC
		`, roundID).
		Scan(&players).Error
	if len(players) == 0 {
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return players, nil
}

func (r *repo) GetUniquePlayers(ctx context.Context, roundID uint) ([]*model.UniquePlayer, error) {
	db := database.FromContext(ctx, r.db)

	var players []*model.UniquePlayer
	err := db.WithContext(ctx).
		Raw(`
			SELECT rt.user_id AS user_id, u.name AS name, u.photo_url AS photo_url, rt.round_id AS round_id, SUM(rt.tickets) AS tickets
			FROM rounds_tickets rt
				LEFT OUTER JOIN users u ON rt.user_id = u.id
			WHERE round_id = $1
			GROUP BY rt.user_id, rt.round_id, u.name, u.photo_url
			ORDER BY tickets DESC
		`, roundID).
		Scan(&players).Error
	if len(players) == 0 {
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return players, nil
}

func (r *repo) GetRoundNftIDs(ctx context.Context, roundID uint) ([]uint, error) {
	db := database.FromContext(ctx, r.db)

	var nftIDs []uint
	err := db.WithContext(ctx).
		Raw(`
			SELECT rn.user_nft_id AS id
			FROM rounds_nfts rn
			WHERE rn.round_id = $1
		`, roundID).
		Scan(&nftIDs).Error
	if len(nftIDs) == 0 {
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return nftIDs, nil
}

func (r *repo) GetRoundGiftIDs(ctx context.Context, roundID uint) ([]uint, error) {
	db := database.FromContext(ctx, r.db)

	var giftIDs []uint
	err := db.WithContext(ctx).
		Raw(`
			SELECT rg.user_gift_id AS id
			FROM rounds_gifts rg
			WHERE rg.round_id = $1
		`, roundID).
		Scan(&giftIDs).Error
	if len(giftIDs) == 0 {
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return giftIDs, nil
}

func (r *repo) GetUserNft(ctx context.Context, userNftID uint) (*model.Gift, error) {
	db := database.FromContext(ctx, r.db)

	var userNft *model.Gift
	err := db.WithContext(ctx).
		Raw(`
			SELECT un.id AS id, un.user_id AS user_id, c.floor AS floor
			FROM users_nfts un
				LEFT OUTER JOIN nfts n ON n.id = un.nft_id
				LEFT OUTER JOIN collections c ON n.collection_id = c.id
			WHERE un.id = $1
		`, userNftID).
		First(&userNft).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return userNft, nil
}

func (r *repo) GetUserGift(ctx context.Context, userGiftID uint) (*model.Gift, error) {
	db := database.FromContext(ctx, r.db)

	var userGift *model.Gift
	err := db.WithContext(ctx).
		Raw(`
			SELECT ug.id AS id, ug.user_id AS user_id, c.floor AS floor
			FROM users_gifts ug
				LEFT OUTER JOIN gifts g ON g.id = ug.gift_id
				LEFT OUTER JOIN collections c ON g.collection_id = c.id
			WHERE ug.id = $1
		`, userGiftID).
		Scan(&userGift).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return userGift, nil
}

func (r *repo) GetWinner(ctx context.Context, roundID uint) (*model.Winner, error) {
	db := database.FromContext(ctx, r.db)

	var winner *model.Winner
	err := db.WithContext(ctx).
		Model(&dbModels.RoundWinnerDB{}).
		Where("round_id = $1", roundID).
		First(&winner).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}
	return winner, nil
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

func (r *repo) GetUserBalance(ctx context.Context, userID uint) (int64, error) {
	db := database.FromContext(ctx, r.db)

	var balance int64
	err := db.WithContext(ctx).
		Raw(`
			SELECT balance
			FROM users
			WHERE id = $1
		`, userID).
		First(&balance).Error
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *repo) GetReferrer(ctx context.Context, refID uint) (*model.Referrer, error) {
	db := database.FromContext(ctx, r.db)

	var ref *model.Referrer
	err := db.WithContext(ctx).
		Raw(`
			SELECT u.id AS id, u.balance AS balance, u.is_spec AS is_spec
			FROM referrals r
				LEFT OUTER JOIN users u ON r.referrer_id = u.id
			WHERE r.ref_id = $1
		`, refID).
		Scan(&ref).Error
	if err != nil {
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return ref, nil
}

func (r *repo) AddRound(ctx context.Context, round *dbModels.RoundDB) error {
	db := database.FromContext(ctx, r.db)
	err := db.WithContext(ctx).
		Select("round", "secret", "hash").
		Create(round).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddUserRoundNft(ctx context.Context, roundNft *dbModels.RoundNftDB) error {
	db := database.FromContext(ctx, r.db)
	err := db.WithContext(ctx).
		Select("round_id", "user_nft_id", "bet").
		Create(roundNft).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddUserRoundGift(ctx context.Context, roundGift *dbModels.RoundGiftDB) error {
	db := database.FromContext(ctx, r.db)
	err := db.WithContext(ctx).
		Select("round_id", "user_gift_id", "bet").
		Create(roundGift).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddUserRoundTicket(ctx context.Context, roundTicket *dbModels.RoundTicketDB) error {
	db := database.FromContext(ctx, r.db)
	err := db.WithContext(ctx).
		Select("user_id", "round_id", "tickets").
		Create(roundTicket).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddWinner(ctx context.Context, winner *dbModels.RoundWinnerDB) error {
	db := database.FromContext(ctx, r.db)
	err := db.WithContext(ctx).
		Create(winner).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddReferralFee(ctx context.Context, fee *dbModels.ReferralFeeDB) error {
	db := database.FromContext(ctx, r.db)
	err := db.WithContext(ctx).
		Select("user_id", "amount", "fee").
		Create(fee).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) UpdateRoundStart(ctx context.Context, roundID uint) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Raw(`
			UPDATE rounds
			SET started_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`, roundID).
		Save(&dbModels.RoundDB{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

func (r *repo) UpdateUserNftOwner(ctx context.Context, ownerID uint, userNftID ...uint) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Model(&dbModels.UserNftDB{}).
		Where("id IN ?", userNftID).
		Updates(map[string]interface{}{
			"user_id": ownerID,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

func (r *repo) UpdateUserGiftOwner(ctx context.Context, ownerID uint, userGiftID ...uint) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Raw(`
			UPDATE users_gifts
			SET user_id = ?
			WHERE id IN (?)
		`, ownerID, userGiftID).
		Save(&dbModels.UserGiftDB{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *repo) UpdateWinner(ctx context.Context, winnerID uint, winner *dbModels.RoundWinnerDB) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Where("id = ?", winnerID).
		Updates(winner)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (r *repo) UpdateFee(ctx context.Context, userID uint) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Raw(`
			UPDATE rounds_winners
			SET is_paid = true
			WHERE user_id = $1 AND is_paid IS NULL
		`, userID).
		Save(&dbModels.RoundWinnerDB{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *repo) UpdateUserBalance(ctx context.Context, userID uint, balance int64) error {
	db := database.FromContext(ctx, r.db)

	res := db.WithContext(ctx).
		Raw(`
			UPDATE users
			SET balance = $2
			WHERE id = $1
		`, userID, balance).
		Save(&dbModels.UserDB{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
