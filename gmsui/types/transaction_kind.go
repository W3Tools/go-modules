package types

import (
	"encoding/json"
	"errors"
)

// -----------SuiTransactionBlockKind-----------
type SuiTransactionBlockKind interface {
	isSuiTransactionBlockKind()
}

type SuiTransactionBlockKindChangeEpoch struct {
	ComputationCharge     string `json:"computation_charge"`
	Epoch                 string `json:"epoch"`
	EpochStartTimestampMs string `json:"epoch_start_timestamp_ms"`
	Kind                  string `json:"kind"`
	StorageCharge         string `json:"storage_charge"`
	StorageRebate         string `json:"storage_rebate"`
}

type SuiTransactionBlockKindGenesis struct {
	Kind    string   `json:"kind"`
	Objects []string `json:"objects"`
}

type SuiTransactionBlockKindConsensusCommitPrologue struct {
	CommitTimestampMs string `json:"commit_timestamp_ms"`
	Epoch             string `json:"epoch"`
	Kind              string `json:"kind"`
	Round             string `json:"round"`
}

type SuiTransactionBlockKindConsensusCommitPrologueV3 struct {
	Kind                                  string                                `json:"kind"`
	Epoch                                 string                                `json:"epoch"`
	Round                                 string                                `json:"round"`
	SubDagIndex                           interface{}                           `json:"sub_dag_index"`
	CommitTimestampMs                     string                                `json:"commit_timestamp_ms"`
	ConsensusCommitDigest                 string                                `json:"consensus_commit_digest"`
	ConsensusDeterminedVersionAssignments ConsensusDeterminedVersionAssignments `json:"consensus_determined_version_assignments"`
}

type ConsensusDeterminedVersionAssignments struct {
	CancelledTransactions []interface{} `json:"CancelledTransactions"`
}

type SuiTransactionBlockKindProgrammableTransaction struct {
	Kind         string                  `json:"kind"`
	Inputs       []SuiCallArgWrapper     `json:"inputs"`
	Transactions []SuiTransactionWrapper `json:"transactions"`
}

type SuiTransactionBlockKindAuthenticatorStateUpdate struct {
	Epoch         string         `json:"epoch"`
	Kind          string         `json:"kind"`
	NewActiveJwks []SuiActiveJwk `json:"new_active_jwks"`
	Round         string         `json:"round"`
}

type SuiTransactionBlockKindEndOfEpochTransaction struct {
	Kind         string                                `json:"kind"`
	Transactions []SuiEndOfEpochTransactionKindWrapper `json:"Transactions"`
}

func (SuiTransactionBlockKindChangeEpoch) isSuiTransactionBlockKind()               {}
func (SuiTransactionBlockKindGenesis) isSuiTransactionBlockKind()                   {}
func (SuiTransactionBlockKindConsensusCommitPrologue) isSuiTransactionBlockKind()   {}
func (SuiTransactionBlockKindConsensusCommitPrologueV3) isSuiTransactionBlockKind() {}
func (SuiTransactionBlockKindProgrammableTransaction) isSuiTransactionBlockKind()   {}
func (SuiTransactionBlockKindAuthenticatorStateUpdate) isSuiTransactionBlockKind()  {}
func (SuiTransactionBlockKindEndOfEpochTransaction) isSuiTransactionBlockKind()     {}

type SuiTransactionBlockKindWrapper struct {
	SuiTransactionBlockKind
}

func (w *SuiTransactionBlockKindWrapper) UnmarshalJSON(data []byte) error {
	type Kind struct {
		Kind string `json:"kind"`
	}

	var kind Kind
	if err := json.Unmarshal(data, &kind); err != nil {
		return err
	}

	switch kind.Kind {
	case "ChangeEpoch":
		var k SuiTransactionBlockKindChangeEpoch
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiTransactionBlockKind = k
	case "Genesis":
		var k SuiTransactionBlockKindGenesis
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiTransactionBlockKind = k
	case "ConsensusCommitPrologue":
		var k SuiTransactionBlockKindConsensusCommitPrologue
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiTransactionBlockKind = k
	case "ConsensusCommitPrologueV3":
		var k SuiTransactionBlockKindConsensusCommitPrologueV3
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiTransactionBlockKind = k
	case "ProgrammableTransaction":
		var k SuiTransactionBlockKindProgrammableTransaction
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiTransactionBlockKind = k
	case "AuthenticatorStateUpdate":
		var k SuiTransactionBlockKindAuthenticatorStateUpdate
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiTransactionBlockKind = k
	case "EndOfEpochTransaction":
		var k SuiTransactionBlockKindEndOfEpochTransaction
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiTransactionBlockKind = k
	default:
		return errors.New("unknown SuiTransactionBlockKind type")
	}

	return nil
}

