package model

type NftDeposit struct {
	ID         uint   `gorm:"id"`
	UserID     uint   `gorm:"user_id"`
	Sender     string `gorm:"sender"`
	NftAddress string `gorm:"nft_address"`
}
