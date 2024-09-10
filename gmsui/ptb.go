package gmsui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/client"
	"github.com/W3Tools/go-modules/gmsui/types"
	"github.com/W3Tools/go-modules/gmsui/utils"
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
		err = fmt.Errorf("invalid target [%s]", target)
		return
	}

	normalized, err := ptb.client.GetNormalizedMoveFunction(types.GetNormalizedMoveFunctionParams{Package: entry[0], Module: entry[1], Function: entry[2]})
	if err != nil {
		return
	}

	hasTxContext := false
	if len(normalized.Parameters) > 0 && isTxContext(normalized.Parameters[len(normalized.Parameters)-1].SuiMoveNormalizedType) {
		hasTxContext = true
	}

	if hasTxContext {
		normalized.Parameters = normalized.Parameters[:len(args)]
	}

	if len(args) != len(normalized.Parameters) {
		return nil, fmt.Errorf("incorrect number of arguments")
	}

	type argumentType struct {
		Pure   any
		Object *sui_types.ObjectArg
	}
	var (
		inputArguments    = make([]*argumentType, len(normalized.Parameters))
		argumentToResolve = make(map[int]struct {
			Mutable  bool
			ObjectId string
		})
	)

	for idx, parameter := range normalized.Parameters {
		var (
			inputarg = args[idx]
		)

		pureType, ok := parameter.SuiMoveNormalizedType.(types.SuiMoveNormalizedType_String)
		if ok {
			var purevalue any
			switch pureType {
			case "Bool", "U8", "U64", "U128", "U256":
				purevalue = inputarg
			case "Address":
				var address *move_types.AccountAddress
				address, err = sui_types.NewAddressFromHex(utils.NormalizeSuiObjectId(inputarg.(string)))
				if err != nil {
					return nil, err
				}
				purevalue = address
			default:
				err = fmt.Errorf("invalid pure type %v", pureType)
				return
			}

			inputArguments[idx] = &argumentType{Pure: purevalue}
			continue
		}

		var resolve struct {
			Mutable  bool
			ObjectId string
		}

		resolve.ObjectId, ok = inputarg.(string)
		if !ok {
			err = fmt.Errorf("invalid obj")
			return
		}
		switch parameter.SuiMoveNormalizedType.(type) {
		case types.SuiMoveNormalizedType_MutableReference:
			resolve.Mutable = true
		}
		argumentToResolve[idx] = resolve
	}

	var ids []string
	for _, resolve := range argumentToResolve {
		ids = append(ids, resolve.ObjectId)
	}

	if len(ids) == 0 {
		return
	}

	var objects []*types.SuiObjectResponse
	objects, err = ptb.client.MultiGetObjects(types.MultiGetObjectsParams{IDs: ids, Options: &types.SuiObjectDataOptions{ShowOwner: true}})
	if err != nil {
		return
	}

	for idx, resolveObject := range argumentToResolve {
		object := gm.FilterOne(objects, func(v *types.SuiObjectResponse) bool {
			return v.Data.ObjectId == utils.NormalizeSuiObjectId(resolveObject.ObjectId)
		})
		if object == nil {
			err = fmt.Errorf("object not found")
			return
		}

		var (
			objecrArgument sui_types.ObjectArg
			objectId       *move_types.AccountAddress
		)

		objectId, err = sui_types.NewObjectIdFromHex(object.Data.ObjectId)
		if err != nil {
			return
		}

		switch t := object.Data.Owner.ObjectOwner.(type) {
		case types.ObjectOwner_Shared:
			objecrArgument.SharedObject = &struct {
				Id                   move_types.AccountAddress
				InitialSharedVersion uint64
				Mutable              bool
			}{
				Id:                   *objectId,
				InitialSharedVersion: t.Shared.InitialSharedVersion,
				Mutable:              resolveObject.Mutable,
			}
		default:
			version, err := strconv.ParseUint(object.Data.Version, 10, 64)
			if err != nil {
				return nil, err
			}

			digest, err := sui_types.NewDigest(object.Data.Digest)
			if err != nil {
				return nil, err
			}

			objecrArgument.ImmOrOwnedObject = &sui_types.ObjectRef{
				ObjectId: *objectId,
				Version:  version,
				Digest:   *digest,
			}
		}
		inputArguments[idx] = &argumentType{Object: &objecrArgument}
	}

	arguments, _ = gm.Map(inputArguments, func(v *argumentType) (sui_types.Argument, error) {
		if v.Pure != nil {
			return ptb.builder.Pure(v.Pure)
		}

		if v.Object != nil {
			return ptb.builder.Obj(*v.Object)
		}

		return sui_types.Argument{}, fmt.Errorf("invalid argument")
	})
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

func (ptb *ProgrammableTransactionBlock) Finish(sender string, gasObject *string, gasBudget uint64, gasPrice *uint64) (*sui_types.TransactionData, []byte, error) {
	hexSender, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		return nil, nil, err
	}

	gasPayment := []*sui_types.ObjectRef{}
	if gasObject == nil {
		coins, err := ptb.client.GetCoins(types.GetCoinsParams{
			Owner:    sender,
			CoinType: gm.NewStringPtr(SuiGasCoinType),
		})
		if err != nil {
			return nil, nil, err
		}

		for _, coin := range coins.Data {
			digest, err := sui_types.NewDigest(coin.Digest)
			if err != nil {
				return nil, nil, err
			}

			uint64Version, err := strconv.ParseUint(coin.Version, 10, 64)
			if err != nil {
				return nil, nil, err
			}

			objectId, err := sui_types.NewObjectIdFromHex(coin.CoinObjectId)
			if err != nil {
				return nil, nil, err
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
			return nil, nil, err
		}

		digest, err := sui_types.NewDigest(gasObjectId.Data.Digest)
		if err != nil {
			return nil, nil, err
		}

		uint64Version, err := strconv.ParseUint(gasObjectId.Data.Version, 10, 64)
		if err != nil {
			return nil, nil, err
		}

		objectId, err := sui_types.NewObjectIdFromHex(gasObjectId.Data.ObjectId)
		if err != nil {
			return nil, nil, err
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
			return nil, nil, err
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
	bs, err := bcs.Marshal(tx)
	return &tx, bs, err
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

// --------
func isTxContext(param types.SuiMoveNormalizedType) bool {
	structType := extractStructTag(param)
	if structType == nil {
		return false
	}

	return structType.Struct.Address == "0x2" && structType.Struct.Module == "tx_context" && structType.Struct.Name == "TxContext"
}

func extractStructTag(normalizedType types.SuiMoveNormalizedType) *types.SuiMoveNormalizedType_Struct {
	_struct, ok := normalizedType.(types.SuiMoveNormalizedType_Struct)
	if ok {
		return &_struct
	}

	ref := extractReference(normalizedType)
	mutRef := extractMutableReference(normalizedType)

	if ref != nil {
		return extractStructTag(ref)
	}

	if mutRef != nil {
		return extractStructTag(mutRef)
	}

	return nil
}

func extractReference(normalizedType types.SuiMoveNormalizedType) types.SuiMoveNormalizedType {
	reference, ok := normalizedType.(types.SuiMoveNormalizedType_Reference)
	if ok {
		return reference.Reference.SuiMoveNormalizedType
	}
	return nil
}

func extractMutableReference(normalizedType types.SuiMoveNormalizedType) types.SuiMoveNormalizedType {
	mutableReference, ok := normalizedType.(types.SuiMoveNormalizedType_MutableReference)
	if ok {
		return mutableReference.MutableReference.SuiMoveNormalizedType
	}
	return nil
}
