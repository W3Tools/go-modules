package gmsui

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
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
		GasBudget: big.NewInt(400000000),
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

func (cli *SuiClient) ImplementationOfDevInspect(ctx context.Context, txBytes string) (*types.DevInspectResults, error) {
	var accountObj *move_types.AccountAddress
	accountObj, err := sui_types.NewAddressFromHex("0x0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}

	txb, err := lib.NewBase64Data(txBytes)
	if err != nil {
		return nil, fmt.Errorf("lib.NewBase64Data %v", err)
	}

	return cli.Provider.DevInspectTransactionBlock(ctx, *accountObj, *txb, nil, nil)
}

func (cli *SuiClient) TryDevInspect(ctx context.Context, target string, args []sui_types.CallArg, typeArgs []move_types.TypeTag) (*types.DevInspectResults, error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}

	packageId, err := sui_types.NewObjectIdFromHex(entry[0])
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewObjectIdFromHex(package) %v", err)
	}

	builder := sui_types.NewProgrammableTransactionBuilder()
	err = builder.MoveCall(
		*packageId,
		move_types.Identifier(entry[1]),
		move_types.Identifier(entry[2]),
		typeArgs,
		args,
	)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewProgrammableTransactionBuilder %v", err)
	}

	transaction := builder.Finish()
	bcsBytes, err := bcs.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("bcs.Marshal %v", err)
	}

	txBytes := append([]byte{0}, bcsBytes...)
	return cli.ImplementationOfDevInspect(ctx, base64.StdEncoding.EncodeToString(txBytes))
}

func (cli *SuiClient) DevInspect(ctx context.Context, target string, args []interface{}, argsType []move_types.TypeTag) (*types.DevInspectResults, error) {
	_args, err := cli.ParseFunctionArgs(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("suiClient.ParseFunctionArgs %v", err)
	}
	return cli.TryDevInspect(ctx, target, _args, argsType)
}

