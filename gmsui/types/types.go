package types

type Balance struct {
	CoinObjectCount int               `json:"coinObjectCount"`
	CoinType        string            `json:"coinType"`
	LockedBalance   map[string]string `json:"lockedBalance"`
	TotalBalance    string            `json:"totalBalance"`
}

type BalanceChange struct {
	Owner    ObjectOwnerWrapper `json:"owner"`
	CoinType string             `json:"coinType"`
	Amount   string             `json:"amount"`
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

type GasCostSummary struct {
	ComputationCost         string `json:"computationCost"`
	StorageCost             string `json:"storageCost"`
	StorageRebate           string `json:"storageRebate"`
	NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
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
	Digest []int `json:"digest"` // TODO: bytes?
}

type CheckpointId string

type Claim struct {
	IndexMod4 uint64 `json:"indexMod4"`
	Value     string `json:"value"`
}

type CoinStruct struct {
	Balance             string `json:"balance"`
	CoinObjectId        string `json:"coinObjectId"`
	CoinType            string `json:"coinType"`
	Digest              string `json:"digest"`
	PreviousTransaction string `json:"previousTransaction"`
	Version             string `json:"version"`
}

type CommitteeInfo struct {
	Epoch      string      `json:"epoch"`
	Validators [][2]string `json:"validators"`
}

type DelegatedStake struct {
	ValidatorAddress string               `json:"validatorAddress"`
	StakingPool      string               `json:"stakingPool"`
	Stakes           []StakeObjectWrapper `json:"stakes"`
}

type DevInspectResults struct {
	Effects TransactionEffects   `json:"effects"`
	Error   string               `json:"error,omitempty"`
	Events  []SuiEvent           `json:"events"`
	Results []SuiExecutionResult `json:"results,omitempty"`
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

type TransactionBlockEffectsModifiedAtVersions struct {
	ObjectId       string `json:"objectId"`
	SequenceNumber string `json:"sequenceNumber"`
}

type SuiObjectRef struct {
	ObjectId string `json:"objectId"`
	Version  uint64 `json:"version"`
	Digest   string `json:"digest"`
}

type OwnedObjectRef struct {
	Owner     ObjectOwnerWrapper `json:"owner"`
	Reference SuiObjectRef       `json:"reference"`
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

type SuiExecutionResult struct {
	MutableReferenceOutputs [][3]interface{} `json:"mutableReferenceOutputs,omitempty"` // TODO: [SuiArgument, bytes, string][]
	ReturnValues            [][2]interface{} `json:"returnValues,omitempty"`            // TODO: interface -> [bytes, string][]
}

type DisplayFieldsResponse struct {
	Data  map[string]*string          `json:"data"`
	Error *ObjectResponseErrorWrapper `json:"error"`
}

type DryRunTransactionBlockResponse struct {
	Effects        TransactionEffects       `json:"effects"`
	Events         []SuiEvent               `json:"events"`
	ObjectChanges  []SuiObjectChangeWrapper `json:"objectChanges"`
	BalanceChanges []BalanceChange          `json:"balanceChanges"`
	Input          TransactionBlockData     `json:"input"`
}

type TransactionBlockData struct {
	MessageVersion string                         `json:"messageVersion"`
	Transaction    SuiTransactionBlockKindWrapper `json:"transaction"`
	Sender         string                         `json:"sender"`
	GasData        SuiGasData                     `json:"gasData"`
}

type SuiGasData struct {
	Payment []SuiObjectRef `json:"payment"`
	Owner   string         `json:"owner"`
	Price   string         `json:"price"`
	Budget  string         `json:"budget"`
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

type ExecuteTransactionRequestType string

var (
	WaitForEffectsCert    ExecuteTransactionRequestType = "WaitForEffectsCert"
	WaitForLocalExecution ExecuteTransactionRequestType = "WaitForLocalExecution"
)

type SuiObjectData struct {
	ObjectId            string                 `json:"objectId"`
	Version             string                 `json:"version"`
	Digest              string                 `json:"digest"`
	Type                *string                `json:"type,omitempty"`
	Owner               *ObjectOwnerWrapper    `json:"owner,omitempty"`
	PreviousTransaction *string                `json:"previousTransaction,omitempty"`
	StorageRebate       *string                `json:"storageRebate,omitempty"`
	Display             *DisplayFieldsResponse `json:"display,omitempty"`
	Content             *SuiParsedDataWrapper  `json:"content,omitempty"`
	Bcs                 *RawDataWrapper        `json:"bcs,omitempty"`
}

type SuiObjectDataOptions struct {
	ShowBcs                 bool `json:"showBcs,omitempty"`
	ShowContent             bool `json:"showContent,omitempty"`
	ShowDisplay             bool `json:"showDisplay,omitempty"`
	ShowOwner               bool `json:"showOwner,omitempty"`
	ShowPreviousTransaction bool `json:"showPreviousTransaction,omitempty"`
	ShowStorageRebate       bool `json:"showStorageRebate,omitempty"`
	ShowType                bool `json:"showType,omitempty"`
}

type SuiObjectResponseQuery struct {
	Filter  *SuiObjectDataFilter  `json:"filter,omitempty"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

type PaginatedCoins struct {
	Data        []CoinStruct `json:"data"`
	HasNextPage bool         `json:"hasNextPage"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
}

type PaginatedDynamicFieldInfos struct {
	Data        []DynamicFieldInfo `json:"data"`
	NextCursor  *string            `json:"nextCursor,omitempty"`
	HasNextPage bool               `json:"hasNextPage"`
}

type PaginatedEvents struct {
	Data        []SuiEvent `json:"data"`
	NextCursor  *EventId   `json:"nextCursor,omitempty"`
	HasNextPage bool       `json:"hasNextPage"`
}

type PaginatedStrings struct {
	Data        []string `json:"data"`
	NextCursor  *string  `json:"nextCursor,omitempty"`
	HasNextPage bool     `json:"hasNextPage"`
}

type PaginatedObjectsResponse struct {
	Data        []SuiObjectResponse `json:"data"`
	NextCursor  *string             `json:"nextCursor,omitempty"`
	HasNextPage bool                `json:"hasNextPage"`
}

type SuiObjectResponse struct {
	Data  *SuiObjectData              `json:"data,omitempty"`
	Error *ObjectResponseErrorWrapper `json:"error,omitempty"`
}

type PaginatedTransactionResponse struct {
	Data        []SuiTransactionBlockResponse `json:"data"`
	NextCursor  *string                       `json:"nextCursor,omitempty"`
	HasNextPage bool                          `json:"hasNextPage"`
}

type SuiTransactionBlockResponse struct {
	Digest                  string                    `json:"digest"`
	Transaction             *SuiTransactionBlock      `json:"transaction,omitempty"`
	RawTransaction          string                    `json:"rawTransaction,omitempty"`
	Effects                 *TransactionEffects       `json:"effects,omitempty"`
	Events                  []*SuiEvent               `json:"events,omitempty"`
	ObjectChanges           []*SuiObjectChangeWrapper `json:"objectChanges,omitempty"`
	BalanceChanges          []*BalanceChange          `json:"balanceChanges,omitempty"`
	TimestampMs             *string                   `json:"timestampMs,omitempty"`
	Checkpoint              *string                   `json:"checkpoint,omitempty"`
	ConfirmedLocalExecution *bool                     `json:"confirmedLocalExecution,omitempty"`
	Errors                  []string                  `json:"errors,omitempty"`
}

type SuiTransactionBlock struct {
	Data         TransactionBlockData `json:"data"`
	TxSignatures []string             `json:"txSignatures"`
}

type ProtocolConfig struct {
	MinSupportedProtocolVersion string                         `json:"minSupportedProtocolVersion"`
	MaxSupportedProtocolVersion string                         `json:"maxSupportedProtocolVersion"`
	ProtocolVersion             string                         `json:"protocolVersion"`
	FeatureFlags                map[string]bool                `json:"featureFlags"`
	Attributes                  map[string]ProtocolConfigValue `json:"attributes"`
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

type SuiAuthenticatorStateExpire struct {
	MinEpoch string `json:"min_epoch"`
}

type SuiChangeEpoch struct {
	ComputationCharge     string `json:"computation_charge"`
	Epoch                 string `json:"epoch"`
	EpochStartTimestampMs string `json:"epoch_start_timestamp_ms"`
	StorageCharge         string `json:"storage_charge"`
	StorageRebate         string `json:"storage_rebate"`
}

type CoinMetadata struct {
	Decimals    uint8  `json:"decimals"`
	Description string `json:"description"`
	IconUrl     string `json:"iconUrl,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
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

type SuiMoveModuleId struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type SuiMoveNormalizedField struct {
	Name string                       `json:"name"`
	Type SuiMoveNormalizedTypeWrapper `json:"type"`
}

type SuiMoveNormalizedFunction struct {
	Visibility     SuiMoveVisibility               `json:"visibility"`
	IsEntry        bool                            `json:"isEntry"`
	TypeParameters []SuiMoveAbilitySet             `json:"typeParameters"`
	Parameters     []*SuiMoveNormalizedTypeWrapper `json:"parameters"`
	Return         []*SuiMoveNormalizedTypeWrapper `json:"return"`
}

type SuiMoveVisibility string

var (
	Private SuiMoveVisibility = "Private"
	Public  SuiMoveVisibility = "Public"
	Friend  SuiMoveVisibility = "Friend"
)

type SuiMoveNormalizedModule struct {
	FileFormatVersion int                                  `json:"fileFormatVersion"`
	Address           string                               `json:"address"`
	Name              string                               `json:"name"`
	Friends           []SuiMoveModuleId                    `json:"friends"`
	Structs           map[string]SuiMoveNormalizedStruct   `json:"structs"`
	ExposedFunctions  map[string]SuiMoveNormalizedFunction `json:"exposedFunctions"`
}

type SuiMoveNormalizedStruct struct {
	Abilities      SuiMoveAbilitySet            `json:"abilities"`
	TypeParameters []SuiMoveStructTypeParameter `json:"typeParameters"`
	Fields         []SuiMoveNormalizedField     `json:"fields"`
}

type SuiMoveStructTypeParameter struct {
	Constraints SuiMoveAbilitySet `json:"constraints"`
	IsPhantom   bool              `json:"isPhantom"`
}

type MoveCallSuiTransaction struct {
	Package       string               `json:"package"`
	Module        string               `json:"module"`
	Function      string               `json:"function"`
	TypeArguments []*string            `json:"type_arguments,omitempty"`
	Arguments     []SuiArgumentWrapper `json:"arguments"`
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
	AtRiskValidators                      [][2]string           `json:"atRiskValidators"`
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

type SuiTransactionBlockBuilderMode string

const (
	Commit     SuiTransactionBlockBuilderMode = "Commit"
	DevInspect SuiTransactionBlockBuilderMode = "DevInspect"
)

type CoinSupply struct {
	Value string `json:"value"`
}

type TransactionBlockBytes struct {
	Gas          []SuiObjectRef           `json:"gas"`
	InputObjects []InputObjectKindWrapper `json:"inputObjects"`
	TxBytes      string                   `json:"txBytes"`
}

type SuiTransactionBlockResponseOptions struct {
	ShowInput          bool `json:"showInput,omitempty"`
	ShowEffects        bool `json:"showEffects,omitempty"`
	ShowEvents         bool `json:"showEvents,omitempty"`
	ShowObjectChanges  bool `json:"showObjectChanges,omitempty"`
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
	ShowRawInput       bool `json:"showRawInput,omitempty"`
}

type SuiTransactionBlockResponseQuery struct {
	Filter  *TransactionFilter                  `json:"filter,omitempty"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
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

type ValidatorApy struct {
	Address string  `json:"address"`
	APY     float64 `json:"apy"`
}

type ValidatorsApy struct {
	APYs  []ValidatorApy `json:"apys"`
	Epoch string         `json:"epoch"`
}

type ZkLoginAuthenticator struct {
	Inputs        ZkLoginInputs    `json:"inputs"`
	MaxEpoch      string           `json:"maxEpoch"`
	UserSignature SignatureWrapper `json:"userSignature"`
}

type ZkLoginInputs struct {
	AddressSeed      string       `json:"addressSeed"`
	HeaderBase64     string       `json:"headerBase64"`
	IssBase64Details Claim        `json:"issBase64Details"`
	ProofPoints      ZkLoginProof `json:"proofPoints"`
}

type ZkLoginProof struct {
	A []string   `json:"a"`
	B [][]string `json:"b"`
	C []string   `json:"c"`
}

type PaginatedCheckpoints struct {
	Data        []Checkpoint `json:"data"`
	HasNextPage bool         `json:"hasNextPage"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
}

type ObjectValueKind string

const (
	ByImmutableReference ObjectValueKind = "ByImmutableReference"
	ByMutableReference   ObjectValueKind = "ByMutableReference"
	ByValue              ObjectValueKind = "ByValue"
)

type ResolvedNameServiceNames struct {
	Data        []string `json:"data"`
	NextCursor  *string  `json:"nextCursor,omitempty"`
	HasNextPage bool     `json:"hasNextPage"`
}

type EpochInfo struct {
	Epoch                  string                `json:"epoch"`
	Validators             []SuiValidatorSummary `json:"validators"`
	EpochTotalTransactions string                `json:"epochTotalTransactions"`
	FirstCheckpointId      string                `json:"firstCheckpointId"`
	EpochStartTimestamp    string                `json:"epochStartTimestamp"`
	EndOfEpochInfo         *EndOfEpochInfo       `json:"endOfEpochInfo"`
	ReferenceGasPrice      *uint64               `json:"referenceGasPrice"`
}

type EndOfEpochInfo struct {
	LastCheckpointId             string `json:"lastCheckpointId"`
	EpochEndTimestamp            string `json:"epochEndTimestamp"`
	ProtocolVersion              string `json:"protocolVersion"`
	ReferenceGasPrice            string `json:"referenceGasPrice"`
	TotalStake                   string `json:"totalStake"`
	StorageFundReinvestment      string `json:"storageFundReinvestment"`
	StorageCharge                string `json:"storageCharge"`
	StorageRebate                string `json:"storageRebate"`
	StorageFundBalance           string `json:"storageFundBalance"`
	StakeSubsidyAmount           string `json:"stakeSubsidyAmount"`
	TotalGasFees                 string `json:"totalGasFees"`
	TotalStakeRewardsDistributed string `json:"totalStakeRewardsDistributed"`
	LeftoverStorageFundInflow    string `json:"leftoverStorageFundInflow"`
}

type EpochPage struct {
	Data        []EpochInfo `json:"data"`
	NextCursor  string      `json:"nextCursor,omitempty"`
	HasNextPage bool        `json:"hasNextPage"`
}

type DynamicFieldPage struct {
	Data        []DynamicFieldInfo `json:"data"`
	NextCursor  *string            `json:"nextCursor,omitempty"`
	HasNextPage bool               `json:"hasNextPage"`
}

type CheckpointPage struct {
	Data        []Checkpoint `json:"data"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
	HasNextPage bool         `json:"hasNextPage"`
}

type SuiMoveNormalizedModules map[string]SuiMoveNormalizedModule

type ProgrammableTransaction struct {
	Transactions []SuiTransactionWrapper `json:"transactions"`
	Inputs       []SuiCallArgWrapper     `json:"inputs"`
}

type GetPastObjectRequest struct {
	ObjectID string `json:"objectId"`
	Version  string `json:"version"`
}

type LoadedChildObject struct {
	ObjectID       string `json:"objectId"`
	SequenceNumber string `json:"sequenceNumber"`
}

type LoadedChildObjectsResponse struct {
	LoadedChildObjects []LoadedChildObject `json:"loadedChildObjects"`
}

type MoveCallParams struct {
	Arguments       []interface{} `json:"arguments"`
	Function        string        `json:"function"`
	Module          string        `json:"module"`
	PackageObjectId string        `json:"packageObjectId"`
	TypeArguments   []string      `json:"typeArguments,omitempty"`
}

type TransferObjectParams struct {
	ObjectID  string `json:"objectId"`
	Recipient string `json:"recipient"`
}
