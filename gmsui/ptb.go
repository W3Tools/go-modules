package gmsui

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/client"
	"github.com/W3Tools/go-modules/gmsui/types"
	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
)

type ProgrammableTransactionBlock struct {
	client  *client.SuiClient
	builder *sui_types.ProgrammableTransactionBuilder
	ctx     context.Context
}

func NewProgrammableTransactionBlock(client *client.SuiClient) *ProgrammableTransactionBlock {
	return &ProgrammableTransactionBlock{
		client:  client,
		builder: sui_types.NewProgrammableTransactionBuilder(),
		ctx:     client.Context(),
	}
}

func (ptb *ProgrammableTransactionBlock) NewMoveCall(target string, args []interface{}, typeArgs []string) (*sui_types.Argument, error) {
	arguments, err := ptb.ParseFunctionArguments(target, args)
	if err != nil {
		return nil, err
	}

	typeArguments, err := ParseFunctionTypeArguments(typeArgs)
	if err != nil {
		return nil, fmt.Errorf("parsed type arguments %v", err)
	}

	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	packageId, err := sui_types.NewAddressFromHex(entry[0])
	if err != nil {
		return nil, fmt.Errorf("sui_types.NewAddressFromHex %v", err)
	}

	returnArgument := ptb.builder.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *packageId,
				Module:        move_types.Identifier(entry[1]),
				Function:      move_types.Identifier(entry[2]),
				Arguments:     arguments,
				TypeArguments: typeArguments,
			},
		},
	)
	return &returnArgument, nil
}

// txContext should be nil
func (ptb *ProgrammableTransactionBlock) ParseFunctionArguments(target string, args []interface{}) (arguments []sui_types.Argument, err error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}

	functionArgumentTypes, err := ptb.client.GetMoveFunctionArgTypes(types.GetMoveFunctionArgTypesParams{Package: entry[0], Module: entry[1], Function: entry[2]})
	if err != nil {
		return nil, err
	}
	if len(functionArgumentTypes) > 0 && args == nil {
		return nil, fmt.Errorf("invalid arg length, required: %d, but got nil", len(functionArgumentTypes))
	}
	if len(functionArgumentTypes) != len(args) {
		return nil, fmt.Errorf("invalid arg length, required: %d, but got %d", len(functionArgumentTypes), len(args))
	}

	for idx, inputArgument := range args {
		stringType, err := json.Marshal(functionArgumentTypes[idx])
		if err != nil {
			return nil, fmt.Errorf("argument type json marshal failed %v", err)
		}

		switch string(stringType) {
		case `"Pure"`:
			var argument = sui_types.Argument{}
			switch inputArgument := inputArgument.(type) {
			case string:
				if strings.HasPrefix(inputArgument, "0x") {
					address, err := sui_types.NewAddressFromHex(inputArgument)
					if err != nil {
						return nil, fmt.Errorf("argument type to address failed %v", err)
					}
					argument, err = ptb.builder.Pure(address)
					if err != nil {
						return nil, fmt.Errorf("input argument to pure data failed %v", err)
					}
				} else {
					argument, err = ptb.builder.Pure(inputArgument)
					if err != nil {
						return nil, fmt.Errorf("input argument to pure data failed %v", err)
					}
				}
			default:
				argument, err = ptb.builder.Pure(inputArgument)
				if err != nil {
					return nil, fmt.Errorf("input argument to pure data failed %v", err)
				}
			}
			arguments = append(arguments, argument)
		case `{"Object":"ByMutableReference"}`, `{"Object":"ByImmutableReference"}`:
			if inputArgument == nil {
				continue
			}

			mutable := false
			if strings.Contains(string(stringType), "ByMutableReference") {
				mutable = true
			}
			objectInfo, _, err := GetObjectAndUnmarshal[any](ptb.client, inputArgument.(string))
			if err != nil {
				return nil, err
			}
			var objectArgs sui_types.ObjectArg
			if objectInfo.Data.Owner != nil {
				owner := *objectInfo.Data.Owner
				sharedObject, isSharedObject := owner.ObjectOwner.(types.ObjectOwner_Shared)
				if isSharedObject {
					objectId, err := sui_types.NewObjectIdFromHex(objectInfo.Data.ObjectId)
					if err != nil {
						return nil, err
					}
					objectArgs.SharedObject = &struct {
						Id                   move_types.AccountAddress
						InitialSharedVersion uint64
						Mutable              bool
					}{
						Id:                   *objectId,
						InitialSharedVersion: sharedObject.Shared.InitialSharedVersion,
						Mutable:              mutable,
					}
				} else {
					objectId, err := sui_types.NewObjectIdFromHex(objectInfo.Data.ObjectId)
					if err != nil {
						return nil, err
					}

					version, err := strconv.ParseUint(objectInfo.Data.Version, 10, 64)
					if err != nil {
						return nil, err
					}

					digest, err := sui_types.NewDigest(objectInfo.Data.Digest)
					if err != nil {
						return nil, err
					}

					objectArgs.ImmOrOwnedObject = &sui_types.ObjectRef{
						ObjectId: *objectId,
						Version:  version,
						Digest:   *digest,
					}
				}
			}

			pureData, err := ptb.builder.Obj(objectArgs)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, pureData)
		default:
			return nil, fmt.Errorf("function argument types %s not match", string(stringType))
		}
	}
	return
}

