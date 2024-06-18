package utils

import (
	"fmt"
	"strings"
)

const (
	SuiAddressLength = 32
)

func NormalizeShortAddress(address string) string {
	return fmt.Sprintf("0x%s", strings.TrimLeft(address, "0x"))
}

func NormalizeShortCoinType(coinType string) string {
	types := strings.Split(coinType, "::")
	if len(types) != 3 {
		return coinType
	}

	return fmt.Sprintf("%s::%s::%s", NormalizeShortAddress(types[0]), types[1], types[2])
}
