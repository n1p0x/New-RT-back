package model

type User struct {
	ID       uint    `json:"id"`
	Name     *string `json:"name"`
	PhotoUrl *string `json:"photoUrl"`
	Balance  int     `json:"balance"`
	Memo     string  `json:"memo"`
}

type UserProfile struct {
	Balance int64 `json:"balance"`
	Refs    uint  `json:"refs"`
	Earned  int64 `json:"earned"`
	Games   uint  `json:"games"`
}

type UpdateUser struct {
	Name     *string
	PhotoUrl *string
	Balance  *uint
}