func (w SuiTransactionBlockKindWrapper) MarshalJSON() ([]byte, error) {
	switch t := w.SuiTransactionBlockKind.(type) {
	case SuiTransactionBlockKindChangeEpoch:
		return json.Marshal(SuiTransactionBlockKindChangeEpoch{
			Kind:                  t.Kind,
			ComputationCharge:     t.ComputationCharge,
			Epoch:                 t.Epoch,
			EpochStartTimestampMs: t.EpochStartTimestampMs,
			StorageCharge:         t.StorageCharge,
			StorageRebate:         t.StorageRebate,
		})
	case SuiTransactionBlockKindGenesis:
		return json.Marshal(SuiTransactionBlockKindGenesis{
			Kind:    t.Kind,
			Objects: t.Objects,
		})
	case SuiTransactionBlockKindConsensusCommitPrologue:
		return json.Marshal(SuiTransactionBlockKindConsensusCommitPrologue{
			Kind:              t.Kind,
			Epoch:             t.Epoch,
			CommitTimestampMs: t.CommitTimestampMs,
			Round:             t.Round,
		})
	case SuiTransactionBlockKindConsensusCommitPrologueV3:
		return json.Marshal(SuiTransactionBlockKindConsensusCommitPrologueV3{
			Kind:                                  t.Kind,
			Epoch:                                 t.Epoch,
			Round:                                 t.Round,
			SubDagIndex:                           t.SubDagIndex,
			CommitTimestampMs:                     t.CommitTimestampMs,
			ConsensusCommitDigest:                 t.ConsensusCommitDigest,
			ConsensusDeterminedVersionAssignments: t.ConsensusDeterminedVersionAssignments,
		})
	case SuiTransactionBlockKindProgrammableTransaction:
		return json.Marshal(SuiTransactionBlockKindProgrammableTransaction{
			Kind:         t.Kind,
			Inputs:       t.Inputs,
			Transactions: t.Transactions,
		})
	case SuiTransactionBlockKindAuthenticatorStateUpdate:
		return json.Marshal(SuiTransactionBlockKindAuthenticatorStateUpdate{
			Kind:          t.Kind,
			Epoch:         t.Epoch,
			NewActiveJwks: t.NewActiveJwks,
			Round:         t.Round,
		})
	case SuiTransactionBlockKindEndOfEpochTransaction:
		return json.Marshal(SuiTransactionBlockKindEndOfEpochTransaction{
			Kind:         t.Kind,
			Transactions: t.Transactions,
		})
	default:
		return nil, errors.New("unknown SuiTransactionBlockKind type")
	}
}

// -----------Sui Transaction-----------
type SuiTransaction interface {
	isSuiTransaction()
}

type SuiTransactionMoveCall struct {
	MoveCall MoveCallSuiTransaction `json:"MoveCall"`
}

type SuiTransactionTransferObjects struct {
	TransferObjects [2]SuiTransactionArgumentWrapper `json:"TransferObjects"`
}

type SuiTransactionSplitCoins struct {
	SplitCoins [2]SuiTransactionArgumentWrapper `json:"SplitCoins"`
}

type SuiTransactionMergeCoins struct {
	MergeCoins [2]SuiTransactionArgumentWrapper `json:"MergeCoins"`
}

type SuiTransactionPublish struct {
	Publish []string `json:"Publish"`
}

type SuiTransactionUpgrade struct {
	Upgrade [3]SuiTransactionArgumentWrapper `json:"Upgrade"`
}

type SuiTransactionMakeMoveVec struct {
	MakeMoveVec [2]*SuiTransactionArgumentWrapper `json:"MakeMoveVec"`
}

func (SuiTransactionMoveCall) isSuiTransaction()        {}
func (SuiTransactionTransferObjects) isSuiTransaction() {}
func (SuiTransactionSplitCoins) isSuiTransaction()      {}
func (SuiTransactionMergeCoins) isSuiTransaction()      {}
func (SuiTransactionPublish) isSuiTransaction()         {}
func (SuiTransactionUpgrade) isSuiTransaction()         {}
func (SuiTransactionMakeMoveVec) isSuiTransaction()     {}

type SuiTransactionWrapper struct {
	SuiTransaction
}

