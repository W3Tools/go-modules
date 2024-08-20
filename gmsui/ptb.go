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

// Constants
const (
	MAX_PURE_ARGUMENT_SIZE = 16 * 1024
	MAX_TX_GAS             = 50_000_000_000
	MAX_GAS_OBJECTS        = 256
	MAX_TX_SIZE_BYTES      = 128 * 1024
	GAS_SAFE_OVERHEAD      = 1000
)

// Transaction Building Parameters Interface Type
type BuildOptions struct {
	OnlyTransactionKind bool `json:"onlyTransactionKind"`
}

type buildOption struct {
	Overrides           buildOptionOverride `json:"Overrides"`
	OnlyTransactionKind bool                `json:"onlyTransactionKind"`
}

type buildOptionOverride struct {
	Sender    string    `json:"sender"`
	GasConfig GasConfig `json:"gasConfig"`
	// Expiration any       `json:"expiration"`
}

type GasConfig struct {
	Budget  uint64                 `json:"budget"`
	Price   uint64                 `json:"price"`
	Payment []*sui_types.ObjectRef `json:"payment"`
	Owner   string                 `json:"owner"`
}

type ProgrammableTransactionBlock struct {
	client  *client.SuiClient
	builder *sui_types.ProgrammableTransactionBuilder
	ctx     context.Context

	sender    string
	gasConfig GasConfig
}

func NewProgrammableTransactionBlock(client *client.SuiClient) *ProgrammableTransactionBlock {
	return &ProgrammableTransactionBlock{
		client:  client,
		builder: sui_types.NewProgrammableTransactionBuilder(),
		ctx:     client.Context(),
	}
}

// The Programmable Transaction Block Parameters Settings
func (ptb *ProgrammableTransactionBlock) SetSender(sender string) {
	ptb.sender = sender
}

func (ptb *ProgrammableTransactionBlock) SetSenderIfNotSet(sender string) {
	if ptb.sender == "" {
		ptb.sender = sender
	}
}

func (ptb *ProgrammableTransactionBlock) SetGasPrice(price uint64) {
	ptb.gasConfig.Price = price
}

func (ptb *ProgrammableTransactionBlock) SetGasBudget(budget uint64) {
	ptb.gasConfig.Budget = budget
}

func (ptb *ProgrammableTransactionBlock) SetGasOwner(owner string) {
	ptb.gasConfig.Owner = owner
}

func (ptb *ProgrammableTransactionBlock) SetGasPayment(payments []*sui_types.ObjectRef) {
	ptb.gasConfig.Payment = payments
}

// The Programmable Transaction Block Parameter Retrieval
func (ptb *ProgrammableTransactionBlock) Builder() *sui_types.ProgrammableTransactionBuilder {
	return ptb.builder
}

func (ptb *ProgrammableTransactionBlock) Client() *client.SuiClient {
	return ptb.client
}

func (ptb *ProgrammableTransactionBlock) Context() context.Context {
	return ptb.ctx
}

