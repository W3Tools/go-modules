package types

import (
	"encoding/json"
)

type PaginatedCoins struct {
	Data        []CoinStruct `json:"data"`
	HasNextPage bool         `json:"hasNextPage"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
}

type CoinStruct struct {
	Balance             string `json:"balance"`
	CoinObjectId        string `json:"coinObjectId"`
	CoinType            string `json:"coinType"`
	Digest              string `json:"digest"`
	PreviousTransaction string `json:"previousTransaction"`
	Version             string `json:"version"`
}

type Balance struct {
	CoinObjectCount int               `json:"coinObjectCount"`
	CoinType        string            `json:"coinType"`
	LockedBalance   map[string]string `json:"lockedBalance"`
	TotalBalance    string            `json:"totalBalance"`
}

type CoinMetadata struct {
	Decimals    int    `json:"decimals"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
}

type CoinSupply struct {
	Value string `json:"value"`
}

type SuiObjectResponse struct {
	Data  *SuiObjectData       `json:"data,omitempty"`
	Error *ObjectResponseError `json:"error,omitempty"`
}

type PaginatedObjectsResponse struct {
	Data        []SuiObjectResponse `json:"data"`
	NextCursor  *string             `json:"nextCursor,omitempty"`
	HasNextPage bool                `json:"hasNextPage"`
}

type ObjectRead struct {
	Details *json.RawMessage `json:"details"`
	Status  string           `json:"status"`
}

type DynamicFieldPage struct {
	Data        []DynamicFieldInfo `json:"data"`
	NextCursor  *string            `json:"nextCursor,omitempty"`
	HasNextPage bool               `json:"hasNextPage"`
}

type SuiTransactionBlockResponse struct {
	Digest                  string               `json:"digest"`
	Transaction             *SuiTransactionBlock `json:"transaction,omitempty"`
	RawTransaction          string               `json:"rawTransaction,omitempty"`
	Effects                 *TransactionEffects  `json:"effects,omitempty"`
	Events                  *[]SuiEvent          `json:"events,omitempty"`
	ObjectChanges           *[]SuiObjectChange   `json:"objectChanges,omitempty"`
	BalanceChanges          *[]BalanceChange     `json:"balanceChanges,omitempty"`
	TimestampMs             *string              `json:"timestampMs,omitempty"`
	Checkpoint              *string              `json:"checkpoint,omitempty"`
	ConfirmedLocalExecution *bool                `json:"confirmedLocalExecution,omitempty"`
	Errors                  []string             `json:"errors,omitempty"`
}

type PaginatedTransactionResponse struct {
	Data        []SuiTransactionBlockResponse `json:"data"`
	NextCursor  *string                       `json:"nextCursor,omitempty"`
	HasNextPage bool                          `json:"hasNextPage"`
}

type PaginatedEvents struct {
	Data        []SuiEvent `json:"data"`
	NextCursor  *EventId   `json:"nextCursor,omitempty"`
	HasNextPage bool       `json:"hasNextPage"`
}

type ProtocolConfig struct {
	MinSupportedProtocolVersion string                 `json:"minSupportedProtocolVersion"`
	MaxSupportedProtocolVersion string                 `json:"maxSupportedProtocolVersion"`
	ProtocolVersion             string                 `json:"protocolVersion"`
	FeatureFlags                map[string]bool        `json:"featureFlags"`
	Attributes                  map[string]interface{} `json:"attributes"`
}

type Checkpoint struct {
	Epoch                      string                 `json:"epoch"`
	SequenceNumber             string                 `json:"sequenceNumber"`
	Digest                     string                 `json:"digest"`
	NetworkTotalTransactions   string                 `json:"networkTotalTransactions"`
	PreviousDigest             string                 `json:"previousDigest,omitempty"`
	EpochRollingGasCostSummary GasCostSummary         `json:"epochRollingGasCostSummary"`
	TimestampMs                string                 `json:"timestampMs"`
	EndOfEpochData             *EndOfEpochData        `json:"endOfEpochData,omitempty"`
	Transactions               []string               `json:"transactions"`
	CheckpointCommitments      []CheckpointCommitment `json:"checkpointCommitments"`
	ValidatorSignature         string                 `json:"validatorSignature"`
}

