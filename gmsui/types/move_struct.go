package types

import (
	"encoding/json"
	"errors"
)

type MoveStruct interface {
	isMoveStruct()
	isMoveValue()
}

type MoveStruct_MoveValue []MoveValueWrapper

type MoveStruct_FieldsType struct {
	Type   string                      `json:"type"`
	Fields map[string]MoveValueWrapper `json:"fields"`
}

type MoveStruct_Map map[string]MoveValueWrapper

func (MoveStruct_MoveValue) isMoveStruct()  {}
func (MoveStruct_FieldsType) isMoveStruct() {}
func (MoveStruct_Map) isMoveStruct()        {}

func (MoveStruct_MoveValue) isMoveValue()  {}
func (MoveStruct_FieldsType) isMoveValue() {}
func (MoveStruct_Map) isMoveValue()        {}

type MoveStructWrapper struct {
	MoveStruct
}

func (w *MoveStructWrapper) UnmarshalJSON(data []byte) error {
	var mvs MoveStruct_MoveValue
	if err := json.Unmarshal(data, &mvs); err == nil {
		w.MoveStruct = mvs
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _, ok := obj["fields"]; ok {
		var ms MoveStruct_FieldsType
		if err := json.Unmarshal(data, &ms); err != nil {
			return err
		}
		w.MoveStruct = ms
		return nil
	} else {
		var ms MoveStruct_Map
		if err := json.Unmarshal(data, &ms); err != nil {
			return err
		}
		w.MoveStruct = ms
		return nil
	}
}

func (w MoveStructWrapper) MarshalJSON() ([]byte, error) {
	switch v := w.MoveStruct.(type) {
	case MoveStruct_MoveValue:
		return json.Marshal([]MoveValueWrapper(v))
	case MoveStruct_FieldsType:
		return json.Marshal(MoveStruct_FieldsType{Fields: v.Fields, Type: v.Type})
	case MoveStruct_Map:
		return json.Marshal(v)
	default:
		return nil, errors.New("unknown MoveStruct type")
	}
}

// ---------- Move Value -----------
type MoveValue interface {
	isMoveValue()
}

type MoveNumberValue uint64
type MoveBooleanValue bool
type MoveStringValue string
type MoveValue_MoveValues []MoveValueWrapper
type MoveIdValue struct {
	Id string `json:"id"`
}
type MoveStructValue MoveStruct

func (MoveNumberValue) isMoveValue()      {}
func (MoveBooleanValue) isMoveValue()     {}
func (MoveStringValue) isMoveValue()      {}
func (MoveValue_MoveValues) isMoveValue() {}
func (MoveIdValue) isMoveValue()          {}

type MoveValueWrapper struct {
	MoveValue
}

func (w *MoveValueWrapper) UnmarshalJSON(data []byte) error {
	var num uint64
	if err := json.Unmarshal(data, &num); err == nil {
		w.MoveValue = MoveNumberValue(num)
		return nil
	}
	var bol bool
	if err := json.Unmarshal(data, &bol); err == nil {
		w.MoveValue = MoveBooleanValue(bol)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		w.MoveValue = MoveStringValue(str)
		return nil
	}

	var mvs []MoveValueWrapper
	if err := json.Unmarshal(data, &mvs); err == nil {
		w.MoveValue = MoveValue_MoveValues(mvs)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if _, ok := obj["id"]; ok {
		var mid MoveIdValue
		if err := json.Unmarshal(data, &mid); err == nil {
			w.MoveValue = mid
			return nil
		}
	} else {
		var ms MoveStructWrapper
		if err := json.Unmarshal(data, &ms); err == nil {
			w.MoveValue = MoveStructValue(ms)
			return nil
		}
	}

	return errors.New("unknown MoveValue type")
}

func (w MoveValueWrapper) MarshalJSON() ([]byte, error) {
	switch v := w.MoveValue.(type) {
	case MoveNumberValue:
		return json.Marshal(uint64(v))
	case MoveBooleanValue:
		return json.Marshal(bool(v))
	case MoveStringValue:
		return json.Marshal(string(v))
	case MoveIdValue:
		return json.Marshal(MoveIdValue{Id: v.Id})
	case MoveStructValue:
		return json.Marshal(MoveStruct(v))
	case MoveValue_MoveValues:
		return json.Marshal(v)
	default:
		return nil, errors.New("unknown MoveValue type")
	}
}

// ------------------SuiMoveFunctionArgType------------------
type SuiMoveFunctionArgType interface {
	isSuiMoveFunctionArgType()
}

type SuiMoveFunctionArgStringType string
type SuiMoveFunctionArgObjectType struct {
	Object ObjectValueKind `json:"Object"`
}

func (SuiMoveFunctionArgStringType) isSuiMoveFunctionArgType() {}
func (SuiMoveFunctionArgObjectType) isSuiMoveFunctionArgType() {}

type SuiMoveFunctionArgTypeWrapper struct {
	SuiMoveFunctionArgType
}

func (w *SuiMoveFunctionArgTypeWrapper) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		w.SuiMoveFunctionArgType = SuiMoveFunctionArgStringType(str)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _, ok := obj["Object"]; ok {
		var at SuiMoveFunctionArgObjectType
		if err := json.Unmarshal(data, &at); err != nil {
			return err
		}

		w.SuiMoveFunctionArgType = at
		return nil
	}

	return errors.New("unknown SuiMoveFunctionArgType type")
}

func (w SuiMoveFunctionArgTypeWrapper) MarshalJSON() ([]byte, error) {
	switch v := w.SuiMoveFunctionArgType.(type) {
	case SuiMoveFunctionArgStringType:
		return json.Marshal(string(v))
	case SuiMoveFunctionArgObjectType:
		return json.Marshal(SuiMoveFunctionArgObjectType{Object: v.Object})
	default:
		return nil, errors.New("unknown SuiMoveFunctionArgType type")
	}
}
