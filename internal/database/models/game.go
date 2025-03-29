package models

import "time"

type RoundDB struct {
	ID          uint      `gorm:"column:id"`
	RoundNumber string    `gorm:"column:round"`
	Secret      string    `gorm:"column:secret"`
	Hash        string    `gorm:"column:hash"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	StartedAt   time.Time `gorm:"column:started_at"`
}

func (RoundDB) TableName() string {
	return "rounds"
}

type RoundTicketDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	RoundID   uint      `gorm:"column:round_id"`
	Tickets   int       `gorm:"column:tickets"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (RoundTicketDB) TableName() string {
	return "rounds_tickets"
}

type RoundNftDB struct {
	ID        uint  `gorm:"column:id"`
	RoundID   uint  `gorm:"column:round_id"`
	UserNftID uint  `gorm:"column:user_nft_id"`
	Bet       int64 `gorm:"column:bet"`
}

func (RoundNftDB) TableName() string {
	return "rounds_nfts"
}

type RoundGiftDB struct {
	ID         uint  `gorm:"column:id"`
	RoundID    uint  `gorm:"column:round_id"`
	UserGiftID uint  `gorm:"column:user_gift_id"`
	Bet        int64 `gorm:"column:bet"`
}

func (RoundGiftDB) TableName() string {
	return "rounds_gifts"
}

type RoundWinnerDB struct {
	ID      uint  `gorm:"column:id"`
	RoundID uint  `gorm:"column:round_id"`
	Ticket  int   `gorm:"column:ticket"`
	UserID  uint  `gorm:"column:user_id"`
	IsPaid  *bool `gorm:"column:is_paid"`
	Fee     int64 `gorm:"column:fee"`
}

func (RoundWinnerDB) TableName() string {
	return "rounds_winners"
}
