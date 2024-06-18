package cryptography

import (
	"bytes"
)

type AppId = uint8

const (
	Sui AppId = iota
)

type IntentVersion = uint8

const (
	V0 IntentVersion = iota
)

type IntentScope = uint8

const (
	TransactionData IntentScope = iota
	TransactionEffects
	CheckpointSummary
	PersonalMessage
)

type Intent = []uint8

func IntentWithScope(scope IntentScope) Intent {
	return []uint8{scope, V0, Sui}
}

func MessageWithIntent(scope IntentScope, message []byte) []byte {
	intent := IntentWithScope(scope)
	intentMessage := new(bytes.Buffer)
	intentMessage.Write(intent)
	intentMessage.Write(message)

	return intentMessage.Bytes()
}
