package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/W3Tools/go-modules/gmsui/types"
	"github.com/W3Tools/go-modules/gmsui/utils"
)

type SuiClientOptions struct {
	URL string
}

type SuiClient struct {
	ctx        context.Context
	rpc        string
	requestId  int
	httpClient *http.Client
}

type SuiTransportRequestOptions struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}

func NewSuiClient(ctx context.Context, rpc string) (*SuiClient, error) {
	_, err := url.ParseRequestURI(rpc)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    5,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}

	return &SuiClient{ctx: ctx, rpc: rpc, requestId: 1, httpClient: httpClient}, nil
}

func (client *SuiClient) Context() context.Context {
	return client.ctx
}

func (client *SuiClient) RPC() string {
	return client.rpc
}

// Invoke any RPC method
func (client *SuiClient) Call(method string, params []any, response any) error {
	return client.request(SuiTransportRequestOptions{Method: method, Params: params}, &response)
}

func (client *SuiClient) GetRpcApiVersion() (string, error) {
	var response struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	err := client.request(SuiTransportRequestOptions{Method: "rpc.discover", Params: []any{}}, &response)
	if err != nil {
		return "", err
	}

	return response.Info.Version, nil
}

// Get all Coin<`coin_type`> objects owned by an address.
func (client *SuiClient) GetCoins(input types.GetCoinsParams) (response *types.PaginatedCoins, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	if input.CoinType != nil || *input.CoinType != "" {
		input.CoinType = gm.NewStringPtr(utils.NormalizeSuiCoinType(*input.CoinType))
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getCoins",
			Params: []any{utils.NormalizeSuiAddress(input.Owner), input.CoinType, input.Cursor, input.Limit},
		},
		&response,
	)
}

// Get all Coin objects owned by an address.
func (client *SuiClient) GetAllCoins(input types.GetAllCoinsParams) (response *types.PaginatedCoins, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getAllCoins",
			Params: []any{utils.NormalizeSuiAddress(input.Owner), input.Cursor, input.Limit},
		},
		&response,
	)
}

// Get the total coin balance for one coin type, owned by the address owner.
func (client *SuiClient) GetBalance(input types.GetBalanceParams) (response *types.Balance, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	if input.CoinType != nil || *input.CoinType != "" {
		input.CoinType = gm.NewStringPtr(utils.NormalizeSuiCoinType(*input.CoinType))
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getBalance",
			Params: []any{utils.NormalizeSuiAddress(input.Owner), input.CoinType},
		},
		&response,
	)
}

// Get the total coin balance for all coin types, owned by the address owner.
func (client *SuiClient) GetAllBalances(input types.GetAllBalancesParams) (response []*types.Balance, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getAllBalances",
			Params: []any{utils.NormalizeSuiAddress(input.Owner)},
		},
		&response,
	)
}

// Fetch CoinMetadata for a given coin type
func (client *SuiClient) GetCoinMetadata(input types.GetCoinMetadataParams) (response *types.CoinMetadata, err error) {
	if input.CoinType != "" {
		input.CoinType = utils.NormalizeSuiCoinType(input.CoinType)
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getCoinMetadata",
			Params: []any{input.CoinType},
		},
		&response,
	)
}

// Fetch total supply for a coin
func (client *SuiClient) GetTotalSupply(input types.GetTotalSupplyParams) (response *types.CoinSupply, err error) {
	if input.CoinType != "" {
		input.CoinType = utils.NormalizeSuiCoinType(input.CoinType)
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getTotalSupply",
			Params: []any{input.CoinType},
		},
		&response,
	)
}

// Get details about an object
func (client *SuiClient) GetObject(input types.GetObjectParams) (response *types.SuiObjectResponse, err error) {
	if input.ID == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(input.ID)) {
		return nil, fmt.Errorf("invalid sui object id")
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getObject",
			Params: []any{utils.NormalizeSuiObjectId(input.ID), input.Options},
		},
		&response,
	)
}

