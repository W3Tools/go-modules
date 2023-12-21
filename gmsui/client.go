package gmsui

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
)

type SuiClient struct {
	Provider  *client.Client
	SuiSigner *SuiSigner
	MultiSig  *SuiMultiSig
	GasBudget *big.Int
}

type SuiGasObject struct {
	Live    string
	Pending []string
}

// Create New Sui Client
func InitSuiClient(suiApi *client.Client) (client *SuiClient) {
	cli := &SuiClient{
		Provider:  suiApi,
		GasBudget: big.NewInt(2000000),
	}
	return cli
}

func (cli *SuiClient) NewSuiSigner(signer *SuiSigner) {
	if cli.SuiSigner == nil {
		cli.SuiSigner = signer
	}
	cli.updateGas(cli.SuiSigner.Signer.Address, cli.SuiSigner.Gas)
}

func (cli *SuiClient) NewSuiMultiSig(multisig *SuiMultiSig) {
	if cli.MultiSig == nil {
		cli.MultiSig = multisig
	}

	cli.updateGas(cli.MultiSig.Address, cli.MultiSig.Gas)
}

// Tools
func (cli *SuiClient) SetSignerDefaultGasObject(obj string) {
	cli.SuiSigner.Gas.Live = obj
}

func (cli *SuiClient) SetMultiSigDefaultGasObject(obj string) {
	cli.MultiSig.Gas.Live = obj
}

func (cli *SuiClient) EnableAutoUpdateGasObjectFromSigner() {
	go cli.AutoUpdateGas(cli.SuiSigner.Signer.Address, cli.SuiSigner.Gas)
}

func (cli *SuiClient) EnableAutoUpdateGasObjectFromMultiSig() {
	go cli.AutoUpdateGas(cli.MultiSig.Address, cli.MultiSig.Gas)
}

func (cli *SuiClient) SetDefaultGasBudget(budget *big.Int) {
	cli.GasBudget = budget
}

// Instance: Move Call
func (cli *SuiClient) NewMoveCall(ctx context.Context, signer, gas, target string, args []interface{}, typeArgs []string) (*types.TransactionBytes, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}

	_signer, err := sui_types.NewAddressFromHex(signer)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex(signer) %v", err)
	}

	packageId, err := sui_types.NewObjectIdFromHex(entry[0])
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewObjectIdFromHex(package) %v", err)
	}

	_gas, err := sui_types.NewObjectIdFromHex(gas)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewObjectIdFromHex(gas) %v", err)
	}

	gasBudget := types.NewSafeSuiBigInt[uint64](cli.GasBudget.Uint64())

	return cli.Provider.MoveCall(ctx, *_signer, *packageId, entry[1], entry[2], typeArgs, args, _gas, gasBudget)
}

func (cli *SuiClient) NewMoveCallFromSigner(ctx context.Context, target string, args []interface{}, typeArgs []string) (*types.TransactionBytes, error) {
	return cli.NewMoveCall(ctx, cli.SuiSigner.Signer.Address, cli.SuiSigner.Gas.Live, target, args, typeArgs)
}

func (cli *SuiClient) NewMoveCallFromMultiSig(ctx context.Context, target string, args []interface{}, typeArgs []string) (*types.TransactionBytes, error) {
	return cli.NewMoveCall(ctx, cli.MultiSig.Address, cli.MultiSig.Gas.Live, target, args, typeArgs)
}

func (cli *SuiClient) ExecuteTransaction(ctx context.Context, b64TxBytes string, signatures []any) (*types.SuiTransactionBlockResponse, error) {
	data, err := lib.NewBase64Data(b64TxBytes)
	if err != nil {
		return nil, err
	}

	return cli.Provider.ExecuteTransactionBlock(ctx, *data, signatures, &types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, types.TxnRequestTypeWaitForLocalExecution,
	)
}

func (cli *SuiClient) MoveCallFromSigner(ctx context.Context, target string, args []interface{}, typeArgs []string) (result *types.SuiTransactionBlockResponse, err error) {
	metadata, err := cli.NewMoveCall(ctx, cli.SuiSigner.Signer.Address, cli.SuiSigner.Gas.Live, target, args, typeArgs)
	if err != nil {
		return nil, fmt.Errorf("moveCall err %v", err)
	}

	signature, err := cli.SuiSigner.SignTransaction(metadata.TxBytes.String())
	if err != nil {
		return nil, fmt.Errorf("cli.SuiSigner.SignTransaction %v", err)
	}

	return cli.ExecuteTransaction(ctx, metadata.TxBytes.String(), []any{signature.Signature})
}

func (cli *SuiClient) GetObject(ctx context.Context, objectId string) (*types.SuiObjectResponse, error) {
	_objectId, err := sui_types.NewObjectIdFromHex(objectId)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewObjectIdFromHex %v", err)
	}

	return cli.Provider.GetObject(ctx, *_objectId, &types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowBcs:                 true,
		ShowOwner:               true,
		ShowPreviousTransaction: true,
		ShowStorageRebate:       true,
		ShowDisplay:             true,
	})
}