type CheckpointPage struct {
	Data        []Checkpoint `json:"data"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
	HasNextPage bool         `json:"hasNextPage"`
}

type SuiSystemStateSummary struct {
	Epoch                                 string                `json:"epoch"`
	ProtocolVersion                       string                `json:"protocolVersion"`
	SystemStateVersion                    string                `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  string                `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       string                `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     string                `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                string                `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            string                `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                string                `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       string                `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 string                `json:"epochStartTimestampMs"`
	EpochDurationMs                       string                `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                string                `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     string                `json:"maxValidatorCount"`
	MinValidatorJoiningStake              string                `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            string                `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        string                `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          string                `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   string                `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       string                `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount string                `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              string                `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              int                   `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            string                `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsId             string                `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           string                `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []string              `json:"pendingRemovals"`
	StakingPoolMappingsId                 string                `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               string                `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       string                `json:"inactivePoolsId"`
	InactivePoolsSize                     string                `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 string                `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               string                `json:"validatorCandidatesSize"`
	AtRiskValidators                      [][2]interface{}      `json:"atRiskValidators"`
	ValidatorReportRecords                [][2]interface{}      `json:"validatorReportRecords"`
}

type SuiValidatorSummary struct {
	SuiAddress                   string  `json:"suiAddress"`
	ProtocolPubkeyBytes          string  `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes           string  `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes            string  `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes       string  `json:"proofOfPossessionBytes"`
	Name                         string  `json:"name"`
	Description                  string  `json:"description"`
	ImageUrl                     string  `json:"imageUrl"`
	ProjectUrl                   string  `json:"projectUrl"`
	NetAddress                   string  `json:"netAddress"`
	P2pAddress                   string  `json:"p2pAddress"`
	PrimaryAddress               string  `json:"primaryAddress"`
	WorkerAddress                string  `json:"workerAddress"`
	NextEpochProtocolPubkeyBytes *string `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   *string `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  *string `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   *string `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          *string `json:"nextEpochNetAddress"`
	NextEpochP2pAddress          *string `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      *string `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       *string `json:"nextEpochWorkerAddress"`
	VotingPower                  string  `json:"votingPower"`
	OperationCapId               string  `json:"operationCapId"`
	GasPrice                     string  `json:"gasPrice"`
	CommissionRate               string  `json:"commissionRate"`
	NextEpochStake               string  `json:"nextEpochStake"`
	NextEpochGasPrice            string  `json:"nextEpochGasPrice"`
	NextEpochCommissionRate      string  `json:"nextEpochCommissionRate"`
	StakingPoolId                string  `json:"stakingPoolId"`
	StakingPoolActivationEpoch   *string `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch *string `json:"stakingPoolDeactivationEpoch"`
	StakingPoolSuiBalance        string  `json:"stakingPoolSuiBalance"`
	RewardsPool                  string  `json:"rewardsPool"`
	PoolTokenBalance             string  `json:"poolTokenBalance"`
	PendingStake                 string  `json:"pendingStake"`
	PendingTotalSuiWithdraw      string  `json:"pendingTotalSuiWithdraw"`
	PendingPoolTokenWithdraw     string  `json:"pendingPoolTokenWithdraw"`
	ExchangeRatesId              string  `json:"exchangeRatesId"`
	ExchangeRatesSize            string  `json:"exchangeRatesSize"`
}

type CommitteeInfo struct {
	Epoch      string      `json:"epoch"`
	Validators [][2]string `json:"validators"`
}

type ValidatorsApy struct {
	APYs  []ValidatorApy `json:"apys"`
	Epoch string         `json:"epoch"`
}

type ValidatorApy struct {
	Address string  `json:"address"`
	APY     float64 `json:"apy"`
}

type DelegatedStake struct {
	ValidatorAddress string        `json:"validatorAddress"`
	StakingPool      string        `json:"stakingPool"`
	Stakes           []StakeObject `json:"stakes"`
}

type StakeObject struct {
	StakedSuiId       string  `json:"stakedSuiId"`               // through Pending/Active/Unstaked
	StakeRequestEpoch string  `json:"stakeRequestEpoch"`         // through Pending/Active/Unstaked
	StakeActiveEpoch  string  `json:"stakeActiveEpoch"`          // through Pending/Active/Unstaked
	Principal         string  `json:"principal"`                 // through Pending/Active/Unstaked
	Status            string  `json:"status"`                    // through Pending/Active/Unstaked
	EstimatedReward   *string `json:"estimatedReward,omitempty"` // through Active
}

type ResolvedNameServiceNames struct {
	Data        []string `json:"data"`
	NextCursor  *string  `json:"nextCursor,omitempty"`
	HasNextPage bool     `json:"hasNextPage"`
}

type SuiMoveNormalizedModules map[string]SuiMoveNormalizedModule

type SuiMoveNormalizedModule struct {
	FileFormatVersion int                                  `json:"fileFormatVersion"`
	Address           string                               `json:"address"`
	Name              string                               `json:"name"`
	Friends           []SuiMoveModuleId                    `json:"friends"`
	Structs           map[string]SuiMoveNormalizedStruct   `json:"structs"`
	ExposedFunctions  map[string]SuiMoveNormalizedFunction `json:"exposedFunctions"`
}

type SuiMoveModuleId struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type SuiMoveNormalizedStruct struct {
	Abilities      SuiMoveAbilitySet            `json:"abilities"`
	TypeParameters []SuiMoveStructTypeParameter `json:"typeParameters"`
	Fields         []SuiMoveNormalizedField     `json:"fields"`
}

type SuiMoveAbilitySet struct {
	Abilities []SuiMoveAbility `json:"abilities"`
}

type SuiMoveAbility string

var (
	Copy  SuiMoveAbility = "Copy"
	Drop  SuiMoveAbility = "Drop"
	Store SuiMoveAbility = "Store"
	Key   SuiMoveAbility = "Key"
)

type SuiMoveStructTypeParameter struct {
	Constraints SuiMoveAbilitySet `json:"constraints"`
	IsPhantom   bool              `json:"isPhantom"`
}

type SuiMoveNormalizedField struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"` // native type: SuiMoveNormalizedType
}

