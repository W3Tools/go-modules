package gmsui

import (
	"encoding/json"
	"fmt"
	"strconv"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/client"
	"github.com/W3Tools/go-modules/gmsui/types"
)

func GetObjectAndUnmarshal[T any](client *client.SuiClient, id string) (raw *types.SuiObjectResponse, value *T, err error) {
	raw, err = client.GetObject(types.GetObjectParams{
		ID: id,
		Options: &types.SuiObjectDataOptions{
			ShowType:                true,
			ShowContent:             true,
			ShowBcs:                 true,
			ShowOwner:               true,
			ShowPreviousTransaction: true,
			ShowStorageRebate:       true,
			ShowDisplay:             true,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	switch raw.Data.Content.DataType {
	case types.Package:
		return nil, nil, fmt.Errorf("unimplemented %s, expected an object id, not package id", types.Package)
	case types.MoveObject:
		jsb, err := json.Marshal(raw.Data.Content.Fields)
		if err != nil {
			return nil, nil, err
		}

		value = new(T)
		err = json.Unmarshal(jsb, &value)
		if err != nil {
			return nil, nil, err
		}
		return raw, value, err
	default:
		return nil, nil, fmt.Errorf("unknown data type %s, expected an object id", raw.Data.Content.DataType)
	}
}

func GetObjectsAndUnmarshal[T any](client *client.SuiClient, ids []string) (raw []*types.SuiObjectResponse, values []*T, err error) {
	raw, err = client.MultiGetObjects(types.MultiGetObjectsParams{
		IDs: ids,
		Options: &types.SuiObjectDataOptions{
			ShowType:                true,
			ShowContent:             true,
			ShowBcs:                 true,
			ShowOwner:               true,
			ShowPreviousTransaction: true,
			ShowStorageRebate:       true,
			ShowDisplay:             true,
		},
	})

	for _, data := range raw {
		switch data.Data.Content.DataType {
		case types.Package:
			return nil, nil, fmt.Errorf("unimplemented %s, %s expected an object id, not package id", types.Package, data.Data.ObjectId)
		case types.MoveObject:
			jsb, err := json.Marshal(data.Data.Content.Fields)
			if err != nil {
				return nil, nil, err
			}

			var value = new(T)
			err = json.Unmarshal(jsb, &value)
			if err != nil {
				return nil, nil, err
			}
			values = append(values, value)
		default:
			return nil, nil, fmt.Errorf("unknown data type %s, %s expected an object id", data.Data.Content.DataType, data.Data.ObjectId)
		}
	}
	return
}

func GetDynamicFieldObjectAndUnmarshal[T any, NameType any](client *client.SuiClient, parentId string, name types.DynamicFieldName) (raw *types.SuiObjectResponse, value *SuiMoveDynamicField[T, NameType], err error) {
	raw, err = client.GetDynamicFieldObject(types.GetDynamicFieldObjectParams{
		ParentId: parentId,
		Name:     name,
	})
	if err != nil {
		return nil, nil, err
	}

	jsb, err := json.Marshal(raw.Data.Content.Fields)
	if err != nil {
		return nil, nil, err
	}

	data := new(SuiMoveDynamicField[T, NameType])
	err = json.Unmarshal(jsb, &data)
	if err != nil {
		return nil, nil, err
	}

	return raw, data, nil
}

// Instance Get All Sui Coins
func GetAllCoins(client *client.SuiClient, owner string, coinType string) (data []types.CoinStruct, err error) {
	firstPage, err := client.GetCoins(types.GetCoinsParams{Owner: owner, CoinType: &coinType, Limit: gm.NewNumberPtr(1)})
	if err != nil {
		return
	}
	data = append(data, firstPage.Data...)

	nextCursor := firstPage.NextCursor
	hasNext := firstPage.HasNextPage
	for hasNext {
		nextPage, err := client.GetCoins(types.GetCoinsParams{Owner: owner, CoinType: &coinType, Limit: gm.NewNumberPtr(1), Cursor: nextCursor})
		if err != nil {
			break
		}

		nextCursor = nextPage.NextCursor
		hasNext = nextPage.HasNextPage
		data = append(data, nextPage.Data...)
	}
	return
}

func GetMaxCoinObject(client *client.SuiClient, owner, coinType string) (*types.CoinStruct, error) {
	coins, err := GetAllCoins(client, owner, coinType)
	if err != nil {
		return nil, err
	}

	if len(coins) == 0 {
		return nil, fmt.Errorf("address: [%s], coins not found, type: %s", owner, coinType)
	}

	max := coins[0]
	for _, coin := range coins {
		balance, err := strconv.ParseUint(coin.Balance, 10, 64)
		if err != nil {
			return nil, err
		}

		maxBalance, err := strconv.ParseUint(max.Balance, 10, 64)
		if err != nil {
			return nil, err
		}

		if balance > maxBalance {
			max = coin
		}
	}
	return &max, nil
}
