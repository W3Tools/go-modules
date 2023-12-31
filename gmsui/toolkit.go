package gmsui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
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
	coins, err := cli.GetAllSuiCoins(context.Background(), owner)
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
func (cli *SuiClient) GetAllSuiCoins(ctx context.Context, owner string) (data []types.Coin, err error) {
	firstPage, err := cli.GetSuiCoins(ctx, owner, nil)
	if err != nil {
		return
	}
	data = append(data, firstPage.Data...)

	nextCursor := firstPage.NextCursor
	hasNext := firstPage.HasNextPage
	for hasNext {
		nextPage, err := cli.GetSuiCoins(ctx, owner, nextCursor)
		if err != nil {
			break
		}

		nextCursor = nextPage.NextCursor
		hasNext = nextPage.HasNextPage
		data = append(data, nextPage.Data...)
	}
	return
}

func (cli *SuiClient) GetSuiCoins(ctx context.Context, owner string, nextCursor *move_types.AccountAddress) (ret *types.Page[types.Coin, move_types.AccountAddress], err error) {
	ownerAddress, err := sui_types.NewAddressFromHex(owner)
	if err != nil {
		return nil, err
	}
	suiCoinType := "0x2::sui::SUI"

	return cli.Provider.GetCoins(
		ctx,
		*ownerAddress,
		&suiCoinType,
		nextCursor,
		50,
	)
}
