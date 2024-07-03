package gmsui

import (
	"encoding/json"
	"fmt"
	"strings"

	gm "github.com/W3Tools/go-modules"
	sdk_client "github.com/W3Tools/go-sui-sdk/v2/client"
	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/go-sui-sdk/v2/types"
)

func (client *SuiClient) GetFunctionArgumentTypes(target string) ([]interface{}, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	packageId, err := sui_types.NewObjectIdFromHex(entry[0])
	if err != nil {
		return nil, err
	}

	var resp []interface{}
	return resp, client.Provider.CallContext(client.ctx, &resp, sdk_client.SuiMethod("getMoveFunctionArgTypes"), packageId, entry[1], entry[2])
}

func GetObjectAndUnmarshal[T any](client *SuiClient, id string) (raw *types.SuiObjectResponse, value *T, err error) {
	objectId, err := sui_types.NewObjectIdFromHex(id)
	if err != nil {
		return nil, nil, err
	}

	raw, err = client.Provider.GetObject(client.ctx, *objectId, &types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowBcs:                 true,
		ShowOwner:               true,
		ShowPreviousTransaction: true,
		ShowStorageRebate:       true,
		ShowDisplay:             true,
	})
	if err != nil {
		return nil, nil, err
	}

	jsb, err := json.Marshal(raw.Data.Content.Data.MoveObject.Fields)
	if err != nil {
		return nil, nil, err
	}

	value = new(T)
	err = json.Unmarshal(jsb, &value)
	if err != nil {
		return nil, nil, err
	}
	return
}

func GetObjectsAndUnmarshal[T any](client *SuiClient, ids []string) (raw []types.SuiObjectResponse, values []*T, err error) {
	objectIdArray, err := gm.Map(ids, func(v string) (move_types.AccountAddress, error) {
		id, err := sui_types.NewObjectIdFromHex(v)
		return *id, err
	})
	if err != nil {
		return nil, nil, err
	}

	raw, err = client.Provider.MultiGetObjects(client.ctx, objectIdArray, &types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowBcs:                 true,
		ShowOwner:               true,
		ShowPreviousTransaction: true,
		ShowStorageRebate:       true,
		ShowDisplay:             true,
	})
	if err != nil {
		return nil, nil, err
	}

	for _, data := range raw {
		jsb, err := json.Marshal(data.Data.Content.Data.MoveObject.Fields)
		if err != nil {
			return nil, nil, err
		}

		var value = new(T)
		err = json.Unmarshal(jsb, &value)
		if err != nil {
			return nil, nil, err
		}
		values = append(values, value)
	}
	return
}

func GetDynamicFieldObjectAndUnmarshal[T any, NameType any](client *SuiClient, parentId string, name sui_types.DynamicFieldName) (raw *types.SuiObjectResponse, value *T, err error) {
	parentIdHex, err := sui_types.NewObjectIdFromHex(parentId)
	if err != nil {
		return nil, nil, err
	}

	raw, err = client.Provider.GetDynamicFieldObject(client.ctx, *parentIdHex, name)
	if err != nil {
		return nil, nil, err
	}

	jsb, err := json.Marshal(raw.Data.Content.Data.MoveObject.Fields)
	if err != nil {
		return nil, nil, err
	}

	data := new(SuiMoveDynamicField[T, NameType])
	err = json.Unmarshal(jsb, &data)
	if err != nil {
		return nil, nil, err
	}

	return raw, &data.Value.Fields, nil
}
