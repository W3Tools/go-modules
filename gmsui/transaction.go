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
	"github.com/W3Tools/go-sui-sdk/v2/lib"
	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
)

type Transaction struct {
	client  *client.SuiClient
	builder *sui_types.ProgrammableTransactionBuilder
	ctx     context.Context
}

func NewTransaction(client *client.SuiClient) *Transaction {
	return &Transaction{
		client:  client,
		builder: sui_types.NewProgrammableTransactionBuilder(),
		ctx:     client.Context(),
	}
}

type TransactionInputGasCoin struct {
	GasCoin bool `json:"gasCoin"`
}

func (txb *Transaction) Gas() *TransactionInputGasCoin {
	return &TransactionInputGasCoin{GasCoin: true}
}

func (txb *Transaction) SplitCoins(coin interface{}, amounts []interface{}) (returnArguments []*sui_types.Argument, err error) {
	if len(amounts) == 0 {
		return nil, fmt.Errorf("got empty amounts")
	}

	var inputCoin sui_types.Argument
	switch coin := coin.(type) {
	case *TransactionInputGasCoin:
		inputCoin = sui_types.Argument{GasCoin: &lib.EmptyEnum{}}
	case sui_types.Argument:
		inputCoin = coin
	case string:
		address, err := sui_types.NewAddressFromHex(utils.NormalizeSuiAddress(coin))
		if err != nil {
			return nil, fmt.Errorf("invalid address [%v]", err)
		}

		inputCoin, err = txb.builder.Pure(address)
		if err != nil {
			return nil, fmt.Errorf("failed to create pure argument, err: %v", err)
		}
	default:
		return nil, fmt.Errorf("invalie input coin type, got %T", coin)
	}

	amountArguments := make([]sui_types.Argument, len(amounts))
	for i, amount := range amounts {
		switch amount := amount.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
			amountArguments[i], err = txb.builder.Pure(amount.(uint64))
			if err != nil {
				return nil, fmt.Errorf("failed to create pure argument, err: %v", err)
			}
		case sui_types.Argument:
			amountArguments[i] = amount
		default:
			return nil, fmt.Errorf("invalid amount type, type: %T, value: %v", amount, amount)
		}
	}

	txb.builder.Command(
		sui_types.Command{
			SplitCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{
				Argument:  inputCoin,
				Arguments: amountArguments,
			},
		},
	)

	return txb.createTransactionResult(len(amounts)), nil
}

func (txb *Transaction) NewMoveCall(target string, args []interface{}, typeArgs []string) (returnArguments []*sui_types.Argument, err error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	var pkg, mod, fn = utils.NormalizeSuiObjectId(entry[0]), entry[1], entry[2]

	arguments, returnsCount, err := txb.ParseFunctionArguments(pkg, mod, fn, args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function arguments, err: %v", err)
	}

	typeArguments, err := ParseFunctionTypeArguments(typeArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function type arguments, err: %v", err)
	}

	packageId, err := sui_types.NewAddressFromHex(pkg)
	if err != nil {
		return nil, fmt.Errorf("invalid package address [%v]", err)
	}

	txb.builder.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *packageId,
				Module:        move_types.Identifier(mod),
				Function:      move_types.Identifier(fn),
				Arguments:     arguments,
				TypeArguments: typeArguments,
			},
		},
	)

	return txb.createTransactionResult(returnsCount), nil
}

func (txb *Transaction) createTransactionResult(count int) []*sui_types.Argument {
	nestedResult1 := uint16(len(txb.builder.Commands) - 1)
	returnArguments := make([]*sui_types.Argument, count)
	for i := 0; i < count; i++ {
		returnArguments[i] = &sui_types.Argument{
			NestedResult: &struct {
				Result1 uint16
				Result2 uint16
			}{
				Result1: nestedResult1,
				Result2: uint16(i),
			},
		}
	}

	return returnArguments
}

