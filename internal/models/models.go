package models

type Nft struct {
	Name          string `json:"name"`
	CollectibleID uint64 `json:"collectibleId"`
	Address       string `json:"address"`
	LottieUrl     string `json:"lottieUrl"`
	CollectionID  int    `json:"collectionId"`
}
