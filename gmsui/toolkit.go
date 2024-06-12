package gmsui

import (
	"context"
	"encoding/json"
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

func GetObjectAndUnmarshal[T any](client *SuiClient, objectId string) (rawData *types.SuiObjectResponse, value *T, err error) {
	rawData, err = client.GetObject(objectId)
	if err != nil {
		return nil, nil, fmt.Errorf("get object -> %v", err)
	}

	jsb, err := json.Marshal(rawData.Data.Content.Data.MoveObject.Fields)
	if err != nil {
		return nil, nil, fmt.Errorf("json marshal -> %v", err)
	}

	value = new(T)
	err = json.Unmarshal(jsb, &value)
	if err != nil {
		return nil, nil, fmt.Errorf("json unmarshal -> %v", err)
	}
	return
}

func GetObjectsAndUnmarshal[T any](client *SuiClient, objectIds []string) (rawData []types.SuiObjectResponse, values []*T, err error) {
	rawData, err = client.GetObjects(objectIds)
	if err != nil {
		return nil, nil, fmt.Errorf("get objects -> %v", err)
	}

	for _, data := range rawData {
		jsb, err := json.Marshal(data.Data.Content.Data.MoveObject.Fields)
		if err != nil {
			return nil, nil, fmt.Errorf("json marshal -> %v", err)
		}

		var value = new(T)
		err = json.Unmarshal(jsb, &value)
		if err != nil {
			return nil, nil, fmt.Errorf("json unmarshal -> %v", err)
		}
		values = append(values, value)
	}
	return
}

func GetDynamicFieldObjectAndUnmarshal[T any](client *SuiClient, parentId string, name sui_types.DynamicFieldName) (rawData *types.SuiObjectResponse, value *T, err error) {
	parentIdHex, err := sui_types.NewObjectIdFromHex(parentId)
	if err != nil {
		return nil, nil, fmt.Errorf("new object id from hex -> %v", err)
	}

	rawData, err = client.Provider.GetDynamicFieldObject(client.ctx, *parentIdHex, name)
	if err != nil {
		return nil, nil, fmt.Errorf("get dynamic field object -> %v", err)
	}

	jsb, err := json.Marshal(rawData.Data.Content.Data.MoveObject.Fields)
	if err != nil {
		return nil, nil, fmt.Errorf("json marshal -> %v", err)
	}

	data := new(SuiMoveDynamicField[T])
	err = json.Unmarshal(jsb, &data)
	if err != nil {
		return nil, nil, fmt.Errorf("json unmarshal -> %v", err)
	}

	return rawData, &data.Value.Fields, nil
}
