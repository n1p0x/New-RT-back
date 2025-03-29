package service

import "errors"

var (
	ErrRoundNotFound    = errors.New("round not found")
	ErrRoundFinished    = errors.New("round is already finished")
	ErrNotEnoughBalance = errors.New("not enough balance")
)

func IsRoundNotFound(err error) bool {
	return errors.Is(err, ErrRoundNotFound)
}

func IsRoundFinished(err error) bool {
	return errors.Is(err, ErrRoundFinished)
}

func IsNotEnoughBalance(err error) bool {
	return errors.Is(err, ErrNotEnoughBalance)
}
