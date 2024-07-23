package gmsui

import (
	"github.com/W3Tools/go-modules/gmsui/client"
	"github.com/W3Tools/go-modules/gmsui/types"
	"github.com/fardream/go-bcs/bcs"
)

func DevInspect(suiClient *client.SuiClient, target string, args []interface{}, typeArgs []string) (*types.DevInspectResults, error) {
	builder := NewProgrammableTransactionBlock(suiClient)

	_, err := builder.NewMoveCall(target, args, typeArgs)
	if err != nil {
		return nil, err
	}

	tx := builder.builder.Finish()
	bs, err := bcs.Marshal(tx)
	if err != nil {
		return nil, err
	}

	txBytes := append([]byte{0}, bs...)
	return suiClient.DevInspectTransactionBlock(types.DevInspectTransactionBlockParams{
		Sender:           SuiZeroAddress,
		TransactionBlock: txBytes,
	})
}