type SuiMoveNormalizedFunction struct {
	Visibility     SuiMoveVisibility   `json:"visibility"`
	IsEntry        bool                `json:"isEntry"`
	TypeParameters []SuiMoveAbilitySet `json:"typeParameters"`
	Parameters     []interface{}       `json:"parameters"` // native type: []SuiMoveNormalizedType
	Return         []interface{}       `json:"return"`     // native type: []SuiMoveNormalizedType
}

type DryRunTransactionBlockResponse struct {
	Effects        TransactionEffects   `json:"effects"`
	Events         []SuiEvent           `json:"events"`
	ObjectChanges  []SuiObjectChange    `json:"objectChanges"`
	BalanceChanges []BalanceChange      `json:"balanceChanges"`
	Input          TransactionBlockData `json:"input"`
}

type DevInspectResults struct {
	Effects TransactionEffects   `json:"effects"`
	Error   string               `json:"error,omitempty"`
	Events  []SuiEvent           `json:"events"`
	Results []SuiExecutionResult `json:"results,omitempty"`
}

type ObjectResponseError struct {
	NotExists            ObjectResponseErrorCode `json:"notExists"`
	DynamicFieldNotFound ObjectResponseErrorCode `json:"dynamicFieldNotFound"`
	Deleted              ObjectResponseErrorCode `json:"deleted"`
	Unknown              ObjectResponseErrorCode `json:"unknown"`
	DisplayError         ObjectResponseErrorCode `json:"displayError"`
}