// Batch get details about a list of objects. If any of the object ids are duplicates the call will fail
func (client *SuiClient) MultiGetObjects(input types.MultiGetObjectsParams) (response []*types.SuiObjectResponse, err error) {
	_, err = gm.Map(input.IDs, func(v string) (any, error) {
		if v == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(v)) {
			return nil, fmt.Errorf("invalid sui object id %s", v)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, id := range input.IDs {
		ids = gm.UniqueAppend(ids, id)
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_multiGetObjects",
			Params: []any{ids, input.Options},
		},
		&response,
	)
}

// Get all objects owned by an address
func (client *SuiClient) GetOwnedObjects(input types.GetOwnedObjectsParams) (response *types.PaginatedObjectsResponse, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getOwnedObjects",
			Params: []any{
				utils.NormalizeSuiAddress(input.Owner),
				types.SuiObjectResponseQuery{
					Filter:  input.SuiObjectResponseQuery.Filter,
					Options: input.SuiObjectResponseQuery.Options,
				},
				input.Cursor,
				input.Limit,
			},
		},
		&response,
	)
}

// Return the object information for a specified version
func (client *SuiClient) TryGetPastObject(input types.TryGetPastObjectParams) (response *types.ObjectRead, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_tryGetPastObject",
			Params: []any{utils.NormalizeSuiObjectId(input.ID), input.Version, input.Options},
		},
		&response,
	)
}

// Return the list of dynamic field objects owned by an object
func (client *SuiClient) GetDynamicFields(input types.GetDynamicFieldsParams) (response *types.DynamicFieldPage, err error) {
	if input.ParentId == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(input.ParentId)) {
		return nil, fmt.Errorf("invalid sui object id")
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getDynamicFields",
			Params: []any{utils.NormalizeSuiObjectId(input.ParentId), input.Cursor, input.Limit},
		},
		&response,
	)
}

// Return the dynamic field object information for a specified object
func (client *SuiClient) GetDynamicFieldObject(input types.GetDynamicFieldObjectParams) (response *types.SuiObjectResponse, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getDynamicFieldObject",
			Params: []any{input.ParentId, input.Name},
		},
		&response,
	)
}

func (client *SuiClient) GetTransactionBlock(input types.GetTransactionBlockParams) (response *types.SuiTransactionBlockResponse, err error) {
	if !utils.IsValidTransactionDigest(input.Digest) {
		return nil, fmt.Errorf("invalid transaction digest")
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getTransactionBlock",
			Params: []any{input.Digest, input.Options},
		},
		&response,
	)
}

func (client *SuiClient) MultiGetTransactionBlocks(input types.MultiGetTransactionBlocksParams) (response []*types.SuiTransactionBlockResponse, err error) {
	_, err = gm.Map(input.Digests, func(v string) (any, error) {
		if v == "" || !utils.IsValidTransactionDigest(v) {
			return nil, fmt.Errorf("invalid transaction digest %s", v)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, id := range input.Digests {
		ids = gm.UniqueAppend(ids, id)
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_multiGetTransactionBlocks",
			Params: []any{ids, input.Options},
		},
		&response,
	)
}

// Get transaction blocks for a given query criteria
func (client *SuiClient) QueryTransactionBlocks(input types.QueryTransactionBlocksParams) (response *types.PaginatedTransactionResponse, err error) {
	var order bool = true
	if input.Order != nil {
		order = (input.Order == &types.Descending)
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_queryTransactionBlocks",
			Params: []any{
				input.SuiTransactionBlockResponseQuery,
				input.Cursor,
				input.Limit,
				&order,
			},
		},
		&response,
	)
}

// Get total number of transactions
func (client *SuiClient) GetTotalTransactionBlocks() (response *big.Int, err error) {
	var response_ string
	err = client.request(
		SuiTransportRequestOptions{
			Method: "sui_getTotalTransactionBlocks",
			Params: []any{},
		},
		&response_,
	)
	if err != nil {
		return nil, err
	}

	response, ok := new(big.Int).SetString(response_, 10)
	if !ok {
		return nil, fmt.Errorf("got invalid string number %s", response_)
	}

	return response, nil
}

func (client *SuiClient) SubscribeTransaction(input types.SubscribeTransactionParams) (response any, err error) {
	return nil, fmt.Errorf("unimplemented")
}

// Get events for a given query criteria
func (client *SuiClient) QueryEvents(input types.QueryEventsParams) (response *types.PaginatedEvents, err error) {
	var order bool = true
	if input.Order != nil {
		order = (input.Order == &types.Descending)
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_queryEvents",
			Params: []any{input.Query, input.Cursor, input.Limit, order},
		},
		&response,
	)
}

