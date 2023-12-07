package gmsui

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
)

type SuiClient struct {
	Provider  sui.ISuiAPI
	SuiSigner *SuiSigner
	MultiSig  *SuiMultiSig
	GasBudget *big.Int
}

type SuiGasObject struct {
	Live    string
	Pending []string
}

// Create New Sui Client
func InitSuiClient(suiApi sui.ISuiAPI) (client *SuiClient) {
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
	go cli.AutoUpdateGas(cli.SuiSigner.Signer.Address, cli.SuiSigner.Gas)
}

func (cli *SuiClient) NewSuiMultiSig(multisig *SuiMultiSig) {
	if cli.MultiSig == nil {
		cli.MultiSig = multisig
	}

	cli.updateGas(cli.MultiSig.Address, cli.MultiSig.Gas)
	go cli.AutoUpdateGas(cli.MultiSig.Address, cli.MultiSig.Gas)
}

// Tools
func (cli *SuiClient) SetDefaultGasBudget(budget *big.Int) {
	cli.GasBudget = budget
}

// Instance: Move Call
func (cli *SuiClient) NewMoveCall(ctx context.Context, signer, gas, target string, args []interface{}, typeArgs []interface{}) (*models.TxnMetaData, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	metadata, err := cli.Provider.MoveCall(ctx, models.MoveCallRequest{
		Signer:          signer,
		PackageObjectId: entry[0],
		Module:          entry[1],
		Function:        entry[2],
		Arguments:       args,
		TypeArguments:   typeArgs,
		Gas:             gas,
		GasBudget:       cli.GasBudget.String(),
	})
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (cli *SuiClient) MoveCallFromSigner(ctx context.Context, target string, args []interface{}, typeArgs []interface{}) (result *models.SuiTransactionBlockResponse, err error) {
	metadata, err := cli.NewMoveCall(ctx, cli.SuiSigner.Signer.Address, cli.SuiSigner.Gas.Live, target, args, typeArgs)
	if err != nil {
		return nil, fmt.Errorf("moveCall err %v", err)
	}

	ret, err := cli.Provider.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: *metadata,
		PriKey:      cli.SuiSigner.Signer.PriKey,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: "WaitForLocalExecution",
	})
	if err != nil {
		return nil, fmt.Errorf("execute err %v", err)
	}
	return &ret, err
}

func (cli *SuiClient) NewMoveCallFromSigner(ctx context.Context, target string, args []interface{}, typeArgs []interface{}) (*models.TxnMetaData, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	metadata, err := cli.Provider.MoveCall(ctx, models.MoveCallRequest{
		Signer:          cli.SuiSigner.Signer.Address,
		PackageObjectId: entry[0],
		Module:          entry[1],
		Function:        entry[2],
		Arguments:       args,
		TypeArguments:   typeArgs,
		Gas:             cli.SuiSigner.Gas.Live,
		GasBudget:       cli.GasBudget.String(),
	})
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (cli *SuiClient) NewMoveCallFromMultiSig(ctx context.Context, target string, args []interface{}, typeArgs []interface{}) (*models.TxnMetaData, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	metadata, err := cli.Provider.MoveCall(ctx, models.MoveCallRequest{
		Signer:          cli.MultiSig.Address,
		PackageObjectId: entry[0],
		Module:          entry[1],
		Function:        entry[2],
		Arguments:       args,
		TypeArguments:   typeArgs,
		Gas:             cli.MultiSig.Gas.Live,
		GasBudget:       cli.GasBudget.String(),
	})
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (cli *SuiClient) ExecuteTransaction(ctx context.Context, b64TxBytes string, signatures []string) (models.SuiTransactionBlockResponse, error) {
	return cli.Provider.SuiExecuteTransactionBlock(ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:   b64TxBytes,
		Signature: signatures,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: "WaitForLocalExecution",
	})
}
