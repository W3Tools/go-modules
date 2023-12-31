package gmsui

import (
	"context"
	"fmt"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
)

func (cli *SuiClient) GetMoveFunctionArgTypes(ctx context.Context, packageId, module, function string) (*[]interface{}, error) {
	pkg, err := sui_types.NewAddressFromHex(packageId)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}

	var resp []interface{}
	return &resp, cli.Provider.CallContext(ctx, &resp, client.SuiMethod("getMoveFunctionArgTypes"), pkg, module, function)
}
