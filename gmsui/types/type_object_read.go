package types

import (
	"encoding/json"
	"errors"
)

type ObjectRead interface {
	isObjectRead()
}

type ObjectReadVersionFound struct {
	Details SuiObjectData `json:"detail"`
	Status  string        `json:"status"`
}

type ObjectReadObjectNotExists struct {
	Details string `json:"detail"`
	Status  string `json:"status"`
}

type ObjectReadObjectDeleted struct {
	Details SuiObjectRef `json:"detail"`
	Status  string       `json:"status"`
}

type ObjectReadVersionNotFound struct {
	Details [2]string `json:"detail"`
	Status  string    `json:"status"`
}

type ObjectReadVersionTooHigh struct {
	Details VersionTooHighDetails `json:"detail"`
	Status  string                `json:"status"`
}

type VersionTooHighDetails struct {
	AskedVersion  string `json:"asked_version"`
	LatestVersion string `json:"latest_version"`
	ObjectID      string `json:"object_id"`
}

func (ObjectReadVersionFound) isObjectRead()    {}
func (ObjectReadObjectNotExists) isObjectRead() {}
func (ObjectReadObjectDeleted) isObjectRead()   {}
func (ObjectReadVersionNotFound) isObjectRead() {}
func (ObjectReadVersionTooHigh) isObjectRead()  {}

type ObjectReadWrapper struct {
	ObjectRead
}

func (w *ObjectReadWrapper) UnmarshalJSON(data []byte) error {
	type Status struct {
		Status string `json:"status"`
	}

	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	switch status.Status {
	case "VersionFound":
		var or ObjectReadVersionFound
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "ObjectNotExists":
		var or ObjectReadObjectNotExists
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "ObjectDeleted":
		var or ObjectReadObjectDeleted
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "VersionNotFound":
		var or ObjectReadVersionNotFound
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "VersionTooHigh":
		var or ObjectReadVersionTooHigh
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	default:
		return errors.New("unknown ObjectRead type")
	}

	return nil
}

func (w ObjectReadWrapper) MarshalJSON() ([]byte, error) {
	switch obj := w.ObjectRead.(type) {
	case ObjectReadVersionFound:
		return json.Marshal(ObjectReadVersionFound{Details: obj.Details, Status: obj.Status})
	case ObjectReadObjectNotExists:
		return json.Marshal(ObjectReadObjectNotExists{Details: obj.Details, Status: obj.Status})
	case ObjectReadObjectDeleted:
		return json.Marshal(ObjectReadObjectDeleted{Details: obj.Details, Status: obj.Status})
	case ObjectReadVersionNotFound:
		return json.Marshal(ObjectReadVersionNotFound{Details: obj.Details, Status: obj.Status})
	case ObjectReadVersionTooHigh:
		return json.Marshal(ObjectReadVersionTooHigh{Details: obj.Details, Status: obj.Status})
	default:
		return nil, errors.New("unknown ObjectRead type")
	}
}