func (w *SuiTransactionWrapper) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _, ok := obj["MoveCall"]; ok {
		var t SuiTransactionMoveCall
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		w.SuiTransaction = t
		return nil
	}

	if _, ok := obj["TransferObjects"]; ok {
		var t SuiTransactionTransferObjects
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		w.SuiTransaction = t
		return nil
	}

	if _, ok := obj["SplitCoins"]; ok {
		var t SuiTransactionSplitCoins
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		w.SuiTransaction = t
		return nil
	}

	if _, ok := obj["MergeCoins"]; ok {
		var t SuiTransactionMergeCoins
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		w.SuiTransaction = t
		return nil
	}

	if _, ok := obj["Publish"]; ok {
		var t SuiTransactionPublish
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		w.SuiTransaction = t
		return nil
	}

	if _, ok := obj["Upgrade"]; ok {
		var t SuiTransactionUpgrade
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		w.SuiTransaction = t
		return nil
	}

	if _, ok := obj["MakeMoveVec"]; ok {
		var t SuiTransactionMakeMoveVec
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		w.SuiTransaction = t
		return nil
	}

	return errors.New("unknown SuiTransaction type")
}

func (w *SuiTransactionWrapper) MarshalJSON() ([]byte, error) {
	switch t := w.SuiTransaction.(type) {
	case SuiTransactionMoveCall:
		return json.Marshal(t)
	case SuiTransactionTransferObjects:
		return json.Marshal(t)
	case SuiTransactionSplitCoins:
		return json.Marshal(t)
	case SuiTransactionMergeCoins:
		return json.Marshal(t)
	case SuiTransactionPublish:
		return json.Marshal(t)
	case SuiTransactionUpgrade:
		return json.Marshal(t)
	case SuiTransactionMakeMoveVec:
		return json.Marshal(t)
	default:
		return nil, errors.New("unknown SuiTransaction type")
	}
}

// -----------SuiEndOfEpochTransactionKind-----------
type SuiEndOfEpochTransactionKind interface {
	isSuiEndOfEpochTransactionKind()
}
type SuiEndOfEpochTransactionKindAuthenticatorStateCreate string

type SuiEndOfEpochTransactionKindChangeEpoch struct {
	ChangeEpoch SuiChangeEpoch `json:"ChangeEpoch"`
}

type SuiEndOfEpochTransactionKindAuthenticatorStateExpire struct {
	AuthenticatorStateExpire SuiAuthenticatorStateExpire `json:"AuthenticatorStateExpire"`
}

func (SuiEndOfEpochTransactionKindAuthenticatorStateCreate) isSuiEndOfEpochTransactionKind() {}
func (SuiEndOfEpochTransactionKindChangeEpoch) isSuiEndOfEpochTransactionKind()              {}
func (SuiEndOfEpochTransactionKindAuthenticatorStateExpire) isSuiEndOfEpochTransactionKind() {}

type SuiEndOfEpochTransactionKindWrapper struct {
	SuiEndOfEpochTransactionKind
}

func (w *SuiEndOfEpochTransactionKindWrapper) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		w.SuiEndOfEpochTransactionKind = SuiEndOfEpochTransactionKindAuthenticatorStateCreate(s)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	switch {
	case obj["ChangeEpoch"] != nil:
		var k SuiEndOfEpochTransactionKindChangeEpoch
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiEndOfEpochTransactionKind = k
	case obj["AuthenticatorStateExpire"] != nil:
		var k SuiEndOfEpochTransactionKindAuthenticatorStateExpire
		if err := json.Unmarshal(data, &k); err != nil {
			return err
		}
		w.SuiEndOfEpochTransactionKind = k
	default:
		return errors.New("unknown SuiEndOfEpochTransactionKind type")
	}
	return nil
}

func (w *SuiEndOfEpochTransactionKindWrapper) MarshalJSON() ([]byte, error) {
	switch t := w.SuiEndOfEpochTransactionKind.(type) {
	case SuiEndOfEpochTransactionKindAuthenticatorStateCreate:
		return json.Marshal(string(t))
	case SuiEndOfEpochTransactionKindChangeEpoch:
		return json.Marshal(SuiEndOfEpochTransactionKindChangeEpoch{ChangeEpoch: t.ChangeEpoch})
	case SuiEndOfEpochTransactionKindAuthenticatorStateExpire:
		return json.Marshal(SuiEndOfEpochTransactionKindAuthenticatorStateExpire{AuthenticatorStateExpire: t.AuthenticatorStateExpire})
	default:
		return nil, errors.New("unknown SuiEndOfEpochTransactionKind type")
	}
}

