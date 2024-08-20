package types

import "github.com/W3Tools/go-modules/gmsui/cryptography"

type GetCoinsParams struct {
	Owner    string  `json:"owner"`
	CoinType *string `json:"coinType,omitempty"`
	Cursor   *string `json:"cursor,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
}

type GetAllCoinsParams struct {
	Owner  string  `json:"owner"`
	Cursor *string `json:"cursor,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

type GetBalanceParams struct {
	Owner    string  `json:"owner"`
	CoinType *string `json:"coinType,omitempty"`
}

type GetAllBalancesParams struct {
	Owner string `json:"owner"`
}

type GetCoinMetadataParams struct {
	CoinType string `json:"coinType"`
}

type GetTotalSupplyParams struct {
	CoinType string `json:"coinType"`
}

type GetObjectParams struct {
	ID      string                `json:"id"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

type MultiGetObjectsParams struct {
	IDs     []string              `json:"ids"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

type GetOwnedObjectsParams struct {
	Owner                  string                 `json:"owner"`
	Cursor                 *string                `json:"cursor,omitempty"`
	Limit                  *int                   `json:"limit,omitempty"`
	SuiObjectResponseQuery SuiObjectResponseQuery `json:",inline"`
}

type TryGetPastObjectParams struct {
	ID      string                `json:"id"`
	Version int                   `json:"version"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

type GetDynamicFieldsParams struct {
	ParentId string  `json:"parentId"`
	Cursor   *string `json:"cursor,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
}

type GetDynamicFieldObjectParams struct {
	ParentId string           `json:"parentId"`
	Name     DynamicFieldName `json:"name"`
}

type GetTransactionBlockParams struct {
	Digest  string                              `json:"digest"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

type MultiGetTransactionBlocksParams struct {
	Digests []string                            `json:"digests"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

type QueryTransactionBlocksParams struct {
	Cursor                           *string                             `json:"cursor,omitempty"`
	Limit                            *int                                `json:"limit,omitempty"`
	Order                            *QueryTransactionBlocksParams_Order `json:"order,omitempty"`
	SuiTransactionBlockResponseQuery SuiTransactionBlockResponseQuery    `json:",inline"`
}

type QueryEventsParams struct {
	Query  SuiEventFilter                      `json:"query"`
	Cursor *EventId                            `json:"cursor,omitempty"`
	Limit  *int                                `json:"limit,omitempty"`
	Order  *QueryTransactionBlocksParams_Order `json:"order,omitempty"`
}

type GetProtocolConfigParams struct {
	Version *string `json:"version,omitempty"`
}

type GetCheckpointParams struct {
	ID CheckpointId `json:"id"`
}

type GetCheckpointsParams struct {
	Cursor          *string `json:"cursor,omitempty"`
	Limit           *int    `json:"limit,omitempty"`
	DescendingOrder bool    `json:"descendingOrder"`
}

type GetCommitteeInfoParams struct {
	Epoch *string `json:"epoch,omitempty"`
}

type SubscribeEventParams struct {
	Filter SuiEventFilter `json:"filter"`
}

type SubscribeTransactionParams struct {
	Filter TransactionFilter `json:"filter"`
}

type GetStakesParams struct {
	Owner string `json:"owner"`
}

type GetStakesByIdsParams struct {
	StakedSuiIds []string `json:"stakedSuiIds"`
}

type ResolveNameServiceNamesParams struct {
	Address string  `json:"address"`
	Cursor  *string `json:"cursor,omitempty"`
	Limit   *int    `json:"limit,omitempty"`
}

type GetMoveFunctionArgTypesParams struct {
	Package  string `json:"package"`
	Module   string `json:"module"`
	Function string `json:"function"`
}

type GetNormalizedMoveModulesByPackageParams struct {
	Package string `json:"package"`
}

type GetNormalizedMoveModuleParams struct {
	Package string `json:"package"`
	Module  string `json:"module"`
}

type GetNormalizedMoveFunctionParams struct {
	Package  string `json:"package"`
	Module   string `json:"module"`
	Function string `json:"function"`
}

type GetNormalizedMoveStructParams struct {
	Package string `json:"package"`
	Module  string `json:"module"`
	Struct  string `json:"struct"`
}

type ResolveNameServiceAddressParams struct {
	Name string `json:"name"`
}

type DryRunTransactionBlockParams struct {
	TransactionBlock interface{} `json:"transactionBlock"`
}

type DevInspectTransactionBlockParams struct {
	Sender           string      `json:"sender"`
	TransactionBlock interface{} `json:"transactionBlock"`
	GasPrice         *uint64     `json:"gasPrice,omitempty"`
	Epoch            *string     `json:"epoch,omitempty"`
}

type ExecuteTransactionBlockParams struct {
	TransactionBlock []byte                              `json:"transactionBlock"`
	Signature        []string                            `json:"signature"`
	Options          *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
	RequestType      *ExecuteTransactionRequestType      `json:"requestType,omitempty"`
}

type SignAndExecuteTransactionBlockParams struct {
	TransactionBlock []byte                              `json:"transactionBlock"`
	Signer           cryptography.Signer                 `json:"signer"`
	Options          *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
	RequestType      *ExecuteTransactionRequestType      `json:"requestType,omitempty"`
}

type SuiObjectDataFilter_MatchAll struct {
	MatchAll []SuiObjectDataFilter `json:"MatchAll"`
}

type SuiObjectDataFilter_MatchAny struct {
	MatchAny []SuiObjectDataFilter `json:"MatchAny"`
}

type SuiObjectDataFilter_MatchNone struct {
	MatchNone []SuiObjectDataFilter `json:"MatchNone"`
}

type SuiObjectDataFilter_Package struct {
	Package string `json:"Package"`
}

type SuiObjectDataFilter_MoveModule struct {
	MoveModule SuiObjectDataFilter_MoveModule_Struct `json:"MoveModule"`
}

type SuiObjectDataFilter_MoveModule_Struct struct {
	Module  string `json:"module"`
	Package string `json:"package"`
}

type SuiObjectDataFilter_StructType struct {
	StructType string `json:"StructType"`
}

type SuiObjectDataFilter_AddressOwner struct {
	AddressOwner string `json:"AddressOwner"`
}

type SuiObjectDataFilter_ObjectOwner struct {
	ObjectOwner string `json:"ObjectOwner"`
}

type SuiObjectDataFilter_ObjectId struct {
	ObjectId string `json:"ObjectId"`
}

type SuiObjectDataFilter_ObjectIds struct {
	ObjectIds []string `json:"ObjectIds"`
}

type SuiObjectDataFilter_Version struct {
	Version string `json:"Version"`
}

type QueryTransactionBlocksParams_Order string

var (
	Ascending  QueryTransactionBlocksParams_Order = "ascending"
	Descending QueryTransactionBlocksParams_Order = "descending"
)

type TransactionFilter_MoveFunction struct {
	Function *string `json:"function,omitempty"`
	Module   *string `json:"module,omitempty"`
	Package  string  `json:"package"`
}

type TransactionFilter_FromAndToAddress struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type TransactionFilter_FromOrToAddress struct {
	Addr string `json:"addr"`
}

type SuiEventFilter_MoveModule struct {
	Module  string `json:"module"`
	Package string `json:"package"`
}

type SuiEventFilter_MoveEventModule struct {
	Module  string `json:"module"`
	Package string `json:"package"`
}

type SuiEventFilter_MoveEventField struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type SuiEventFilter_TimeRange struct {
	EndTime   string `json:"endTime"`
	StartTime string `json:"startTime"`
}

type SuiEventFilters []SuiEventFilter
