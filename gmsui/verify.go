package gmsui

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/blake2b"
)

func VerifyPersonalMessage(message string, signature string) (signer string, pass bool, err error) {
	bcsData, _ := bcs.Marshal(message)
	b64Message := base64.StdEncoding.EncodeToString(bcsData)
	return VerifyMessage(b64Message, signature, PersonalMessageIntentScope)
}

func VerifyTransactionMessage(b64Message string, signature string) (signer string, pass bool, err error) {
	return VerifyMessage(b64Message, signature, TransactionDataIntentScope)
}

func VerifyMessage(message, signature string, scope IntentScope) (signer string, pass bool, err error) {
	b64Bytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", false, fmt.Errorf("base64.StdEncoding.DecodeString %v", err)
	}

	messageBytes := NewSuiMessageWithIntent(b64Bytes, scope)

	serializedSignature, err := FromSerializedSignature(signature)
	if err != nil {
		return "", false, err
	}
	digest := blake2b.Sum256(messageBytes)

	pass = ed25519.Verify(serializedSignature.PubKey[:], digest[:], serializedSignature.Signature)

	pubKey := Ed25519PublicKeyToB64PublicKey(serializedSignature.PubKey)
	signer, err = B64PublicKeyToSuiAddress(pubKey)
	if err != nil {
		return "", false, fmt.Errorf("invalid signer %v", err)
	}

	return
}
