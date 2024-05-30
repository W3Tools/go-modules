package gmsui

import (
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
		return nil, fmt.Errorf("get function argument types failed, invalid package id [%s], %v", entry[0], err)
	}

	var resp []interface{}
	return resp, client.Provider.CallContext(client.ctx, &resp, sdk_client.SuiMethod("getMoveFunctionArgTypes"), packageId, entry[1], entry[2])
}

// GET OBJECT/MULTI_OBJECTS
func (client *SuiClient) GetObject(objectId string) (*types.SuiObjectResponse, error) {
	_objectId, err := sui_types.NewObjectIdFromHex(objectId)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewObjectIdFromHex %v", err)
	}

	return client.Provider.GetObject(client.ctx, *_objectId, &types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowBcs:                 true,
		ShowOwner:               true,
		ShowPreviousTransaction: true,
		ShowStorageRebate:       true,
		ShowDisplay:             true,
	})
}

func (client *SuiClient) GetObjects(objectIds []string) ([]types.SuiObjectResponse, error) {
	ids, err := gm.Map(objectIds, func(v string) (move_types.AccountAddress, error) {
		hex, err := sui_types.NewObjectIdFromHex(v)
		return *hex, err
	})
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewObjectIdFromHex %v", err)
	}

	return client.Provider.MultiGetObjects(client.ctx, ids, &types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowBcs:                 true,
		ShowOwner:               true,
		ShowPreviousTransaction: true,
		ShowStorageRebate:       true,
		ShowDisplay:             true,
	})
}
