package models

import "time"

type UserDB struct {
	ID        uint      `gorm:"column:id"`
	Name      *string   `gorm:"column:name"`
	PhotoUrl  *string   `gorm:"column:photo_url"`
	Balance   int       `gorm:"column:balance"`
	Memo      string    `gorm:"column:memo"`
	IsSpec    *bool     `gorm:"column:is_spec"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserDB) TableName() string {
	return "users"
}

type ReferralDB struct {
	ID         uint `gorm:"column:id"`
	ReferrerID uint `gorm:"column:referrer_id"`
	RefID      uint `gorm:"column:ref_id"`
}

func (ReferralDB) TableName() string {
	return "referrals"
}

type ReferralFeeDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	Amount    int64     `gorm:"column:amount"`
	Fee       int64     `gorm:"column:fee"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