// -----------SuiArgument-----------
type SuiArgument interface {
	isSuiArgument()
}

type SuiArgumentGasCoin string

type SuiArgumentInput struct {
	Input uint64 `json:"Input"`
}

type SuiArgumentResult struct {
	Result uint64 `json:"Result"`
}

type SuiArgumentNestedResult struct {
	NestedResult [2]uint64 `json:"NestedResult"`
}

func (SuiArgumentGasCoin) isSuiArgument()      {}
func (SuiArgumentInput) isSuiArgument()        {}
func (SuiArgumentResult) isSuiArgument()       {}
func (SuiArgumentNestedResult) isSuiArgument() {}

type SuiArgumentWrapper struct {
	SuiArgument
}

func (w *SuiArgumentWrapper) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		w.SuiArgument = SuiArgumentGasCoin(s)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _, ok := obj["Input"]; ok {
		var sa SuiArgumentInput
		if err := json.Unmarshal(data, &sa); err != nil {
			return err
		}
		w.SuiArgument = sa
		return nil
	}

	if _, ok := obj["Result"]; ok {
		var sa SuiArgumentResult
		if err := json.Unmarshal(data, &sa); err != nil {
			return err
		}
		w.SuiArgument = sa
		return nil
	}

	if _, ok := obj["NestedResult"]; ok {
		var sa SuiArgumentNestedResult
		if err := json.Unmarshal(data, &sa); err != nil {
			return err
		}
		w.SuiArgument = sa
		return nil
	}

	return errors.New("unknown SuiArgument type")
}

func (w SuiArgumentWrapper) MarshalJSON() ([]byte, error) {
	switch arg := w.SuiArgument.(type) {
	case SuiArgumentGasCoin:
		return json.Marshal(SuiArgumentGasCoin(arg))
	case SuiArgumentInput:
		return json.Marshal(SuiArgumentInput{Input: arg.Input})
	case SuiArgumentResult:
		return json.Marshal(SuiArgumentResult{Result: arg.Result})
	case SuiArgumentNestedResult:
		return json.Marshal(SuiArgumentNestedResult{NestedResult: arg.NestedResult})
	default:
		return nil, errors.New("unknown SuiArgument type")
	}
}

// -----------SuiTransactionArgument-----------
type SuiTransactionArgument interface {
	isSuiTransactionArgument()
}

type SuiTransactionArgumentOne SuiArgumentWrapper
type SuiTransactionArgumentArray []SuiArgumentWrapper
type SuiTransactionArgumentString string
type SuiTransactionArgumentStringArray []string

func (SuiTransactionArgumentOne) isSuiTransactionArgument()         {}
func (SuiTransactionArgumentArray) isSuiTransactionArgument()       {}
func (SuiTransactionArgumentString) isSuiTransactionArgument()      {}
func (SuiTransactionArgumentStringArray) isSuiTransactionArgument() {}

type SuiTransactionArgumentWrapper struct {
	SuiTransactionArgument
}

func (w *SuiTransactionArgumentWrapper) UnmarshalJSON(data []byte) error {
	var one SuiArgumentWrapper
	if err := json.Unmarshal(data, &one); err == nil {
		w.SuiTransactionArgument = SuiTransactionArgumentOne(one)
		return nil
	}

	var array []SuiArgumentWrapper
	if err := json.Unmarshal(data, &array); err == nil {
		w.SuiTransactionArgument = SuiTransactionArgumentArray(array)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		w.SuiTransactionArgument = SuiTransactionArgumentString(s)
		return nil
	}

	var sa []string
	if err := json.Unmarshal(data, &sa); err == nil {
		w.SuiTransactionArgument = SuiTransactionArgumentStringArray(sa)
		return nil
	}

	return errors.New("unknown SuiTransactionArgument type")
}

func (w SuiTransactionArgumentWrapper) MarshalJSON() ([]byte, error) {
	switch arg := w.SuiTransactionArgument.(type) {
	case SuiTransactionArgumentOne:
		return json.Marshal(SuiArgumentWrapper(arg))
	case SuiTransactionArgumentArray:
		return json.Marshal([]SuiArgumentWrapper(arg))
	case SuiTransactionArgumentString:
		return json.Marshal(SuiTransactionArgumentString(arg))
	case SuiTransactionArgumentStringArray:
		return json.Marshal(SuiTransactionArgumentStringArray(arg))
	default:
		return nil, errors.New("unknown SuiTransactionArgument type")
	}
}
