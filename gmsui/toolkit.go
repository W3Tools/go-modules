package gmsui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/go-sui-sdk/v2/types"
)

func (cli *SuiClient) AutoUpdateGas(owner string, gas *SuiGasObject) {
	timer := time.NewTimer(0)
	morejobTimeout := 2 * time.Minute
	failedTimeout := 5 * time.Second
	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("exit trigger\n")
			return
		case <-timer.C:
		}

		err := cli.updateGas(owner, gas)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			timer.Reset(failedTimeout)
			continue
		}

		timer.Reset(morejobTimeout)
	}
}

func (cli *SuiClient) updateGas(owner string, gas *SuiGasObject) error {
	coinType := "0x2::sui::SUI"
	coins, err := cli.GetAllCoins(owner, coinType)
	if err != nil {
		return err
	}

	var live string
	var pending []string

	for _, coin := range coins {
		if strings.EqualFold(live, "") {
			if coin.Balance.Uint64() > cli.GasBudget.Uint64() {
				live = coin.CoinObjectId.String()
			}
		}
		pending = append(pending, coin.CoinObjectId.String())
	}

	gas.Live = live
	gas.Pending = pending
	return nil
}

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