func (txb *Transaction) ParseFunctionArguments(pkg, mod, fn string, args []interface{}) (arguments []sui_types.Argument, returnsCount int, err error) {
	normalized, err := txb.client.GetNormalizedMoveFunction(types.GetNormalizedMoveFunctionParams{Package: pkg, Module: mod, Function: fn})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get normalized move function, err: %v", err)
	}

	hasTxContext := false
	if len(normalized.Parameters) > 0 && isTxContext(normalized.Parameters[len(normalized.Parameters)-1].SuiMoveNormalizedType) {
		hasTxContext = true
	}

	if hasTxContext {
		normalized.Parameters = normalized.Parameters[:len(args)]
	}

	if len(args) != len(normalized.Parameters) {
		return nil, 0, fmt.Errorf("incorrect number of arguments")
	}

	type argumentType struct {
		Pure         any
		Object       *sui_types.ObjectArg
		typeArgument *sui_types.Argument
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

		byvalue, ok := inputarg.(*sui_types.Argument)
		if ok {
			inputArguments[idx] = &argumentType{typeArgument: byvalue}
			continue
		}

		_, ok = parameter.SuiMoveNormalizedType.(types.SuiMoveNormalizedType_Vector)
		if ok {
			inputArguments[idx] = &argumentType{Pure: inputarg}
			continue
		}

		pureType, ok := parameter.SuiMoveNormalizedType.(types.SuiMoveNormalizedType_String)
		if ok {
			var purevalue any
			switch pureType {
			case "Bool", "U8", "U16", "U32", "U64", "U128", "U256":
				purevalue = inputarg
			case "Address":
				var address *move_types.AccountAddress
				address, err = sui_types.NewAddressFromHex(utils.NormalizeSuiAddress(inputarg.(string)))
				if err != nil {
					return nil, 0, err
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
			return nil, 0, fmt.Errorf("invalid object [%s]", inputarg)
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

	if len(ids) > 0 {
		var objects []*types.SuiObjectResponse
		objects, err = txb.client.MultiGetObjects(types.MultiGetObjectsParams{IDs: ids, Options: &types.SuiObjectDataOptions{ShowOwner: true}})
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
					return nil, 0, err
				}

				digest, err := sui_types.NewDigest(object.Data.Digest)
				if err != nil {
					return nil, 0, err
				}

				objecrArgument.ImmOrOwnedObject = &sui_types.ObjectRef{
					ObjectId: *objectId,
					Version:  version,
					Digest:   *digest,
				}
			}
			inputArguments[idx] = &argumentType{Object: &objecrArgument}
		}
	}

	arguments, _ = gm.Map(inputArguments, func(v *argumentType) (sui_types.Argument, error) {
		if v.Pure != nil {
			return txb.builder.Pure(v.Pure)
		}

		if v.Object != nil {
			return txb.builder.Obj(*v.Object)
		}

		if v.typeArgument != nil {
			return *v.typeArgument, nil
		}

		return sui_types.Argument{}, fmt.Errorf("invalid argument")
	})
	return arguments, len(normalized.Return), nil
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

func (txb *Transaction) Finish(sender string, gasObject *string, gasBudget uint64, gasPrice *uint64) (*sui_types.TransactionData, []byte, error) {
	hexSender, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		return nil, nil, err
	}

	gasPayment := []*sui_types.ObjectRef{}
	if gasObject == nil {
		coins, err := txb.client.GetCoins(types.GetCoinsParams{
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
		gasObjectId, _, err := GetObjectAndUnmarshal[any](txb.client, *gasObject)
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
		refGasPrice, err := txb.client.GetReferenceGasPrice()
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
		txb.builder.Finish(),
		gasBudget,
		referenceGasPrice,
	)
	bs, err := bcs.Marshal(tx)
	return &tx, bs, err
}

func (txb *Transaction) Builder() *sui_types.ProgrammableTransactionBuilder {
	return txb.builder
}

func (txb *Transaction) Client() *client.SuiClient {
	return txb.client
}

func (txb *Transaction) Context() context.Context {
	return txb.ctx
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
