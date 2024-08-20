package types

import (
	"encoding/json"
	"errors"
)

type InputObjectKind interface {
	isInputObjectKind()
}

type InputObjectKindMovePackage struct {
	MovePackage string `json:"MovePackage"`
}

type InputObjectKindImmOrOwnedMoveObject struct {
	ImmOrOwnedMoveObject SuiObjectRef `json:"ImmOrOwnedMoveObject"`
}

type InputObjectKindSharedMoveObject struct {
	SharedMoveObject KindSharedMoveObject `json:"SharedMoveObject"`
}

type KindSharedMoveObject struct {
	ID                   string `json:"id"`
	InitialSharedVersion string `json:"initial_shared_version"`
	Mutable              bool   `json:"mutable,omitempty"`
}

func (InputObjectKindMovePackage) isInputObjectKind()          {}
func (InputObjectKindImmOrOwnedMoveObject) isInputObjectKind() {}
func (InputObjectKindSharedMoveObject) isInputObjectKind()     {}

type InputObjectKindWrapper struct {
	InputObjectKind
}

func (w *InputObjectKindWrapper) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	switch {
	case obj["MovePackage"] != nil:
		var s InputObjectKindMovePackage
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.InputObjectKind = s
	case obj["ImmOrOwnedMoveObject"] != nil:
		var s InputObjectKindImmOrOwnedMoveObject
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.InputObjectKind = s
	case obj["SharedMoveObject"] != nil:
		var s InputObjectKindSharedMoveObject
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.InputObjectKind = s
	default:
		return errors.New("unknown InputObjectKind type")
	}

	return nil
}

func (w *InputObjectKindWrapper) MarshalJSON() ([]byte, error) {
	switch t := w.InputObjectKind.(type) {
	case InputObjectKindMovePackage:
		return json.Marshal(InputObjectKindMovePackage{
			MovePackage: t.MovePackage,
		})
	case InputObjectKindImmOrOwnedMoveObject:
		return json.Marshal(InputObjectKindImmOrOwnedMoveObject{
			ImmOrOwnedMoveObject: t.ImmOrOwnedMoveObject,
		})
	case InputObjectKindSharedMoveObject:
		return json.Marshal(InputObjectKindSharedMoveObject{
			SharedMoveObject: t.SharedMoveObject,
		})
	default:
		return nil, errors.New("unknown InputObjectKind type")
	}
}
