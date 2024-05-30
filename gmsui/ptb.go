package gmsui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/W3Tools/go-bcs/bcs"
	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
)

type ProgrammableTransactionBlock struct {
	client  *SuiClient
	builder *sui_types.ProgrammableTransactionBuilder
	ctx     context.Context
}

func (client *SuiClient) NewProgrammableTransactionBlock(ctx context.Context) *ProgrammableTransactionBlock {
	return &ProgrammableTransactionBlock{
		client:  client,
		builder: sui_types.NewProgrammableTransactionBuilder(),
		ctx:     ctx,
	}
}

func (ptb *ProgrammableTransactionBlock) NewMoveCall(target string, args []interface{}, typeArgs []string) (*sui_types.Argument, error) {
	arguments, err := ptb.ParseFunctionArguments(target, args)
	if err != nil {
		return nil, fmt.Errorf("cli.ParseFunctionArgs %v", err)
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
	functionArgumentTypes, err := ptb.client.GetFunctionArgumentTypes(target)
	if err != nil {
		return nil, fmt.Errorf("get function argument types failed %v", err)
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
		case `{"Object":"ByMutableReference"}`:
			if inputArgument == nil {
				continue
			}
			objectInfo, err := ptb.client.GetObject(inputArgument.(string))
			if err != nil {
				return nil, fmt.Errorf("get object %s failed %v", inputArgument, err)
			}
			var objectArgs sui_types.ObjectArg
			if objectInfo.Data.Owner.Shared == nil {
				objectArgs.ImmOrOwnedObject = &sui_types.ObjectRef{
					ObjectId: objectInfo.Data.ObjectId,
					Version:  objectInfo.Data.Version.Uint64(),
					Digest:   objectInfo.Data.Digest,
				}
			} else {
				objectArgs.SharedObject = &struct {
					Id                   move_types.AccountAddress
					InitialSharedVersion uint64
					Mutable              bool
				}{
					Id:                   objectInfo.Data.ObjectId,
					InitialSharedVersion: *objectInfo.Data.Owner.Shared.InitialSharedVersion,
					Mutable:              true,
				}
			}
			pureData, err := ptb.builder.Obj(objectArgs)
			if err != nil {
				return nil, fmt.Errorf("input argument to pure object failed %v", err)
			}
			arguments = append(arguments, pureData)
		case `{"Object":"ByImmutableReference"}`:
			if inputArgument == nil {
				continue
			}
			objectInfo, err := ptb.client.GetObject(inputArgument.(string))
			if err != nil {
				return nil, fmt.Errorf("get object %s failed %v", inputArgument, err)
			}
			var objectArgs sui_types.ObjectArg
			if objectInfo.Data.Owner.Shared == nil {
				objectArgs.ImmOrOwnedObject = &sui_types.ObjectRef{
					ObjectId: objectInfo.Data.ObjectId,
					Version:  objectInfo.Data.Version.Uint64(),
					Digest:   objectInfo.Data.Digest,
				}
			} else {
				objectArgs.SharedObject = &struct {
					Id                   move_types.AccountAddress
					InitialSharedVersion uint64
					Mutable              bool
				}{
					Id:                   objectInfo.Data.ObjectId,
					InitialSharedVersion: *objectInfo.Data.Owner.Shared.InitialSharedVersion,
					Mutable:              false,
				}
			}
			pureData, err := ptb.builder.Obj(objectArgs)
			if err != nil {
				return nil, fmt.Errorf("input argument to pure object failed %v", err)
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

func (ptb *ProgrammableTransactionBlock) FinishFromSigner() ([]byte, error) {
	if ptb.client.SuiSigner == nil {
		return nil, fmt.Errorf("finish ptb from signer failed, invalid signer")
	}
	return ptb.Finish(ptb.client.SuiSigner.Signer.Address, ptb.client.SuiSigner.Gas.Live, ptb.client.GasBudget.Uint64(), 751)
}

func (ptb *ProgrammableTransactionBlock) FinishFromMultisig() ([]byte, error) {
	if ptb.client.MultiSig == nil {
		return nil, fmt.Errorf("finish ptb from multisig failed, invalid multisig")
	}

	return ptb.Finish(ptb.client.MultiSig.Address, ptb.client.MultiSig.Gas.Live, ptb.client.GasBudget.Uint64(), 751)
}

func (ptb *ProgrammableTransactionBlock) Finish(sender, gasObject string, gasBudget, gasPrice uint64) ([]byte, error) {
	hexSender, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		return nil, fmt.Errorf("finish ptb failed, %s can not convert to address hex %v", sender, err)
	}
	gasObjectId, err := ptb.client.GetObject(gasObject)
	if err != nil {
		return nil, fmt.Errorf("finish ptb failed, get object %v", err)
	}

	gasReference := gasObjectId.Data.Reference()
	tx := sui_types.NewProgrammable(
		*hexSender,
		[]*sui_types.ObjectRef{&gasReference},
		ptb.builder.Finish(),
		gasBudget,
		gasPrice,
	)
	return bcs.Marshal(tx)
}

func (ptb *ProgrammableTransactionBlock) Builder() *sui_types.ProgrammableTransactionBuilder {
	return ptb.builder
}

func (ptb *ProgrammableTransactionBlock) Client() *SuiClient {
	return ptb.client
}

func (ptb *ProgrammableTransactionBlock) Context() context.Context {
	return ptb.ctx
}
