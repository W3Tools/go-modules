package gmsui

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	sdk_client "github.com/W3Tools/go-sui-sdk/v2/client"
	"github.com/W3Tools/go-sui-sdk/v2/lib"
	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/go-sui-sdk/v2/types"
	"github.com/fardream/go-bcs/bcs"
)

type SuiClient struct {
	ctx       context.Context
	Provider  *sdk_client.Client
	Keypair   cryptography.Keypair
	MultiSig  *SuiMultiSig
	GasBudget *big.Int
}

type SuiGasObject struct {
	Live    string
	Pending []string
}

// Create New Sui Client
func InitClient(ctx context.Context, suiApi *sdk_client.Client) *SuiClient {
	return &SuiClient{
		ctx:       ctx,
		Provider:  suiApi,
		GasBudget: big.NewInt(400000000),
	}
}

func (client *SuiClient) Context() context.Context {
	return client.ctx
}

func (client *SuiClient) SetKeypair(keypair cryptography.Keypair) {
	client.Keypair = keypair
}

func (client *SuiClient) NewMultiSig(multisig *SuiMultiSig) {
	if client.MultiSig == nil {
		client.MultiSig = multisig
	}

	client.updateGas(client.MultiSig.Address, client.MultiSig.Gas)
}

func (client *SuiClient) SetDefaultGasObjectToMultiSig(obj string) {
	client.MultiSig.Gas.Live = obj
}

func (client *SuiClient) EnableAutoUpdateGasObjectFromMultiSig() {
	go client.AutoUpdateGas(client.MultiSig.Address, client.MultiSig.Gas)
}

func (client *SuiClient) SetDefaultGasBudget(budget *big.Int) {
	client.GasBudget = budget
}

// Instance: Move Call
func (client *SuiClient) NewMoveCall(signer string, gas *string, target string, args []interface{}, typeArgs []string) (*types.TransactionBytes, error) {
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

	var referenceGas *move_types.AccountAddress
	if gas == nil {
		referenceGas, err = client.getMaxGasCoin(signer)
		if err != nil {
			return nil, err
		}
	} else {
		referenceGas, err = sui_types.NewObjectIdFromHex(*gas)
		if err != nil {
			return nil, err
		}
	}

	gasBudget := types.NewSafeSuiBigInt[uint64](client.GasBudget.Uint64())

	return client.Provider.MoveCall(client.ctx, *_signer, *packageId, entry[1], entry[2], typeArgs, args, referenceGas, gasBudget)
}

func (client *SuiClient) NewMoveCallFromSigner(target string, args []interface{}, typeArgs []string) (*types.TransactionBytes, error) {
	return client.NewMoveCall(client.Keypair.ToSuiAddress(), nil, target, args, typeArgs)
}

func (client *SuiClient) NewMoveCallFromMultiSig(target string, args []interface{}, typeArgs []string) (*types.TransactionBytes, error) {
	return client.NewMoveCall(client.MultiSig.Address, nil, target, args, typeArgs)
}

func (client *SuiClient) ExecuteTransaction(b64TxBytes string, signatures []any) (*types.SuiTransactionBlockResponse, error) {
	data, err := lib.NewBase64Data(b64TxBytes)
	if err != nil {
		return nil, err
	}

	return client.Provider.ExecuteTransactionBlock(client.ctx, *data, signatures, &types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, types.TxnRequestTypeWaitForLocalExecution,
	)
}

func (client *SuiClient) MoveCallFromSigner(target string, args []interface{}, typeArgs []string) (result *types.SuiTransactionBlockResponse, err error) {
	metadata, err := client.NewMoveCall(client.Keypair.ToSuiAddress(), nil, target, args, typeArgs)
	if err != nil {
		return nil, err
	}

	signature, err := client.Keypair.SignTransactionBlock(metadata.TxBytes.Data())
	if err != nil {
		return nil, err
	}

	return client.ExecuteTransaction(metadata.TxBytes.String(), []any{signature.Signature})
}

func (client *SuiClient) ImplementationOfDevInspect(txBytes string) (*types.DevInspectResults, error) {
	var accountObj *move_types.AccountAddress
	accountObj, err := sui_types.NewAddressFromHex("0x0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}

	txb, err := lib.NewBase64Data(txBytes)
	if err != nil {
		return nil, fmt.Errorf("lib.NewBase64Data %v", err)
	}

	return client.Provider.DevInspectTransactionBlock(client.ctx, *accountObj, *txb, nil, nil)
}

func (client *SuiClient) DevInspect(target string, args []interface{}, typeArgs []string) (*types.DevInspectResults, error) {
	ptb := client.NewProgrammableTransactionBlock(client.ctx)
	_, err := ptb.NewMoveCall(target, args, typeArgs)
	if err != nil {
		return nil, fmt.Errorf("new move call failed %v", err)
	}

	tx := ptb.builder.Finish()
	bcsBytes, err := bcs.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("dev inspect failed, bcs marshal %v", err)
	}
	txBytes := append([]byte{0}, bcsBytes...)
	return client.ImplementationOfDevInspect(base64.StdEncoding.EncodeToString(txBytes))
}

func ParseDevInspectReturnValue(v interface{}) []byte {
	returnValues := v.([]any)
	values := returnValues[0].([]interface{})
	bs := make([]byte, 0)
	for _, i := range values {
		float64Value := i.(float64)
		bs = append(bs, uint8(float64Value))
	}
	return bs
}

func (client *SuiClient) getMaxGasCoin(address string) (*sui_types.ObjectID, error) {
	addressHex, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return nil, err
	}

	result, err := client.Provider.GetCoins(client.ctx, *addressHex, &SuiGasCoinType, nil, 50)
	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("gas coin not found")
	}

	maxGas := result.Data[0]
	for _, c := range result.Data {
		if c.Balance.Uint64() > maxGas.Balance.Uint64() {
			maxGas = c
		}
	}
	return &maxGas.CoinObjectId, nil
}
