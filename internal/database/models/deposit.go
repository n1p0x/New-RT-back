package models

import (
	"time"
)

type DepositTime uint

const (
	TonTime DepositTime = 1
)

type TonDepositDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	Amount    int       `gorm:"column:amount"`
	Payload   *string   `gorm:"column:payload"`
	MsgHash   string    `gorm:"column:msg_hash"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (TonDepositDB) TableName() string {
	return "ton_deposits"
}

type StarDepositDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	Amount    uint      `gorm:"column:amount"`
	Payload   *string   `gorm:"column:payload"`
	PaymentID string    `gorm:"column:payment_id"`
	PaidAt    time.Time `gorm:"column:paid_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (StarDepositDB) TableName() string {
	return "star_deposits"
}

type NftDepositDB struct {
	ID          uint      `gorm:"column:id"`
	UserID      uint      `gorm:"column:user_id"`
	Sender      string    `gorm:"column:sender"`
	NftAddress  string    `gorm:"column:nft_address"`
	TraceID     *string   `gorm:"column:trace_id"`
	IsConfirmed *bool     `gorm:"column:is_confirmed"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (NftDepositDB) TableName() string {
	return "nft_deposits"
}

type GiftDepositDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	GiftID    uint      `gorm:"column:gift_id"`
	MsgID     string    `gorm:"column:msg_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (GiftDepositDB) TableName() string {
	return "gift_deposits"
}

type DepositTimeDB struct {
	ID    uint      `gorm:"column:id"`
	Start time.Time `gorm:"column:start"`
}

func (DepositTimeDB) TableName() string {
	return "deposit_time"
}
