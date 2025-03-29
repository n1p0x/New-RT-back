package handler

import (
	"roulette/internal/game/model"
)

type RoundResponse struct {
	*model.RoundWithPlayers
}

func NewRoundResponse(round *model.RoundWithPlayers) *RoundResponse {
	return &RoundResponse{
		RoundWithPlayers: round,
	}
}

type WinnerResponse struct {
	*model.Winner
}

func NewWinnerResponse(winner *model.Winner) *WinnerResponse {
	return &WinnerResponse{
		Winner: winner,
	}
}
