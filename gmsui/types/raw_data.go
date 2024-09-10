package types

import (
	"encoding/json"
	"errors"
)

type RawData interface {
	isRawData()
}

type RawDataMoveObject struct {
	DataType          string `json:"dataType"`
	Type              string `json:"type"`
	HasPublicTransfer bool   `json:"hasPublicTransfer"`
	Version           uint64 `json:"version"`
	BcsBytes          string `json:"bcsBytes"`
}

type RawDataPackage struct {
	DataType        string                 `json:"dataType"`
	ID              string                 `json:"id"`
	Version         uint64                 `json:"version"`
	ModuleMap       map[string]string      `json:"moduleMap"`
	TypeOriginTable []TypeOrigin           `json:"typeOriginTable"`
	LinkageTable    map[string]UpgradeInfo `json:"linkageTable"`
}

func (RawDataMoveObject) isRawData() {}
func (RawDataPackage) isRawData()    {}

type RawDataWrapper struct {
	RawData
}

func (w *RawDataWrapper) UnmarshalJSON(data []byte) error {
	type DataType struct {
		DataType string `json:"dataType"`
	}
	var dataType DataType
	if err := json.Unmarshal(data, &dataType); err != nil {
		return err
	}

	switch dataType.DataType {
	case "moveObject":
		var rd RawDataMoveObject
		if err := json.Unmarshal(data, &rd); err != nil {
			return err
		}
		w.RawData = rd
	case "package":
		var rd RawDataPackage
		if err := json.Unmarshal(data, &rd); err != nil {
			return err
		}
		w.RawData = rd
	default:
		return errors.New("unknown RawData type")
	}

	return nil
}

func (w *RawDataWrapper) MarshalJSON() ([]byte, error) {
	switch rd := w.RawData.(type) {
	case RawDataMoveObject:
		return json.Marshal(RawDataMoveObject{
			BcsBytes:          rd.BcsBytes,
			DataType:          rd.DataType,
			HasPublicTransfer: rd.HasPublicTransfer,
			Type:              rd.Type,
			Version:           rd.Version,
		})
	case RawDataPackage:
		return json.Marshal(RawDataPackage{
			DataType:        rd.DataType,
			ID:              rd.ID,
			LinkageTable:    rd.LinkageTable,
			ModuleMap:       rd.ModuleMap,
			TypeOriginTable: rd.TypeOriginTable,
			Version:         rd.Version,
		})
	default:
		return nil, errors.New("unknown RawData type")
	}
}
