package gmsui

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
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
		balance, ok := big.NewInt(0).SetString(coin.Balance, 10)
		if !ok {
			continue
		}

		if strings.EqualFold(live, "") {
			if balance.Cmp(cli.GasBudget) == 1 {
				live = coin.CoinObjectId
			}
		}
		pending = append(pending, coin.CoinObjectId)
	}

	gas.Live = live
	gas.Pending = pending
	return nil
}

// Instance Get All Sui Coins
func (cli *SuiClient) GetAllSuiCoins(ctx context.Context, owner string) (data []models.CoinData, err error) {
	firstPage, err := cli.GetSuiCoins(ctx, owner, "")
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

func (cli *SuiClient) GetSuiCoins(ctx context.Context, owner, nextCursor string) (ret models.PaginatedCoinsResponse, err error) {
	var cursor interface{} = nil
	if !strings.EqualFold(nextCursor, "") {
		cursor = nextCursor
	}

	return cli.Provider.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
		Owner:    owner,
		CoinType: "0x2::sui::SUI",
		Cursor:   cursor,
		Limit:    50,
	})
}
