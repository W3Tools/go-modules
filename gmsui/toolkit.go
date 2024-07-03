package gmsui

import (
	"fmt"

	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/go-sui-sdk/v2/types"
)

// Instance Get All Sui Coins
func (client *SuiClient) GetAllCoins(owner string, coinType string) (data []types.Coin, err error) {
	firstPage, err := client.GetCoins(owner, coinType, nil)
	if err != nil {
		return
	}
	data = append(data, firstPage.Data...)

	nextCursor := firstPage.NextCursor
	hasNext := firstPage.HasNextPage
	for hasNext {
		nextPage, err := client.GetCoins(owner, coinType, nextCursor)
		if err != nil {
			break
		}

		nextCursor = nextPage.NextCursor
		hasNext = nextPage.HasNextPage
		data = append(data, nextPage.Data...)
	}
	return
}

func (client *SuiClient) GetCoins(owner, coinType string, nextCursor *move_types.AccountAddress) (ret *types.Page[types.Coin, move_types.AccountAddress], err error) {
	ownerAddress, err := sui_types.NewAddressFromHex(owner)
	if err != nil {
		return nil, err
	}

	return client.Provider.GetCoins(
		client.ctx,
		*ownerAddress,
		&coinType,
		nextCursor,
		50,
	)
}

func (cli *SuiClient) GetMaxCoinObject(address, coinType string) (*types.Coin, error) {
	coins, err := cli.GetAllCoins(address, coinType)
	if err != nil {
		return nil, fmt.Errorf("p.SuiClient.GetAllSuiCoins %v", err)
	}
	if len(coins) == 0 {
		return nil, fmt.Errorf("address: [%s], coins not found, type: %s", address, coinType)
	}

	max := coins[0]
	for _, coin := range coins {
		if coin.Balance.Uint64() > max.Balance.Uint64() {
			max = coin
		}
	}
	return &max, nil
}
