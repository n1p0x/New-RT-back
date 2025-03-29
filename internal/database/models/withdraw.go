package models

import "time"

type TonWithdrawDB struct {
	ID          uint      `gorm:"column:id"`
	UserID      uint      `gorm:"column:user_id"`
	Destination string    `gorm:"column:destination"`
	Amount      uint      `gorm:"column:amount"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (TonWithdrawDB) TableName() string {
	return "ton_withdraws"
}

type NftWithdrawDB struct {
	ID          uint      `gorm:"column:id"`
	UserID      uint      `gorm:"column:user_id"`
	Destination string    `gorm:"column:destination"`
	NftAddress  string    `gorm:"column:address"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (NftWithdrawDB) TableName() string {
	return "nft_withdraws"
}

type GiftWithdrawDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	GiftID    uint      `gorm:"column:gift_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (GiftWithdrawDB) TableName() string {
	return "gift_withdraws"
}
