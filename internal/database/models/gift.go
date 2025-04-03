package models

import "time"

type CollectionDB struct {
	ID      uint    `gorm:"column:id"`
	Name    string  `gorm:"column:name"`
	Address *string `gorm:"column:address"`
	Floor   *int64  `gorm:"column:floor"`
}

func (CollectionDB) TableName() string {
	return "collections"
}

type NftDB struct {
	ID            uint   `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	CollectibleID uint   `gorm:"column:collectible_id"`
	Address       string `gorm:"column:address"`
	LottieUrl     string `gorm:"column:lottie_url"`
	CollectionID  uint   `gorm:"column:collection_id"`
}

func (NftDB) TableName() string {
	return "nfts"
}

type GiftDB struct {
	ID            int64  `gorm:"column:id"`
	MsgID         int64  `gorm:"column:msg_id"`
	Name          string `gorm:"column:name"`
	CollectibleID int    `gorm:"column:collectible_id"`
	LottieUrl     string `gorm:"column:lottie_url"`
	CollectionID  int    `gorm:"column:collection_id"`
}

func (GiftDB) TableName() string {
	return "gifts"
}

type UserNftDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	NftID     uint      `gorm:"column:nft_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserNftDB) TableName() string {
	return "users_nfts"
}

type UserGiftDB struct {
	ID        uint      `gorm:"column:id"`
	UserID    int64     `gorm:"column:user_id"`
	GiftID    int64     `gorm:"column:gift_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserGiftDB) TableName() string {
	return "users_gifts"
}
