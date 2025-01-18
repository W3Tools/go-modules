package transactions

import (
	"context"
	"fmt"
	"strings"

	"github.com/W3Tools/go-modules/gmsui/client"
	"github.com/W3Tools/go-modules/gmsui/types"
	"github.com/W3Tools/go-modules/gmsui/utils"
	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/fardream/go-bcs/bcs"
)

type Transaction struct {
	client  *client.SuiClient
	builder *sui_types.ProgrammableTransactionBuilder
	ctx     context.Context

	Sender    *sui_types.SuiAddress `json:"sender"`
	GasConfig *GasData              `json:"gasConfig"`
}

type GasData struct {
	Budget  uint64                 `json:"budget"`
	Price   uint64                 `json:"price"`
	Owner   string                 `json:"owner"`
	Payment []*sui_types.ObjectRef `json:"payment"`
}

func NewTransaction(client *client.SuiClient) *Transaction {
	return &Transaction{
		client:  client,
		builder: sui_types.NewProgrammableTransactionBuilder(),
		ctx:     client.Context(),

		GasConfig: new(GasData),
	}
}

func (txb *Transaction) TransactionBuilder() *sui_types.ProgrammableTransactionBuilder {
	return txb.builder
}

func (txb *Transaction) Client() *client.SuiClient {
	return txb.client
}

func (txb *Transaction) Context() context.Context {
	return txb.ctx
}

type TransactionInputGasCoin struct {
	GasCoin bool `json:"gasCoin"`
}

func (txb *Transaction) Gas() *TransactionInputGasCoin {
	return &TransactionInputGasCoin{GasCoin: true}
}

// Split coins into multiple parts
func (txb *Transaction) SplitCoins(coin interface{}, amounts []interface{}) (returnArguments []*sui_types.Argument, err error) {
	if len(amounts) == 0 {
		return nil, nil
	}

	unresolvedParameter, err := txb.resolveSplitCoinsCoin(coin)
	if err != nil {
		return nil, fmt.Errorf("can not resolve coin in command %d: %v", len(txb.builder.Commands), err)
	}

	amountArguments, err := txb.resolveSplitCoinsAmounts(amounts)
	if err != nil {
		return nil, fmt.Errorf("can not resolve amounts in command %d: %v", len(txb.builder.Commands), err)
	}

	unresolvedParameter.merge(amountArguments)
	arguments, err := unresolvedParameter.resolveAndPArseToArguments(txb.client, txb)
	if err != nil {
		return nil, fmt.Errorf("can not resolve and parse to arguments in command %d, err: %v", len(txb.builder.Commands), err)
	}

	txb.builder.Command(
		sui_types.Command{
			SplitCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{
				Argument:  arguments[0],
				Arguments: arguments[1:],
			},
		},
	)

	return txb.createTransactionResult(len(amounts)), nil
}

// Transfer objects to address
func (txb *Transaction) TransferObjects(objects []interface{}, address interface{}) error {
	if len(objects) == 0 {
		return nil
	}

	unresolvedParameter, err := txb.resolveTransferObjectsObjects(objects)
	if err != nil {
		return fmt.Errorf("can not resolve objects in command %d: %v", len(txb.builder.Commands), err)
	}

	unresolvedAddressArgument, err := txb.resolveTransferObjectsAddress(address)
	if err != nil {
		return fmt.Errorf("can not resolve address in command %d: %v", len(txb.builder.Commands), err)
	}

	unresolvedParameter.merge(unresolvedAddressArgument)
	arguments, err := unresolvedParameter.resolveAndPArseToArguments(txb.client, txb)
	if err != nil {
		return fmt.Errorf("can not resolve and parse to arguments in command %d, err: %v", len(txb.builder.Commands), err)
	}

	txb.builder.Command(
		sui_types.Command{
			TransferObjects: &struct {
				Arguments []sui_types.Argument
				Argument  sui_types.Argument
			}{
				Arguments: arguments[:len(arguments)-1],
				Argument:  arguments[len(arguments)-1],
			},
		},
	)

	return nil
}