// Subscribe to get notifications whenever an event matching the filter occurs
func (client *SuiClient) SubscribeEvent(input types.SubscribeEventParams) (response any, err error) {
	return nil, fmt.Errorf("unimplemented")
}

func (client *SuiClient) GetProtocolConfig(input types.GetProtocolConfigParams) (response *types.ProtocolConfig, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getProtocolConfig",
			Params: []any{input.Version},
		},
		&response,
	)
}

// Get the sequence number of the latest checkpoint that has been executed
func (client *SuiClient) GetLatestCheckpointSequenceNumber() (response string, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getLatestCheckpointSequenceNumber",
			Params: []any{},
		},
		&response,
	)
}

// Returns information about a given checkpoint
func (client *SuiClient) GetCheckpoint(input types.GetCheckpointParams) (response *types.Checkpoint, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getCheckpoint",
			Params: []any{input.ID},
		},
		&response,
	)
}

// Returns historical checkpoints paginated
func (client *SuiClient) GetCheckpoints(input types.GetCheckpointsParams) (response *types.CheckpointPage, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getCheckpoints",
			Params: []any{input.Cursor, input.Limit, input.DescendingOrder},
		},
		&response,
	)
}

// Getting the reference gas price for the network
func (client *SuiClient) GetReferenceGasPrice() (response *big.Int, err error) {
	var response_ string
	err = client.request(
		SuiTransportRequestOptions{
			Method: "suix_getReferenceGasPrice",
			Params: []any{},
		},
		&response_,
	)
	if err != nil {
		return nil, err
	}

	response, ok := new(big.Int).SetString(response_, 10)
	if !ok {
		return nil, fmt.Errorf("got invalid string number %s", response_)
	}

	return response, nil
}

// Return the latest system state content.
func (client *SuiClient) GetLatestSuiSystemState() (response *types.SuiSystemStateSummary, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getLatestSuiSystemState",
			Params: []any{},
		},
		&response,
	)
}

// Return the committee information for the asked epoch
func (client *SuiClient) GetCommitteeInfo(input types.GetCommitteeInfoParams) (response *types.CommitteeInfo, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getCommitteeInfo",
			Params: []any{input.Epoch},
		},
		&response,
	)
}

// Return the Validators APYs
func (client *SuiClient) GetValidatorsApy() (response *types.ValidatorsApy, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getValidatorsApy",
			Params: []any{},
		},
		&response,
	)
}

func (client *SuiClient) GetChainIdentifier() (response string, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getChainIdentifier",
			Params: []any{},
		},
		&response,
	)
}

// Return the delegated stakes for an address
func (client *SuiClient) GetStakes(input types.GetStakesParams) (response []*types.DelegatedStake, err error) {
	if input.Owner == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getStakes",
			Params: []any{utils.NormalizeSuiObjectId(input.Owner)},
		},
		&response,
	)
}

// Return the delegated stakes queried by id.
func (client *SuiClient) GetStakesByIds(input types.GetStakesByIdsParams) (response []*types.DelegatedStake, err error) {
	_, err = gm.Map(input.StakedSuiIds, func(v string) (any, error) {
		if v == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(v)) {
			return nil, fmt.Errorf("invalid sui object id %s", v)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, id := range input.StakedSuiIds {
		ids = gm.UniqueAppend(ids, id)
	}

	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_getStakesByIds",
			Params: []any{ids},
		},
		&response,
	)
}

