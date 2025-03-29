package service

import "errors"

var (
	ErrNftsNotFound = errors.New("nfts not found")
)

func IsNftsNotFound(err error) bool {
	return errors.Is(err, ErrNftsNotFound)
}