func (cli *SuiClient) ParseFunctionArgs(ctx context.Context, target string, args []interface{}) (ret []sui_types.CallArg, err error) {
	functionArgTypes, err := cli.GetFunctionArgTypes(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("cli.GetFunctionArgTypes %v", err)
	}
	if len(*functionArgTypes) > 0 && args == nil {
		return nil, fmt.Errorf("invalid arg length, required: %d, but got nil", len(*functionArgTypes))
	}
	if len(*functionArgTypes) != len(args) {
		return nil, fmt.Errorf("invalid arg length, required: %d, but got %d", len(*functionArgTypes), len(args))
	}

	for idx, functionArgType := range *functionArgTypes {
		switch argType := functionArgType.(type) {
		case string:
			if argType != "Pure" {
				return nil, fmt.Errorf("unknown string type: %v", argType)
			}

			switch arg := args[idx].(type) {
			case uint8, uint16, uint32, uint64, big.Int:
				bcsArg, err := bcs.Marshal(arg)
				if err != nil {
					return nil, fmt.Errorf("givenArgs index: [%d], bcs.Marshal %v", idx, err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			case string:
				address, err := sui_types.NewAddressFromHex(arg)
				if err != nil {
					return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
				}

				bcsArg, err := bcs.Marshal(address)
				if err != nil {
					return nil, fmt.Errorf("givenArgs index: [%d], bcs.Marshal %v", idx, err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			case []uint8:
				vector := VectorU8{Data: arg}
				bcsArg, err := vector.Marshal()
				if err != nil {
					return nil, fmt.Errorf("vector.Marshal %v", err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			case []uint16:
				vector := VectorU16{Data: arg}
				bcsArg, err := vector.Marshal()
				if err != nil {
					return nil, fmt.Errorf("vector.Marshal %v", err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			case []uint32:
				vector := VectorU32{Data: arg}
				bcsArg, err := vector.Marshal()
				if err != nil {
					return nil, fmt.Errorf("vector.Marshal %v", err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			case []uint64:
				vector := VectorU64{Data: arg}
				bcsArg, err := vector.Marshal()
				if err != nil {
					return nil, fmt.Errorf("vector.Marshal %v", err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			case []big.Int:
				vector := VectorBigInt{Data: arg}
				bcsArg, err := vector.Marshal()
				if err != nil {
					return nil, fmt.Errorf("vector.Marshal %v", err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			case []string:
				vector := VectorAddress{Data: arg}
				bcsArg, err := vector.Marshal()
				if err != nil {
					return nil, fmt.Errorf("vector.Marshal %v", err)
				}
				ret = append(ret, sui_types.CallArg{Pure: &bcsArg})
			default:
				return nil, fmt.Errorf("invalid givenArgs: %v, type: %T", args[idx], args[idx])
			}
		case map[string]interface{}:
			objType, ok := argType["Object"]
			if !ok {
				return nil, fmt.Errorf("unknown map type: %v", argType)
			}

			mutable := false
			if objType == "ByMutableReference" {
				mutable = true
			}

			if fmt.Sprintf("%T", args[idx]) != "string" {
				return nil, fmt.Errorf("invalid object type, index: [%d], value: %v", idx, args[idx])
			}

			obj, err := cli.GetObject(ctx, args[idx].(string))
			if err != nil {
				return nil, fmt.Errorf("suiClient.GetObject %v", err)
			}

			var objectArgs sui_types.ObjectArg
			if obj.Data.Owner.Shared == nil {
				objectArgs.ImmOrOwnedObject = &sui_types.ObjectRef{
					ObjectId: obj.Data.ObjectId,
					Version:  obj.Data.Version.Uint64(),
					Digest:   obj.Data.Digest,
				}
			} else {
				objectArgs.SharedObject = &struct {
					Id                   move_types.AccountAddress
					InitialSharedVersion uint64
					Mutable              bool
				}{
					Id:                   obj.Data.ObjectId,
					InitialSharedVersion: *obj.Data.Owner.Shared.InitialSharedVersion,
					Mutable:              mutable,
				}
			}

			ret = append(ret, sui_types.CallArg{Object: &objectArgs})
		}
	}
	return
}

func (cli *SuiClient) NewProgrammableTransactionMoveCall(ctx context.Context, builder *sui_types.ProgrammableTransactionBuilder, target string, args []interface{}) (*sui_types.Argument, error) {
	functionArgs, err := cli.ParseFunctionArgs(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("cli.ParseFunctionArgs %v", err)
	}

	var arguments []sui_types.Argument
	for idx, functionArg := range functionArgs {
		if functionArg.Object != nil {
			argument, err := builder.Obj(*functionArg.Object)
			if err != nil {
				return nil, fmt.Errorf("givenArgs index: [%d], builder.Obj %v", idx, err)
			}
			arguments = append(arguments, argument)
			continue
		}

		argument, err := builder.Pure(args[idx])
		if err != nil {
			return nil, fmt.Errorf("givenArgs index: [%d], builder.Pure %v", idx, err)
		}
		arguments = append(arguments, argument)
	}

	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	packageId, err := sui_types.NewAddressFromHex(entry[0])
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}

	returnArgument := builder.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:   *packageId,
				Module:    move_types.Identifier(entry[1]),
				Function:  move_types.Identifier(entry[2]),
				Arguments: arguments,
			},
		},
	)
	return &returnArgument, nil
}

func (cli *SuiClient) ProgrammableTransactionFinishFromSigner(ctx context.Context, builder *sui_types.ProgrammableTransactionBuilder) ([]byte, error) {
	if cli.SuiSigner == nil {
		return nil, fmt.Errorf("invalid signer")
	}

	sender, err := sui_types.NewAddressFromHex(cli.SuiSigner.Signer.Address)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}
	gasObj, err := cli.GetObject(ctx, cli.SuiSigner.Gas.Live)
	if err != nil {
		return nil, fmt.Errorf("cli.GetObject %v", err)
	}

	gas := gasObj.Data.Reference()
	transaction := sui_types.NewProgrammable(
		*sender,
		[]*sui_types.ObjectRef{
			&gas,
		},
		builder.Finish(),
		cli.GasBudget.Uint64(),
		uint64(750),
	)
	bcsTransaction, err := bcs.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("bcs.Marshal %v", err)
	}
	return bcsTransaction, nil
}

func (cli *SuiClient) ProgrammableTransactionFinishFromMultisig(ctx context.Context, builder *sui_types.ProgrammableTransactionBuilder) ([]byte, error) {
	if cli.MultiSig == nil {
		return nil, fmt.Errorf("invalid multisig")
	}

	sender, err := sui_types.NewAddressFromHex(cli.MultiSig.Address)
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}
	gasObj, err := cli.GetObject(ctx, cli.MultiSig.Gas.Live)
	if err != nil {
		return nil, fmt.Errorf("cli.GetObject %v", err)
	}

	gas := gasObj.Data.Reference()
	transaction := sui_types.NewProgrammable(
		*sender,
		[]*sui_types.ObjectRef{
			&gas,
		},
		builder.Finish(),
		cli.GasBudget.Uint64(),
		uint64(750),
	)
	bcsTransaction, err := bcs.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("bcs.Marshal %v", err)
	}
	return bcsTransaction, nil
}
