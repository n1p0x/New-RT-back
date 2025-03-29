package utils

import (
	"fmt"
	
	"github.com/xssnick/tonutils-go/address"
)

func GetAddress(addr string) (*address.Address, error) {
	rawAddr, errRaw := address.ParseRawAddr(addr)
	if errRaw != nil {
		friendlyAddr, err := address.ParseAddr(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid address %s", addr)
		}
		return friendlyAddr, nil
	}
	return rawAddr, nil
}
