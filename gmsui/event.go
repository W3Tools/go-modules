package gmsui

import (
	"encoding/json"
	"fmt"

	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/go-sui-sdk/v2/types"
)

// Define the type as MoveModule/MoveEventModule. Events emitted, defined on the specified Move module.
// Reference: https://docs.sui.io/guides/developer/sui-101/using-events#filtering-event-queries
type MoveEventModuleConfig struct {
	Package string `toml:"Package,omitempty"`
	Module  string `toml:"Module,omitempty"`
}

func (ec *MoveEventModuleConfig) Join() string {
	return fmt.Sprintf("%s::%s", ec.Package, ec.Module)
}

func (ec *MoveEventModuleConfig) JoinEventName(name string) string {
	return fmt.Sprintf("%s::%s::%s", ec.Package, ec.Module, name)
}

func (ec *MoveEventModuleConfig) PackageToAddress() *move_types.AccountAddress {
	address, err := sui_types.NewObjectIdFromHex(ec.Package)
	if err != nil {
		return nil
	}
	return address
}

// Parsing custom event json
// Reference: https://docs.sui.io/guides/developer/sui-101/using-events#move-event-structure
func ParseEvent[T any](event types.SuiEvent) (*T, error) {
	jsonBytes, err := json.Marshal(event.ParsedJson)
	if err != nil {
		return nil, err
	}

	parsedJson := new(T)
	if err := json.Unmarshal(jsonBytes, &parsedJson); err != nil {
		return nil, err
	}
	return parsedJson, nil
}