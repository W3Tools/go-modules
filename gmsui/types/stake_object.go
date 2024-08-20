package types

import (
	"encoding/json"
	"errors"
)

type StakeObject interface {
	isStakeObject()
}

type StakeObjectPending struct {
	Principal         string `json:"principal"`
	StakeActiveEpoch  string `json:"stakeActiveEpoch"`
	StakeRequestEpoch string `json:"stakeRequestEpoch"`
	StakedSuiId       string `json:"stakedSuiId"`
	Status            string `json:"status"`
}

type StakeObjectActive struct {
	Principal         string `json:"principal"`
	StakeActiveEpoch  string `json:"stakeActiveEpoch"`
	StakeRequestEpoch string `json:"stakeRequestEpoch"`
	StakedSuiId       string `json:"stakedSuiId"`
	EstimatedReward   string `json:"estimatedReward"`
	Status            string `json:"status"`
}

type StakeObjectUnstaked struct {
	Principal         string `json:"principal"`
	StakeActiveEpoch  string `json:"stakeActiveEpoch"`
	StakeRequestEpoch string `json:"stakeRequestEpoch"`
	StakedSuiId       string `json:"stakedSuiId"`
	Status            string `json:"status"`
}

func (StakeObjectPending) isStakeObject()  {}
func (StakeObjectActive) isStakeObject()   {}
func (StakeObjectUnstaked) isStakeObject() {}

type StakeObjectWrapper struct {
	StakeObject
}

func (w *StakeObjectWrapper) UnmarshalJSON(data []byte) error {
	type Status struct {
		Status string `json:"status"`
	}
	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	switch status.Status {
	case "Pending":
		var so StakeObjectPending
		if err := json.Unmarshal(data, &so); err != nil {
			return err
		}
		w.StakeObject = so
	case "Active":
		var so StakeObjectActive
		if err := json.Unmarshal(data, &so); err != nil {
			return err
		}
		w.StakeObject = so
	case "Unstaked":
		var so StakeObjectUnstaked
		if err := json.Unmarshal(data, &so); err != nil {
			return err
		}
		w.StakeObject = so
	default:
		return errors.New("unknown StakeObject type")
	}

	return nil
}

func (w *StakeObjectWrapper) MarshalJSON() ([]byte, error) {
	switch obj := w.StakeObject.(type) {
	case StakeObjectPending:
		return json.Marshal(StakeObjectPending{
			Status:            obj.Status,
			Principal:         obj.Principal,
			StakeActiveEpoch:  obj.StakeActiveEpoch,
			StakeRequestEpoch: obj.StakeRequestEpoch,
			StakedSuiId:       obj.StakedSuiId,
		})
	case StakeObjectActive:
		return json.Marshal(StakeObjectActive{
			Status:            obj.Status,
			Principal:         obj.Principal,
			StakeActiveEpoch:  obj.StakeActiveEpoch,
			StakeRequestEpoch: obj.StakeRequestEpoch,
			StakedSuiId:       obj.StakedSuiId,
			EstimatedReward:   obj.EstimatedReward,
		})
	case StakeObjectUnstaked:
		return json.Marshal(StakeObjectUnstaked{
			Status:            obj.Status,
			Principal:         obj.Principal,
			StakeActiveEpoch:  obj.StakeActiveEpoch,
			StakeRequestEpoch: obj.StakeRequestEpoch,
			StakedSuiId:       obj.StakedSuiId,
		})
	default:
		return nil, errors.New("unknown StakeObject type")
	}
}
