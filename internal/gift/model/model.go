package model

type Collection struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Address *string `json:"address"`
	Floor   int64   `json:"floor"`
	//ImgUrl  string `json:"imgUrl"`
}

type UserNft struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"user_id"`
	NftID   uint   `json:"nft_id"`
	Address string `json:"address"`
}

type UserGift struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	GiftID uint `json:"gift_id"`
	MsgID  int  `json:"message_id"`
}

// Gift used for users gifts request
type Gift struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	CollectibleID uint   `json:"collectibleId"`
	LottieUrl     string `json:"lottieUrl"`
	Floor         int64  `json:"floor"`
	IsBet         bool   `json:"isBet"`
}

type UserGifts struct {
	Nfts  []*Gift `json:"nfts"`
	Gifts []*Gift `json:"gifts"`
}
