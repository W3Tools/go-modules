package types

import (
	"encoding/json"
	"errors"
)

type SuiParsedData interface {
	isSuiParsedData()
}

type SuiParsedMoveObjectData struct {
	DataType          string            `json:"dataType"`
	Type              string            `json:"type"`
	HasPublicTransfer bool              `json:"hasPublicTransfer"`
	Fields            MoveStructWrapper `json:"fields"`
}

type SuiParsedPackageData struct {
	DataType     string                  `json:"dataType"`
	Disassembled *map[string]interface{} `json:"disassembled,omitempty"`
}

func (SuiParsedMoveObjectData) isSuiParsedData() {}
func (SuiParsedPackageData) isSuiParsedData()    {}

type SuiParsedDataWrapper struct {
	SuiParsedData
}

func (w *SuiParsedDataWrapper) UnmarshalJSON(data []byte) error {
	type DataType struct {
		DataType string `json:"dataType"`
	}

	var dataType DataType
	if err := json.Unmarshal(data, &dataType); err != nil {
		return err
	}

	switch dataType.DataType {
	case "moveObject":
		var p SuiParsedMoveObjectData
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		w.SuiParsedData = p
		return nil
	case "package":
		var p SuiParsedPackageData
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		w.SuiParsedData = p
		return nil
	default:
		return errors.New("unknown SuiParsedData type")
	}
}

func (w *SuiParsedDataWrapper) MarshalJSON() ([]byte, error) {
	switch data := w.SuiParsedData.(type) {
	case SuiParsedMoveObjectData:
		return json.Marshal(SuiParsedMoveObjectData{DataType: data.DataType, Type: data.Type, HasPublicTransfer: data.HasPublicTransfer, Fields: data.Fields})
	case SuiParsedPackageData:
		return json.Marshal(SuiParsedPackageData{DataType: data.DataType, Disassembled: data.Disassembled})
	default:
		return nil, errors.New("unknown SuiParsedData type")
	}
}
