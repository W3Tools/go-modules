package gmsui

import (
	"context"
	"fmt"
	"strings"

	"github.com/W3Tools/go-sui-sdk/v2/client"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
)

func (cli *SuiClient) GetFunctionArgTypes(ctx context.Context, target string) (*[]interface{}, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}

	return cli.GetMoveFunctionArgTypes(ctx, entry[0], entry[1], entry[2])
}

func (cli *SuiClient) GetMoveFunctionArgTypes(ctx context.Context, packageId, module, function string) (*[]interface{}, error) {
	pkg, err := sui_types.NewAddressFromHex(packageId)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}

	var resp []interface{}
	return &resp, cli.Provider.CallContext(ctx, &resp, client.SuiMethod("getMoveFunctionArgTypes"), pkg, module, function)
}