type ObjectResponseErrorCode struct {
	Code           string `jons:"code"`
	ObjectId       string `json:"object_id,omitempty"`
	ParentObjectId string `json:"parent_object_id,omitempty"`
	Digest         string `json:"digest"`
	Version        string `json:"version"`
	Error          string `json:"error"`
}

type SuiObjectData struct {
	ObjectId            string                 `json:"objectId"`
	Version             string                 `json:"version"`
	Digest              string                 `json:"digest"`
	Type                *string                `json:"type,omitempty"`
	Owner               *json.RawMessage       `json:"owner,omitempty"` // native type: ObjectOwner
	PreviousTransaction *string                `json:"previousTransaction,omitempty"`
	StorageRebate       *string                `json:"storageRebate,omitempty"`
	Display             *DisplayFieldsResponse `json:"display,omitempty"`
	Content             *SuiParsedData         `json:"content,omitempty"`
	Bcs                 *RawData               `json:"bcs,omitempty"`
}

type DisplayFieldsResponse struct {
	Data  *map[string]string   `json:"data,omitempty"`
	Error *ObjectResponseError `json:"error,omitempty"`
}

type SuiParsedData struct {
	DataType          string                  `json:"dataType"`                    //
	Type              *string                 `json:"type,omitempty"`              // through moveObject
	HasPublicTransfer *bool                   `json:"hasPublicTransfer,omitempty"` // through moveObject
	Fields            *interface{}            `json:"fields,omitempty"`            // through moveObject, native type: MoveStruct
	Disassembled      *map[string]interface{} `json:"disassembled,omitempty"`      // through package
}

type RawData struct {
	DataType          string                 `json:"dataType"`                    //
	Id                string                 `json:"id,omitempty"`                // through package
	Type              string                 `json:"type,omitempty"`              // through moveObject
	HasPublicTransfer *bool                  `json:"hasPublicTransfer,omitempty"` // through moveObject
	Version           int64                  `json:"version"`                     //
	BcsBytes          *string                `json:"bcsBytes,omitempty"`          // through moveObject
	ModuleMap         map[string]string      `json:"moduleMap,omitempty"`         // through package
	TypeOriginTable   []TypeOrigin           `json:"typeOriginTable,omitempty"`   // through package
	LinkageTable      map[string]UpgradeInfo `json:"linkageTable,omitempty"`      // through package
}

type TypeOrigin struct {
	ModuleName   string `json:"module_name"`
	DataTypeName string `json:"datatype_name"`
	Package      string `json:"package"`
}

type UpgradeInfo struct {
	UpgradedId      string `json:"upgraded_id"`
	UpgradedVersion int64  `json:"upgraded_version"`
}

type DynamicFieldInfo struct {
	Name       DynamicFieldName `json:"name"`
	BcsName    string           `json:"bcsName"`
	Type       DynamicFieldType `json:"type"`
	ObjectType string           `json:"objectType"`
	ObjectId   string           `json:"objectId"`
	Version    int64            `json:"version"`
	Digest     string           `json:"digest"`
}

type DynamicFieldName struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type DynamicFieldType string

var (
	DynamicField  DynamicFieldType = "DynamicField"
	DynamicObject DynamicFieldType = "DynamicObject"
)

type SuiTransactionBlock struct {
	Data         TransactionBlockData `json:"data"`
	TxSignatures []string             `json:"txSignatures"`
}

type TransactionBlockData struct {
	MessageVersion string                  `json:"messageVersion"`
	Transaction    SuiTransactionBlockKind `json:"transaction"`
	Sender         string                  `json:"sender"`
	GasData        SuiGasData              `json:"gasData"`
}

