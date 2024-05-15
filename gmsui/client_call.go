package gmsui

import (
	"context"
	"fmt"
	"strings"

	sdk_client "github.com/W3Tools/go-sui-sdk/v2/client"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
)

func (client *SuiClient) GetFunctionArgumentTypes(ctx context.Context, target string) ([]interface{}, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	packageId, err := sui_types.NewObjectIdFromHex(entry[0])
	if err != nil {
		return nil, fmt.Errorf("get function argument types failed, invalid package id [%s], %v", entry[0], err)
	}

	var resp []interface{}
	return resp, client.Provider.CallContext(ctx, &resp, sdk_client.SuiMethod("getMoveFunctionArgTypes"), packageId, entry[1], entry[2])
}
