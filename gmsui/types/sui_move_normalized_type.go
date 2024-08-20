package types

import (
	"encoding/json"
	"errors"
)

type SuiMoveNormalizedType interface {
	isSuiMoveNormalizedType()
}

type SuiMoveNormalizedType_String string

type SuiMoveNormalizedType_Struct struct {
	Struct SuiMoveNormalizedTypeStruct `json:"Struct"`
}

type SuiMoveNormalizedTypeStruct struct {
	Address       string                         `json:"address"`
	Module        string                         `json:"module"`
	Name          string                         `json:"name"`
	TypeArguments []SuiMoveNormalizedTypeWrapper `json:"typeArguments"`
}

type SuiMoveNormalizedType_Vector struct {
	Vector SuiMoveNormalizedTypeWrapper `json:"Vector"`
}

type SuiMoveNormalizedType_TypeParameter struct {
	TypeParameter uint64 `json:"TypeParameter"`
}

type SuiMoveNormalizedType_Reference struct {
	Reference SuiMoveNormalizedTypeWrapper `json:"Reference"`
}

type SuiMoveNormalizedType_MutableReference struct {
	MutableReference SuiMoveNormalizedTypeWrapper `json:"MutableReference"`
}

func (SuiMoveNormalizedType_String) isSuiMoveNormalizedType()           {}
func (SuiMoveNormalizedType_Struct) isSuiMoveNormalizedType()           {}
func (SuiMoveNormalizedType_Vector) isSuiMoveNormalizedType()           {}
func (SuiMoveNormalizedType_TypeParameter) isSuiMoveNormalizedType()    {}
func (SuiMoveNormalizedType_Reference) isSuiMoveNormalizedType()        {}
func (SuiMoveNormalizedType_MutableReference) isSuiMoveNormalizedType() {}

type SuiMoveNormalizedTypeWrapper struct {
	SuiMoveNormalizedType
}

func (w *SuiMoveNormalizedTypeWrapper) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		w.SuiMoveNormalizedType = SuiMoveNormalizedType_String(s)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _struct, ok := obj["Struct"]; ok {
		var s SuiMoveNormalizedType_Struct
		if err := json.Unmarshal(_struct, &s.Struct); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if vector, ok := obj["Vector"]; ok {
		var s SuiMoveNormalizedType_Vector
		if err := json.Unmarshal(vector, &s.Vector); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if typeParameter, ok := obj["TypeParameter"]; ok {
		var s SuiMoveNormalizedType_TypeParameter
		if err := json.Unmarshal(typeParameter, &s.TypeParameter); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if reference, ok := obj["Reference"]; ok {
		var s SuiMoveNormalizedType_Reference
		if err := json.Unmarshal(reference, &s.Reference); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if mutableReference, ok := obj["MutableReference"]; ok {
		var s SuiMoveNormalizedType_MutableReference
		if err := json.Unmarshal(mutableReference, &s.MutableReference); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	return errors.New("unknown SuiMoveNormalizedType type")
}

func (w SuiMoveNormalizedTypeWrapper) MarshalJSON() ([]byte, error) {
	switch t := w.SuiMoveNormalizedType.(type) {
	case SuiMoveNormalizedType_String:
		return json.Marshal(string(t))
	case SuiMoveNormalizedType_Struct:
		return json.Marshal(SuiMoveNormalizedType_Struct{Struct: t.Struct})
	case SuiMoveNormalizedType_Vector:
		return json.Marshal(SuiMoveNormalizedType_Vector{Vector: t.Vector})
	case SuiMoveNormalizedType_TypeParameter:
		return json.Marshal(SuiMoveNormalizedType_TypeParameter{TypeParameter: t.TypeParameter})
	case SuiMoveNormalizedType_Reference:
		return json.Marshal(SuiMoveNormalizedType_Reference{Reference: t.Reference})
	case SuiMoveNormalizedType_MutableReference:
		return json.Marshal(SuiMoveNormalizedType_MutableReference{MutableReference: t.MutableReference})
	default:
		return nil, errors.New("unknown SuiMoveNormalizedType")
	}
}
