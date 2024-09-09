package types

type TransactionFilter struct {
	Checkpoint        *string                             `json:"Checkpoint,omitempty"`
	MoveFunction      *TransactionFilter_MoveFunction     `json:"MoveFunction,omitempty"`
	InputObject       *string                             `json:"InputObject,omitempty"`
	ChangedObject     *string                             `json:"ChangedObject,omitempty"`
	FromAddress       *string                             `json:"FromAddress,omitempty"`
	ToAddress         *string                             `json:"ToAddress,omitempty"`
	FromAndToAddress  *TransactionFilter_FromAndToAddress `json:"FromAndToAddress,omitempty"`
	FromOrToAddress   *TransactionFilter_FromOrToAddress  `json:"FromOrToAddress,omitempty"`
	TransactionKind   *string                             `json:"TransactionKind,omitempty"`
	TransactionKindIn []*string                           `json:"TransactionKindIn,omitempty"`
}

type SuiObjectDataFilter struct {
	*SuiObjectDataFilter_MatchAll
	*SuiObjectDataFilter_MatchAny
	*SuiObjectDataFilter_MatchNone
	*SuiObjectDataFilter_Package
	*SuiObjectDataFilter_MoveModule
	*SuiObjectDataFilter_StructType
	*SuiObjectDataFilter_AddressOwner
	*SuiObjectDataFilter_ObjectOwner
	*SuiObjectDataFilter_ObjectId
	*SuiObjectDataFilter_ObjectIds
	*SuiObjectDataFilter_Version
}

type SuiEventFilter struct {
	Sender          *string                         `json:"Sender,omitempty"`
	Transaction     *string                         `json:"Transaction,omitempty"`
	Package         *string                         `json:"Package,omitempty"`
	MoveModule      *SuiEventFilter_MoveModule      `json:"MoveModule,omitempty"`
	MoveEventType   *string                         `json:"MoveEventType,omitempty"`
	MoveEventModule *SuiEventFilter_MoveEventModule `json:"MoveEventModule,omitempty"`
	MoveEventField  *SuiEventFilter_MoveEventField  `json:"MoveEventField,omitempty"`
	TimeRange       *SuiEventFilter_TimeRange       `json:"TimeRange,omitempty"`
	All             *SuiEventFilters                `json:"All,omitempty"`
	Any             *SuiEventFilters                `json:"Any,omitempty"`
	And             *SuiEventFilters                `json:"And,omitempty"`
	Or              *SuiEventFilters                `json:"Or,omitempty"`
}