func ParseFunctionTypeArguments(typeArgs []string) (typeArguments []move_types.TypeTag, err error) {
	typeArguments = []move_types.TypeTag{}

	for _, arg := range typeArgs {
		entry := strings.Split(arg, "::")
		if len(entry) != 3 {
			return nil, fmt.Errorf("type arguments parsing failed, invalid target [%s]", arg)
		}

		typeAddress, err := sui_types.NewObjectIdFromHex(entry[0])
		if err != nil {
			return nil, fmt.Errorf("invalid package address [%v]", err)
		}

		typeTag := move_types.TypeTag{
			Struct: &move_types.StructTag{
				Address: *typeAddress,
				Module:  move_types.Identifier(entry[1]),
				Name:    move_types.Identifier(entry[2]),
			},
		}
		typeArguments = append(typeArguments, typeTag)
	}
	return
}

func (ptb *ProgrammableTransactionBlock) Finish(sender string, gasObject *string, gasBudget uint64, gasPrice *uint64) ([]byte, error) {
	hexSender, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		return nil, err
	}

	gasPayment := []*sui_types.ObjectRef{}
	if gasObject == nil {
		coins, err := ptb.client.GetCoins(types.GetCoinsParams{
			Owner:    sender,
			CoinType: gm.NewStringPtr(SuiGasCoinType),
		})
		if err != nil {
			return nil, err
		}

		for _, coin := range coins.Data {
			digest, err := sui_types.NewDigest(coin.Digest)
			if err != nil {
				return nil, err
			}

			uint64Version, err := strconv.ParseUint(coin.Version, 10, 64)
			if err != nil {
				return nil, err
			}

			objectId, err := sui_types.NewObjectIdFromHex(coin.CoinObjectId)
			if err != nil {
				return nil, err
			}
			reference := sui_types.ObjectRef{
				Digest:   *digest,
				Version:  uint64Version,
				ObjectId: *objectId,
			}
			gasPayment = append(gasPayment, &reference)
		}
	} else {
		gasObjectId, _, err := GetObjectAndUnmarshal[any](ptb.client, *gasObject)
		if err != nil {
			return nil, err
		}

		digest, err := sui_types.NewDigest(gasObjectId.Data.Digest)
		if err != nil {
			return nil, err
		}

		uint64Version, err := strconv.ParseUint(gasObjectId.Data.Version, 10, 64)
		if err != nil {
			return nil, err
		}

		objectId, err := sui_types.NewObjectIdFromHex(gasObjectId.Data.ObjectId)
		if err != nil {
			return nil, err
		}
		gasReference := sui_types.ObjectRef{
			Digest:   *digest,
			Version:  uint64Version,
			ObjectId: *objectId,
		}
		gasPayment = append(gasPayment, &gasReference)
	}

	var referenceGasPrice uint64
	if gasPrice == nil {
		refGasPrice, err := ptb.client.GetReferenceGasPrice()
		if err != nil {
			return nil, err
		}
		referenceGasPrice = refGasPrice.Uint64() + 1
	} else {
		referenceGasPrice = *gasPrice
	}

	tx := sui_types.NewProgrammable(
		*hexSender,
		gasPayment,
		ptb.builder.Finish(),
		gasBudget,
		referenceGasPrice,
	)
	return bcs.Marshal(tx)
}

func (ptb *ProgrammableTransactionBlock) Builder() *sui_types.ProgrammableTransactionBuilder {
	return ptb.builder
}

func (ptb *ProgrammableTransactionBlock) Client() *client.SuiClient {
	return ptb.client
}

func (ptb *ProgrammableTransactionBlock) Context() context.Context {
	return ptb.ctx
}
