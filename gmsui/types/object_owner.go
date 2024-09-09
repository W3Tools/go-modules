package types

import (
	"encoding/json"
	"errors"
)

type ObjectOwner interface {
	isObjectOwner()
}

type ObjectOwner_AddressOwner struct {
	AddressOwner string `json:"AddressOwner"`
}

type ObjectOwner_ObjectOwner struct {
	ObjectOwner string `json:"ObjectOwner"`
}

type ObjectOwner_Shared struct {
	Shared ObjectOwner_SharedData `json:"Shared"`
}

type ObjectOwner_SharedData struct {
	InitialSharedVersion uint64 `json:"initial_shared_version"`
}

type ObjectOwner_Immutable string

func (ObjectOwner_AddressOwner) isObjectOwner() {}
func (ObjectOwner_ObjectOwner) isObjectOwner()  {}
func (ObjectOwner_Shared) isObjectOwner()       {}
func (ObjectOwner_Immutable) isObjectOwner()    {}

type ObjectOwnerWrapper struct {
	ObjectOwner
}

// UnmarshalJSON custom unmarshaller for ObjectOwnerWrapper
func (w *ObjectOwnerWrapper) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		w.ObjectOwner = ObjectOwner_Immutable(s)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if addressOwner, ok := obj["AddressOwner"]; ok {
		var o ObjectOwner_AddressOwner
		if err := json.Unmarshal(addressOwner, &o.AddressOwner); err != nil {
			return err
		}
		w.ObjectOwner = o
		return nil
	}

	if objectOwner, ok := obj["ObjectOwner"]; ok {
		var o ObjectOwner_ObjectOwner
		if err := json.Unmarshal(objectOwner, &o.ObjectOwner); err != nil {
			return err
		}
		w.ObjectOwner = o
		return nil
	}

	if shared, ok := obj["Shared"]; ok {
		var o ObjectOwner_Shared
		if err := json.Unmarshal(shared, &o.Shared); err != nil {
			return err
		}

		w.ObjectOwner = o
		return nil
	}

	return errors.New("unknown ObjectOwner type")
}

func (w ObjectOwnerWrapper) MarshalJSON() ([]byte, error) {
	switch o := w.ObjectOwner.(type) {
	case ObjectOwner_AddressOwner:
		return json.Marshal(ObjectOwner_AddressOwner{AddressOwner: o.AddressOwner})
	case ObjectOwner_ObjectOwner:
		return json.Marshal(ObjectOwner_ObjectOwner{ObjectOwner: o.ObjectOwner})
	case ObjectOwner_Shared:
		return json.Marshal(ObjectOwner_Shared{Shared: o.Shared})
	case ObjectOwner_Immutable:
		return json.Marshal(string(o))
	default:
		return nil, errors.New("unknown ObjectOwner type")
	}
}