type SuiTransactionBlockKind struct {
	Kind                  string             `json:"kind"`                               // through ChangeEpoch/Genesis/ConsensusCommitPrologue/ProgrammableTransaction/AuthenticatorStateUpdate/EndOfEpochTransaction
	ComputationCharge     *string            `json:"computation_charge,omitempty"`       // through ChangeEpoch
	Epoch                 *string            `json:"epoch,omitempty"`                    // through ChangeEpoch/ConsensusCommitPrologue/AuthenticatorStateUpdate
	EpochStartTimestampMs *string            `json:"epoch_start_timestamp_ms,omitempty"` // through ChangeEpoch
	StorageCharge         *string            `json:"storage_charge,omitempty"`           // through ChangeEpoch
	StorageRebate         *string            `json:"storage_rebate,omitempty"`           // through ChangeEpoch
	Objects               *[]string          `json:"objects,omitempty"`                  // through Genesis
	CommitTimestampMs     *string            `json:"commit_timestamp_ms,omitempty"`      // through ConsensusCommitPrologue
	Round                 *string            `json:"round,omitempty"`                    // through ConsensusCommitPrologue/AuthenticatorStateUpdate
	Inputs                *[]SuiCallArg      `json:"inputs,omitempty"`                   // through ProgrammableTransaction
	Transactions          *[]json.RawMessage `json:"transactions,omitempty"`             // through ProgrammableTransaction
	NewActiveJwks         *[]SuiActiveJwk    `json:"new_active_jwks,omitempty"`          // through AuthenticatorStateUpdate
}

type SuiCallArg struct {
	Type                 string       `json:"type"`                           // through immOrOwnedObject/sharedObject/receiving/pure
	ObjectType           *string      `json:"objectType,omitempty"`           // through immOrOwnedObject/sharedObject/receiving
	ObjectId             *string      `json:"objectId,omitempty"`             // through immOrOwnedObject/sharedObject/receiving
	InitialSharedVersion *string      `json:"initialSharedVersion,omitempty"` // through sharedObject
	Version              *string      `json:"version,omitempty"`              // through immOrOwnedObject/receiving
	Mutable              *bool        `json:"mutable,omitempty"`              // through sharedObject
	Digest               *string      `json:"digest,omitempty"`               // through immOrOwnedObject/receiving
	ValueType            *string      `json:"valueType,omitempty"`            // through pure
	Value                *interface{} `json:"value,omitempty"`                // through pure
}

type SuiActiveJwk struct {
	Epoch string   `json:"epoch"`
	Jwk   SuiJWK   `json:"jwk"`
	JwkID SuiJwkID `json:"jwk_id"`
}

type SuiJWK struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kty string `json:"kty"`
	N   string `json:"n"`
}

type SuiJwkID struct {
	Iss string `json:"iss"`
	Kid string `json:"kid"`
}

type SuiGasData struct {
	Payment []SuiObjectRef `json:"payment"`
	Owner   string         `json:"owner"`
	Price   string         `json:"price"`
	Budget  string         `json:"budget"`
}

type SuiObjectRef struct {
	ObjectId string `json:"objectId"`
	Version  int    `json:"version"`
	Digest   string `json:"digest"`
}

type TransactionEffects struct {
	MessageVersion       string                                      `json:"messageVersion"`
	Status               ExecutionStatus                             `json:"status"`
	ExecutedEpoch        string                                      `json:"executedEpoch"`
	GasUsed              GasCostSummary                              `json:"gasUsed"`
	ModifiedAtVersions   []TransactionBlockEffectsModifiedAtVersions `json:"modifiedAtVersions,omitempty"`
	SharedObjects        []SuiObjectRef                              `json:"sharedObjects,omitempty"`
	TransactionDigest    string                                      `json:"transactionDigest"`
	Created              []OwnedObjectRef                            `json:"created,omitempty"`
	Mutated              []OwnedObjectRef                            `json:"mutated,omitempty"`
	Deleted              []SuiObjectRef                              `json:"deleted,omitempty"`
	GasObject            OwnedObjectRef                              `json:"gasObject"`
	EventsDigest         *string                                     `json:"eventsDigest,omitempty"`
	Dependencies         []string                                    `json:"dependencies,omitempty"`
	Unwrapped            []OwnedObjectRef                            `json:"unwrapped,omitempty"`
	UnwrappedThenDeleted []SuiObjectRef                              `json:"unwrappedThenDeleted,omitempty"`
	Wrapped              []SuiObjectRef                              `json:"wrapped,omitempty"`
}

type ExecutionStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type GasCostSummary struct {
	ComputationCost         string `json:"computationCost"`
	StorageCost             string `json:"storageCost"`
	StorageRebate           string `json:"storageRebate"`
	NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
}

type TransactionBlockEffectsModifiedAtVersions struct {
	ObjectId       string `json:"objectId"`
	SequenceNumber string `json:"sequenceNumber"`
}

type OwnedObjectRef struct {
	Owner     interface{}  `json:"owner"` // native type: ObjectOwner
	Reference SuiObjectRef `json:"reference"`
}

type SuiEvent struct {
	Id                EventId     `json:"id"`
	PackageId         string      `json:"packageId"`
	TransactionModule string      `json:"transactionModule"`
	Sender            string      `json:"sender"`
	Type              string      `json:"type"`
	ParsedJson        interface{} `json:"parsedJson"`
	Bcs               string      `json:"bcs"`
	TimestampMs       string      `json:"timestampMs,omitempty"`
}

type EventId struct {
	TxDigest string `json:"txDigest"`
	EventSeq string `json:"eventSeq"`
}

type SuiObjectChange struct {
	Type            string           `json:"type"`                      //
	Sender          *string          `json:"sender,omitempty"`          // through transferred/mutated/deleted/wrapped/created
	Recipient       *json.RawMessage `json:"recipient,omitempty"`       // through transferred/ native type: ObjectOwner
	Owner           *json.RawMessage `json:"owner,omitempty"`           // through mutated/created native type: ObjectOwner
	ObjectType      *string          `json:"objectType,omitempty"`      // through transferred/mutated/deleted/wrapped/created
	ObjectId        *string          `json:"objectId,omitempty"`        // through transferred/mutated/deleted/wrapped/created
	Version         string           `json:"version"`                   //
	PreviousVersion *string          `json:"previousVersion,omitempty"` // through mutated/
	Digest          *string          `json:"digest,omitempty"`          // through published/transferred/mutated/created
	Modules         *[]string        `json:"modules,omitempty"`         // through published/
	PackageId       *string          `json:"packageId,omitempty"`       // through published/
}

type BalanceChange struct {
	Owner    interface{} `json:"owner"` // native type: ObjectOwner
	CoinType string      `json:"coinType"`
	Amount   string      `json:"amount"`
}

type EndOfEpochData struct {
	EpochCommitments         []CheckpointCommitment `json:"epochCommitments"`
	NextEpochCommittee       [][2]string            `json:"nextEpochCommittee"`
	NextEpochProtocolVersion string                 `json:"nextEpochProtocolVersion"`
}

type CheckpointCommitment struct {
	ECMHLiveObjectSetDigest ECMHLiveObjectSetDigest `json:"ecmhLiveObjectSetDigest"`
}

type ECMHLiveObjectSetDigest struct {
	Digest []int `json:"digest"`
}

type SuiExecutionResult struct {
	MutableReferenceOutputs [][3]interface{} `json:"mutableReferenceOutputs,omitempty"`
	ReturnValues            [][2]interface{} `json:"returnValues,omitempty"`
}

type SuiMoveVisibility string

var (
	Private SuiMoveVisibility = "Private"
	Public  SuiMoveVisibility = "Public"
	Friend  SuiMoveVisibility = "Friend"
)
