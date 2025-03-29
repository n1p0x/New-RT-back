package handler

import (
	"roulette/internal/gift/model"
	"roulette/internal/models"
)

type CollectionsResponse struct {
	Collections []*model.Collection
}

func NewCollectionsResponse(collections []*model.Collection) *CollectionsResponse {
	return &CollectionsResponse{Collections: collections}
}

type NftsResponse struct {
	Nfts []*models.Nft `json:"nfts"`
}

func NewNftsResponse(nfts []*models.Nft) *NftsResponse {
	return &NftsResponse{Nfts: nfts}
}

type UserGiftsResponse struct {
	*model.UserGifts
	IsAvailable bool  `json:"isAvailable"`
	Fee         int64 `json:"fee"`
}

func NewUserGiftsResponse(userGifts *model.UserGifts, isAvailable bool, fee int64) *UserGiftsResponse {
	return &UserGiftsResponse{UserGifts: userGifts, IsAvailable: isAvailable, Fee: fee}
}
