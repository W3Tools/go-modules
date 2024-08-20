package types

import (
	"encoding/json"
	"errors"
)

type ObjectResponseError interface {
	isObjectResponseError()
}

type ObjectResponseNotExistsError struct {
	Code     string `json:"code"`
	ObjectId string `json:"object_id"`
}

type ObjectResponseDynamicFieldNotFoundError struct {
	Code           string `json:"code"`
	ParentObjectId string `json:"parent_object_id"`
}

type ObjectResponseDeletedError struct {
	Code     string `json:"code"`
	Digest   string `json:"digest"`
	ObjectId string `json:"object_id"`
	Version  uint64 `json:"version"`
}

type ObjectResponseUnknownError struct {
	Code string `json:"code"`
}

type ObjectResponseDisplayErrorError struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (ObjectResponseNotExistsError) isObjectResponseError()            {}
func (ObjectResponseDynamicFieldNotFoundError) isObjectResponseError() {}
func (ObjectResponseDeletedError) isObjectResponseError()              {}
func (ObjectResponseUnknownError) isObjectResponseError()              {}
func (ObjectResponseDisplayErrorError) isObjectResponseError()         {}

type ObjectResponseErrorWrapper struct {
	ObjectResponseError
}

func (w *ObjectResponseErrorWrapper) UnmarshalJSON(data []byte) error {
	type ErrorCode struct {
		Code string `json:"code"`
	}

	var errorCode ErrorCode
	if err := json.Unmarshal(data, &errorCode); err != nil {
		return err
	}

	switch errorCode.Code {
	case "notExists":
		var e ObjectResponseNotExistsError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "deleted":
		var e ObjectResponseDeletedError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "dynamicFieldNotFound":
		var e ObjectResponseDynamicFieldNotFoundError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "unknown":
		var e ObjectResponseUnknownError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "displayError":
		var e ObjectResponseDisplayErrorError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	default:
		return errors.New("unknown ObjectResponseError type")
	}
}

func (w *ObjectResponseErrorWrapper) MarshalJSON() ([]byte, error) {
	switch e := w.ObjectResponseError.(type) {
	case ObjectResponseNotExistsError:
		return json.Marshal(ObjectResponseNotExistsError{Code: e.Code, ObjectId: e.ObjectId})
	case ObjectResponseDynamicFieldNotFoundError:
		return json.Marshal(ObjectResponseDynamicFieldNotFoundError{Code: e.Code, ParentObjectId: e.ParentObjectId})
	case ObjectResponseDeletedError:
		return json.Marshal(ObjectResponseDeletedError{Code: e.Code, Digest: e.Digest, ObjectId: e.ObjectId, Version: e.Version})
	case ObjectResponseUnknownError:
		return json.Marshal(ObjectResponseUnknownError{Code: e.Code})
	case ObjectResponseDisplayErrorError:
		return json.Marshal(ObjectResponseDisplayErrorError{Code: e.Code, Error: e.Error})
	default:
		return nil, errors.New("unknown ObjectResponseError type")
	}
}
