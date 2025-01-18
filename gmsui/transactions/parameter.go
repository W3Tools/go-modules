package transactions

import (
	"fmt"

	"github.com/W3Tools/go-modules/gmsui/client"
	"github.com/W3Tools/go-modules/gmsui/types"
	"github.com/W3Tools/go-modules/gmsui/utils"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
)

type UnresolvedParameter struct {
	Arguments UnresolvedArguments      `json:"argument"`
	Objects   map[int]UnresolvedObject `json:"object"` // map key is index of argument
}

type UnresolvedArgument struct {
	Pure     any
	Object   *sui_types.ObjectArg
	Argument *sui_types.Argument
}
type UnresolvedArguments []*UnresolvedArgument

type UnresolvedObject struct {
	ObjectId string
	Mutable  bool
}

func NewUnresolvedParameter(count int) *UnresolvedParameter {
	return &UnresolvedParameter{
		Arguments: make(UnresolvedArguments, count),
		Objects:   make(map[int]UnresolvedObject),
	}
}

func (up *UnresolvedParameter) merge(dest *UnresolvedParameter) {
	if dest == nil {
		return
	}

	count := len(up.Arguments)

	up.Arguments = append(up.Arguments, dest.Arguments...)
	for idx, obj := range dest.Objects {
		up.Objects[idx+count] = obj
	}
}

func (up *UnresolvedParameter) resolveAndPArseToArguments(suiClient *client.SuiClient, txb *Transaction) ([]sui_types.Argument, error) {
	err := up.resolveObjects(suiClient)
	if err != nil {
		return nil, fmt.Errorf("can not resolve objects: %v", err)
	}

	return up.toArguments(txb)
}

func (up *UnresolvedParameter) resolveObjects(suiClient *client.SuiClient) error {
	if len(up.Objects) > 0 {
		var ids []string
		for _, resolve := range up.Objects {
			ids = append(ids, resolve.ObjectId)
		}

		var objects []*types.SuiObjectResponse
		objects, err := suiClient.MultiGetObjects(types.MultiGetObjectsParams{IDs: ids, Options: &types.SuiObjectDataOptions{ShowOwner: true}})
		if err != nil {
			return fmt.Errorf("can not call jsonrpc to multi get objects: %v", err)
		}

		objectMap := utils.SliceToMap(objects, func(v *types.SuiObjectResponse) string {
			if v.Data != nil {
				return v.Data.ObjectId
			}
			return ""
		})

		for idx, resolveObject := range up.Objects {
			if idx >= len(up.Arguments) {
				return fmt.Errorf("can not resolve object at index %d, out of range", idx)
			}

			object := objectMap[utils.NormalizeSuiObjectId(resolveObject.ObjectId)]
			if object == nil {
				return fmt.Errorf("can not fetch object with id [%s] at index %d", resolveObject.ObjectId, idx)
			}

			objectArg, err := objectResponseToObjectArg(object, resolveObject.Mutable)
			if err != nil {
				return fmt.Errorf("can not convert object response to object arg at index %d: %v", idx, err)
			}

			up.Arguments[idx] = &UnresolvedArgument{Object: objectArg}
		}
	}

	return nil
}

func (up *UnresolvedParameter) toArguments(txb *Transaction) ([]sui_types.Argument, error) {
	arguments := make([]sui_types.Argument, len(up.Arguments))
	for idx, input := range up.Arguments {
		if input.Pure != nil {
			value, err := txb.builder.Pure(input.Pure)
			if err != nil {
				return nil, fmt.Errorf("can not create pure argument at index %d: %v", idx, err)
			}
			arguments[idx] = value
		} else if input.Object != nil {
			value, err := txb.builder.Obj(*input.Object)
			if err != nil {
				return nil, fmt.Errorf("can not create object argument at index %d: %v", idx, err)
			}
			arguments[idx] = value
		} else if input.Argument != nil {
			arguments[idx] = *input.Argument
		} else {
			return nil, fmt.Errorf("invalid input argument at index: %v, pure: %v, object: %v, argument: %v", idx, input.Pure, input.Object, input.Argument)
		}
	}

	return arguments, nil
}