func (client *SuiClient) ResolveNameServiceNames(input types.ResolveNameServiceNamesParams) (response *types.ResolvedNameServiceNames, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_resolveNameServiceNames",
			Params: []any{input.Address, input.Cursor, input.Limit},
		},
		&response,
	)
}

func (client *SuiClient) ResolveNameServiceAddress(input types.ResolveNameServiceAddressParams) (response string, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "suix_resolveNameServiceAddress",
			Params: []any{input.Name},
		},
		&response,
	)
}

// Get Move function argument types like read, write and full access
// The native response type is []types.SuiMoveFunctionArgType, It contains two different types: string and types.SuiMoveFunctionArgType_Object
func (client *SuiClient) GetMoveFunctionArgTypes(input types.GetMoveFunctionArgTypesParams) (response []interface{}, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getMoveFunctionArgTypes",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module, input.Function},
		},
		&response,
	)
}

// Get a map from module name to structured representations of Move modules
func (client *SuiClient) GetNormalizedMoveModulesByPackage(input types.GetNormalizedMoveModulesByPackageParams) (response *types.SuiMoveNormalizedModules, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveModulesByPackage",
			Params: []any{utils.NormalizeSuiObjectId(input.Package)},
		},
		&response,
	)
}

// Get a structured representation of Move module
func (client *SuiClient) GetNormalizedMoveModule(input types.GetNormalizedMoveModuleParams) (response *types.SuiMoveNormalizedModule, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveModule",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module},
		},
		&response,
	)
}

// Get a structured representation of Move function
func (client *SuiClient) GetNormalizedMoveFunction(input types.GetNormalizedMoveFunctionParams) (response *types.SuiMoveNormalizedFunction, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveFunction",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module, input.Function},
		},
		&response,
	)
}

// Get a structured representation of Move struct
func (client *SuiClient) GetNormalizedMoveStruct(input types.GetNormalizedMoveStructParams) (response *types.SuiMoveNormalizedStruct, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveStruct",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module, input.Struct},
		},
		&response,
	)
}

// Dry run a transaction block and return the result.
func (client *SuiClient) DryRunTransactionBlock(input types.DryRunTransactionBlockParams) (response *types.DryRunTransactionBlockResponse, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_dryRunTransactionBlock",
			Params: []any{input.TransactionBlock},
		},
		&response,
	)
}

// Runs the transaction block in dev-inspect mode.
// Which allows for nearly any transaction (or Move call) with any arguments.
// Detailed results are provided, including both the transaction effects and any return values.
func (client *SuiClient) DevInspectTransactionBlock(input types.DevInspectTransactionBlockParams) (response *types.DevInspectResults, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_devInspectTransactionBlock",
			Params: []any{input.Sender, input.TransactionBlock, input.GasPrice, input.Epoch},
		},
		&response,
	)
}

func (client *SuiClient) ExecuteTransactionBlock(input types.ExecuteTransactionBlockParams) (response *types.SuiTransactionBlockResponse, err error) {
	return response, client.request(
		SuiTransportRequestOptions{
			Method: "sui_executeTransactionBlock",
			Params: []any{b64.ToBase64(input.TransactionBlock), input.Signature, input.Options, input.RequestType},
		},
		&response,
	)
}

func (client *SuiClient) SignAndExecuteTransactionBlock(input types.SignAndExecuteTransactionBlockParams) (response *types.SuiTransactionBlockResponse, err error) {
	signatureData, err := input.Signer.SignTransactionBlock(input.TransactionBlock)
	if err != nil {
		return nil, err
	}

	return client.ExecuteTransactionBlock(types.ExecuteTransactionBlockParams{
		TransactionBlock: input.TransactionBlock,
		Signature:        []string{signatureData.Signature},
		Options:          input.Options,
		RequestType:      input.RequestType,
	})
}