// Transaction Option
func (ptb *ProgrammableTransactionBlock) NewMergeCoins(distination string, sources []string) (*sui_types.Argument, error) {
	if len(sources) == 0 || distination == "" {
		return nil, fmt.Errorf("missing distination coin or sources coins")
	}

	coinObjects, err := ptb.client.MultiGetObjects(
		types.MultiGetObjectsParams{
			IDs: append([]string{distination}, sources...),
			Options: &types.SuiObjectDataOptions{
				ShowContent: true,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if len(coinObjects) != len(sources)+1 {
		return nil, fmt.Errorf("invalid coin list or coin duplicates")
	}

	var distinationArgument sui_types.Argument
	var sourceArguments []sui_types.Argument
	for _, object := range coinObjects {
		if object.Error != nil {
			ebytes, err := object.Error.MarshalJSON()
			return nil, fmt.Errorf("invalid coin object, cause: %v, marshal error: %v", string(ebytes), err)
		}

		objectId, err := sui_types.NewObjectIdFromHex(object.Data.ObjectId)
		if err != nil {
			return nil, err
		}

		version, err := strconv.ParseUint(object.Data.Version, 10, 64)
		if err != nil {
			return nil, err
		}

		digest, err := sui_types.NewDigest(object.Data.Digest)
		if err != nil {
			return nil, err
		}

		arg, err := ptb.builder.Obj(
			sui_types.ObjectArg{
				ImmOrOwnedObject: &sui_types.ObjectRef{
					ObjectId: *objectId,
					Version:  version,
					Digest:   *digest,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		if utils.NormalizeSuiObjectId(object.Data.ObjectId) == utils.NormalizeSuiObjectId(distination) {
			distinationArgument = arg
		} else {
			sourceArguments = append(sourceArguments, arg)
		}
	}

	argument := ptb.builder.Command(
		sui_types.Command{
			MergeCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{Argument: distinationArgument, Arguments: sourceArguments},
		},
	)
	return &argument, nil
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

type ArgumentType string

const (
	Pure   ArgumentType = "Pure"
	Object ArgumentType = "Object"
)

type ResolveArgument struct {
	Type           ArgumentType                `json:"type"`
	PureValue      interface{}                 `json:"pureValue"`
	ObjectValue    sui_types.ObjectArg         `json:"objectValue"`
	Id             string                      `json:"id"`
	NormalizedType types.SuiMoveNormalizedType `json:"normalizedType"`
}

func (ptb *ProgrammableTransactionBlock) ParseFunctionArguments(target string, args []interface{}) (arguments []sui_types.Argument, err error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}

	{
		normalized, err := ptb.client.GetNormalizedMoveFunction(
			types.GetNormalizedMoveFunctionParams{
				Package:  utils.NormalizeSuiObjectId(entry[0]),
				Module:   entry[1],
				Function: entry[2],
			},
		)
		if err != nil {
			return nil, err
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

		argumentToResolve := []*ResolveArgument{}
		for idx, param := range normalized.Parameters {
			arg := args[idx]

			structVal := extractStructTag(param.SuiMoveNormalizedType)
			if structVal != nil {
				argumentToResolve = append(argumentToResolve, &ResolveArgument{Type: Object, Id: arg.(string), NormalizedType: param.SuiMoveNormalizedType})
				continue
			}

			stringType, ok := param.SuiMoveNormalizedType.(types.SuiMoveNormalizedType_String)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type %v", param.SuiMoveNormalizedType)
			}

			switch stringType {
			case "Bool", "U8", "U16", "U32", "U64", "U128", "U256":
			case "Address":
				arg, err = sui_types.NewAddressFromHex(utils.NormalizeSuiObjectId(arg.(string)))
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unimplemented sui move normalized type, index: %d, type: %v", idx, param)
			}
			argumentToResolve = append(argumentToResolve, &ResolveArgument{Type: Pure, PureValue: arg})
		}

		var ids []string
		for _, argResolve := range argumentToResolve {
			if argResolve.Type == Object {
				ids = append(ids, argResolve.Id)
			}
		}

		if len(ids) > 0 {
			objects, err := ptb.client.MultiGetObjects(types.MultiGetObjectsParams{IDs: ids, Options: &types.SuiObjectDataOptions{ShowOwner: true}})
			if err != nil {
				return nil, err
			}

			for _, resolve := range argumentToResolve {
				if resolve.Type == Pure {
					continue
				}
				objectInfo := gm.FilterOne(objects, func(v *types.SuiObjectResponse) bool {
					return v.Data.ObjectId == utils.NormalizeSuiObjectId(resolve.Id)
				})
				objectId, err := sui_types.NewObjectIdFromHex(objectInfo.Data.ObjectId)
				if err != nil {
					return nil, err
				}

				var initialSharedVersion *uint64
				sharedObject, ok := objectInfo.Data.Owner.ObjectOwner.(types.ObjectOwner_Shared)
				if ok {
					initialSharedVersion = &sharedObject.Shared.InitialSharedVersion
				}

				if initialSharedVersion != nil {
					mutable := extractMutableReference(resolve.NormalizedType) != nil
					resolve.ObjectValue.SharedObject = &struct {
						Id                   move_types.AccountAddress
						InitialSharedVersion uint64
						Mutable              bool
					}{
						Id:                   *objectId,
						InitialSharedVersion: *initialSharedVersion,
						Mutable:              mutable,
					}
				} else {
					version, err := strconv.ParseUint(objectInfo.Data.Version, 10, 64)
					if err != nil {
						return nil, err
					}

					digest, err := sui_types.NewDigest(objectInfo.Data.Digest)
					if err != nil {
						return nil, err
					}
					resolve.ObjectValue.ImmOrOwnedObject = &sui_types.ObjectRef{
						ObjectId: *objectId,
						Version:  version,
						Digest:   *digest,
					}
				}
			}
		}

		for _, arg := range argumentToResolve {
			switch arg.Type {
			case Pure:
				argument, err := ptb.builder.Pure(arg.PureValue)
				if err != nil {
					return nil, err
				}
				arguments = append(arguments, argument)
			case Object:
				argument, err := ptb.builder.Obj(arg.ObjectValue)
				if err != nil {
					return nil, err
				}
				arguments = append(arguments, argument)
			}
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

func (ptb *ProgrammableTransactionBlock) Build(options BuildOptions) (*sui_types.TransactionData, []byte, error) {
	if err := ptb.prepare(options); err != nil {
		return nil, nil, err
	}

	return ptb.build(buildOption{})
}

// Internal Functions
func (ptb *ProgrammableTransactionBlock) build(options buildOption) (*sui_types.TransactionData, []byte, error) {
	programmableTransaction := ptb.builder.Finish()
	if options.OnlyTransactionKind {
		kind := sui_types.TransactionKind{
			ProgrammableTransaction: &programmableTransaction,
		}

		bs, err := bcs.Marshal(kind)
		return nil, bs, err
	}

	sender := ptb.sender
	if options.Overrides.Sender != "" {
		sender = options.Overrides.Sender
	}

	gasBudget := ptb.gasConfig.Budget
	if options.Overrides.GasConfig.Budget != 0 {
		gasBudget = options.Overrides.GasConfig.Budget
	}

	gasPayment := ptb.gasConfig.Payment
	if options.Overrides.GasConfig.Payment != nil {
		gasPayment = options.Overrides.GasConfig.Payment
	}

	gasPrice := ptb.gasConfig.Price
	if options.Overrides.GasConfig.Price != 0 {
		gasPrice = options.Overrides.GasConfig.Price
	}

	if sender == "" {
		return nil, nil, fmt.Errorf("missing transaction sender")
	}

	if gasBudget == 0 {
		return nil, nil, fmt.Errorf("missing gas budget")
	}

	if gasPayment == nil {
		return nil, nil, fmt.Errorf("missing gas payment")
	}

	if gasPrice == 0 {
		return nil, nil, fmt.Errorf("missing gas price")
	}

	s, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		return nil, nil, err
	}
	tx := sui_types.NewProgrammable(
		*s,
		gasPayment,
		programmableTransaction,
		gasBudget,
		gasPrice,
	)
	bs, err := bcs.Marshal(tx)
	return &tx, bs, err
}

// prepare transaction and transaction config
func (ptb *ProgrammableTransactionBlock) prepare(options BuildOptions) error {
	if !options.OnlyTransactionKind && ptb.sender == "" {
		return fmt.Errorf("missing transaction sender")
	}

	if err := ptb.prepareGasPrice(options); err != nil {
		return err
	}

	if !options.OnlyTransactionKind {
		if err := ptb.prepareGasPayment(options); err != nil {
			return err
		}

		if ptb.gasConfig.Budget == 0 {
			_, blockData, err := ptb.build(buildOption{Overrides: buildOptionOverride{GasConfig: GasConfig{Budget: MAX_TX_GAS, Payment: []*sui_types.ObjectRef{}}}})
			if err != nil {
				return err
			}
			dryRunResult, err := ptb.client.DryRunTransactionBlock(types.DryRunTransactionBlockParams{
				TransactionBlock: blockData,
			})
			if err != nil {
				return err
			}

			if dryRunResult.Effects.Status.Status != "success" {
				return fmt.Errorf("dry run failed, could not automatically determine a budget: %v", dryRunResult.Effects.Status.Error)
			}

			computationCost, err := strconv.ParseUint(dryRunResult.Effects.GasUsed.ComputationCost, 10, 64)
			if err != nil {
				return nil
			}

			storageCost, err := strconv.ParseUint(dryRunResult.Effects.GasUsed.StorageCost, 10, 64)
			if err != nil {
				return nil
			}

			storageRebate, err := strconv.ParseUint(dryRunResult.Effects.GasUsed.StorageRebate, 10, 64)
			if err != nil {
				return nil
			}

			_price := ptb.gasConfig.Price
			if _price == 0 {
				_price = 1
			}
			safeOverhead := GAS_SAFE_OVERHEAD * _price
			baseComputationCostWithOverhead := computationCost + safeOverhead
			gasBudget := baseComputationCostWithOverhead + storageCost - storageRebate

			if gasBudget < baseComputationCostWithOverhead {
				gasBudget = baseComputationCostWithOverhead
			}
			ptb.SetGasBudget(gasBudget)
		}
	}

	return nil
}

func (ptb *ProgrammableTransactionBlock) prepareGasPrice(options BuildOptions) error {
	if options.OnlyTransactionKind || ptb.gasConfig.Price != 0 {
		return nil
	}

	referenceGasPrice, err := ptb.client.GetReferenceGasPrice()
	if err != nil {
		return err
	}

	ptb.SetGasPrice(referenceGasPrice.Uint64() + 1)
	return nil
}

func (ptb *ProgrammableTransactionBlock) prepareGasPayment(options BuildOptions) error {
	if len(ptb.gasConfig.Payment) > MAX_GAS_OBJECTS {
		return fmt.Errorf("payment objects exceed maximum amount: %v", MAX_GAS_OBJECTS)
	}

	if options.OnlyTransactionKind || len(ptb.gasConfig.Payment) > 0 {
		return nil
	}

	gasOwner := ptb.sender
	if ptb.gasConfig.Owner != "" {
		gasOwner = ptb.gasConfig.Owner
	}

	coins, err := ptb.client.GetCoins(
		types.GetCoinsParams{
			Owner:    gasOwner,
			CoinType: &SuiGasCoinType,
		},
	)
	if err != nil {
		return err
	}

	var paymentCoins []*sui_types.ObjectRef
	for _, coin := range coins.Data {
		if findObjectFromBuilderArgs(ptb.builder.Inputs, coin.CoinObjectId) {
			continue
		}

		if len(paymentCoins) > MAX_GAS_OBJECTS {
			continue
		}

		digest, err := sui_types.NewDigest(coin.Digest)
		if err != nil {
			return err
		}

		uint64Version, err := strconv.ParseUint(coin.Version, 10, 64)
		if err != nil {
			return err
		}

		objectId, err := sui_types.NewObjectIdFromHex(coin.CoinObjectId)
		if err != nil {
			return err
		}

		paymentCoins = append(paymentCoins, &sui_types.ObjectRef{ObjectId: *objectId, Version: uint64Version, Digest: *digest})
	}

	if len(paymentCoins) == 0 {
		return fmt.Errorf("no valid gas coins found for the transaction")
	}

	ptb.SetGasPayment(paymentCoins)
	return nil
}

func findObjectFromBuilderArgs(args map[string]sui_types.CallArg, objectId string) bool {
	for _, arg := range args {
		if arg.Object != nil && arg.Object.ImmOrOwnedObject != nil {
			if utils.NormalizeSuiObjectId(objectId) == utils.NormalizeSuiObjectId(arg.Object.ImmOrOwnedObject.ObjectId.String()) {
				return true
			}
		}
	}

	return false
}

// checking
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