// Merge Coins into one
func (txb *Transaction) MergeCoins(destination interface{}, sources []interface{}) (err error) {
	if len(sources) == 0 {
		return nil
	}

	unresolvedParameter, err := txb.resolveMergeCoinsDestination(destination)
	if err != nil {
		return fmt.Errorf("failed to resolve destination in command %d, err: %v", len(txb.builder.Commands), err)
	}

	unresolvedSourceParameter, err := txb.resolveMergeCoinsSources(sources)
	if err != nil {
		return fmt.Errorf("failed to resolve sources in command %d, err: %v", len(txb.builder.Commands), err)
	}

	unresolvedParameter.merge(unresolvedSourceParameter)
	arguments, err := unresolvedParameter.resolveAndPArseToArguments(txb.client, txb)
	if err != nil {
		return fmt.Errorf("can not resolve and parse to arguments in command %d, err: %v", len(txb.builder.Commands), err)
	}

	txb.builder.Command(
		sui_types.Command{
			MergeCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{
				Argument:  arguments[0],
				Arguments: arguments[1:],
			},
		},
	)

	return nil
}

func (txb *Transaction) MoveCall(target string, arguments []interface{}, typeArguments []string) (returnArguments []*sui_types.Argument, err error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	var pkg, mod, fn = utils.NormalizeSuiObjectId(entry[0]), entry[1], entry[2]

	inputArguments, inputTypeArguments, returnsCount, err := txb.resolveMoveFunction(pkg, mod, fn, arguments, typeArguments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function arguments, err: %v", err)
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
				Arguments:     inputArguments,
				TypeArguments: inputTypeArguments,
			},
		},
	)

	return txb.createTransactionResult(returnsCount), nil
}

func (txb *Transaction) Build(sender string) (*sui_types.TransactionData, []byte, error) {
	txb.SetSenderIfNotSet(sender)

	if err := setGasPrice(txb); err != nil {
		return nil, nil, fmt.Errorf("can not set gas price when building transaction: %v", err)
	}
	if err := setGasBudget(txb); err != nil {
		return nil, nil, fmt.Errorf("can not set gas budget when building transaction: %v", err)
	}
	if err := setGasPayment(txb); err != nil {
		return nil, nil, fmt.Errorf("can not set gas payment when building transaction: %v", err)
	}

	tx := sui_types.NewProgrammable(
		*txb.Sender,
		txb.GasConfig.Payment,
		txb.builder.Finish(),
		txb.GasConfig.Budget,
		txb.GasConfig.Price,
	)
	bs, err := bcs.Marshal(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("can not marshal transaction: %v", err)
	}
	return &tx, bs, err
}

func (txb *Transaction) DryRunTransactionBlock() (*types.DryRunTransactionBlockResponse, error) {
	if txb.Sender == nil {
		return nil, fmt.Errorf("missing transaction sender")
	}

	if err := setGasPrice(txb); err != nil {
		return nil, fmt.Errorf("failed to set gas price, err: %v", err)
	}

	tx := sui_types.NewProgrammable(
		*txb.Sender,
		nil,
		txb.builder.Finish(),
		utils.MAX_GAS,
		txb.GasConfig.Price,
	)
	bs, err := bcs.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("can not marshal transaction, err: %v", err)
	}

	return txb.client.DryRunTransactionBlock(types.DryRunTransactionBlockParams{TransactionBlock: bs})
}

func (txb *Transaction) DevInspectTransactionBlock() (*types.DevInspectResults, error) {
	if txb.Sender == nil {
		return nil, fmt.Errorf("missing transaction sender")
	}

	bs, err := bcs.Marshal(txb.builder.Finish())
	if err != nil {
		return nil, fmt.Errorf("can not marshal transaction: %v", err)
	}

	txBytes := append([]byte{0}, bs...)
	return txb.client.DevInspectTransactionBlock(types.DevInspectTransactionBlockParams{Sender: txb.Sender.String(), TransactionBlock: txBytes})
}

func (txb *Transaction) SetSender(sender string) {
	address, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		panic(fmt.Errorf("failed to create address from hex [%s], err: %v", sender, err))
	}

	txb.Sender = address
}

func (txb *Transaction) SetSenderIfNotSet(sender string) {
	if txb.Sender == nil {
		txb.SetSender(sender)
	}
}

func (txb *Transaction) SetGasPrice(price uint64) {
	txb.GasConfig.Price = price
}

func (txb *Transaction) SetGasBudget(budget uint64) {
	txb.GasConfig.Budget = budget
}

func (txb *Transaction) SetGasBudgetIfNotSet(budget uint64) {
	if txb.GasConfig.Budget == 0 {
		txb.GasConfig.Budget = budget
	}
}

func (txb *Transaction) SetGasOwner(owner string) {
	txb.GasConfig.Owner = owner
}

func (txb *Transaction) SetGasPayment(payments []*sui_types.ObjectRef) {
	txb.GasConfig.Payment = payments
}
