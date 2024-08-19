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
			return nil, fmt.Errorf("invalid coin object %s, cause: %v", object.Error.ObjectId, object.Error.Code)
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

	objectIds := []string{}
	for idx, arg := range functionArgumentTypes {
		jsb, err := json.Marshal(arg)
		if err != nil {
			return nil, err
		}
		if string(jsb) == `{"Object":"ByImmutableReference"}` || string(jsb) == `{"Object":"ByMutableReference"}` {
			if args[idx] == nil {
				continue
			}
			input, ok := args[idx].(string)
			if !ok {
				return nil, fmt.Errorf("invalid object input, index %d, value: %v", idx, args[idx])
			}
			if !utils.IsHex(utils.NormalizeSuiObjectId(input)) {
				return nil, fmt.Errorf("input data not object, index %d, value: %v", idx, input)
			}

			objectIds = append(objectIds, utils.NormalizeSuiObjectId(input))
		}
	}

	inputObjects, _, err := GetObjectsAndUnmarshal[any](ptb.client, objectIds)
	if err != nil {
		return nil, err
	}

	// dd, _ := json.Marshal(inputObjects)
	// fmt.Printf("ids: %v\n", string(dd))
	// fmt.Printf("functionArgumentTypes: %v\n", functionArgumentTypes)

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

			objectInfo := gm.FilterOne(inputObjects, func(v *types.SuiObjectResponse) bool {
				return v.Data.ObjectId == utils.NormalizeSuiObjectId(inputArgument.(string))
			})

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
