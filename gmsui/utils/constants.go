package utils

import "fmt"

var (
	SUI_DECIMALS                      = 9
	MIST_PER_SUI                      = 1000000000
	GAS_SAFE_OVERHEAD          uint64 = 1000
	MAX_GAS                    uint64 = 50000000000
	MOVE_STDLIB_ADDRESS               = NormalizeSuiObjectId("0x1")
	SUI_FRAMEWORK_ADDRESS             = NormalizeSuiObjectId("0x2")
	SUI_SYSTEM_ADDRESS                = NormalizeSuiObjectId("0x3")
	SUI_CLOCK_OBJECT_ID               = NormalizeSuiObjectId("0x6")
	SUI_SYSTEM_MODULE_NAME            = "sui_system"
	SUI_TYPE_ARG                      = fmt.Sprintf("%s::sui::SUI", SUI_FRAMEWORK_ADDRESS)
	SUI_SYSTEM_STATE_OBJECT_ID        = NormalizeSuiObjectId("0x5")
)
