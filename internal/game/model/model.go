package model

import "time"

type Round struct {
	ID          uint       `json:"id"`
	RoundNumber string     `json:"-"`
	Secret      string     `json:"-"`
	Hash        string     `json:"hash"`
	CreatedAt   time.Time  `json:"createdAt"`
	StartedAt   *time.Time `json:"startedAt"`
	IsFinished  bool       `json:"-"`
}

type RoundStats struct {
	TotalGifts   int   `json:"totalGifts"`
	TotalBet     int64 `json:"totalBet"`
	TotalTickets int   `json:"totalTickets"`
}

type Player struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	RoundID   uint      `json:"roundId"`
	Tickets   int       `json:"tickets"`
	CreatedAt time.Time `json:"createdAt"`
}

type UniquePlayer struct {
	UserID   uint    `json:"userId"`
	Name     *string `json:"name"`
	PhotoUrl *string `json:"photoUrl"`
	Tickets  int     `json:"tickets"`
}

type RoundWithPlayers struct {
	*Round
	*RoundStats
	Players       []*Player       `json:"-"`
	UniquePlayers []*UniquePlayer `json:"players"`
}

type Winner struct {
	ID     uint `json:"id"`
	UserID uint `json:"userId"`
	Ticket uint `json:"ticket"`
	Fee    int  `json:"fee"`
}

type Gift struct {
	ID     uint
	UserID uint
	Floor  int64
}

type Referrer struct {
	ID      uint  `json:"id"`
	Balance int64 `json:"balance"`
	IsSpec  bool  `json:"is_spec"`
}
