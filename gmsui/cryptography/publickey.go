package cryptography

import (
	"bytes"
	"encoding/hex"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/W3Tools/go-modules/gmsui/utils"
	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/blake2b"
)

type PublicKey interface {
	Equals(publicKey PublicKey) bool
	ToBase64() string
	ToSuiPublicKey() string
	VerifyWithIntent(bs []byte, signature SerializedSignature, intent IntentScope) (bool, error)
	VerifyPersonalMessage(message []byte, signature SerializedSignature) (bool, error)
	VerifyTransactionBlock(transactionBlock []byte, signature SerializedSignature) (bool, error)
	ToSuiBytes() []byte
	ToSuiAddress() string
	ToRawBytes() []byte
	Flag() uint8
	Verify(message []byte, signature SerializedSignature) (bool, error)
}

type BasePublicKey struct {
	self PublicKey
}

// Checks if two public keys are equal
func (base *BasePublicKey) Equals(publicKey PublicKey) bool {
	return gm.BytesEqual(base.self.ToRawBytes(), publicKey.ToRawBytes())
}

// Return the base-64 representation of the public key
func (base *BasePublicKey) ToBase64() string {
	return b64.ToBase64(base.self.ToRawBytes())
}

/*
 * Return the Sui representation of the public key encoded in
 * base-64. A Sui public key is formed by the concatenation
 * of the scheme flag with the raw bytes of the public key
 */
func (base *BasePublicKey) ToSuiPublicKey() string {
	suiBytes := base.ToSuiBytes()
	return b64.ToBase64(suiBytes)
}

func (base *BasePublicKey) VerifyWithIntent(bs []byte, signature SerializedSignature, intent IntentScope) (bool, error) {
	intentMessage := MessageWithIntent(intent, bs)
	digest := blake2b.Sum256(intentMessage)
	return base.self.Verify(digest[:], signature)
}

// Verifies that the signature is valid for for the provided PersonalMessage
func (base *BasePublicKey) VerifyPersonalMessage(message []byte, signature SerializedSignature) (bool, error) {
	msg, err := bcs.Marshal(message)
	if err != nil {
		return false, err
	}
	return base.VerifyWithIntent(msg, signature, PersonalMessage)
}

// Verifies that the signature is valid for for the provided TransactionBlock
func (base *BasePublicKey) VerifyTransactionBlock(transactionBlock []byte, signature SerializedSignature) (bool, error) {
	return base.VerifyWithIntent(transactionBlock, signature, TransactionData)
}

/**
 * Returns the bytes representation of the public key
 * prefixed with the signature scheme flag
 */
func (base *BasePublicKey) ToSuiBytes() []byte {
	rawBytes := base.self.ToRawBytes()
	suiBytes := new(bytes.Buffer)
	suiBytes.Write([]byte{base.self.Flag()})
	suiBytes.Write(rawBytes)

	return suiBytes.Bytes()
}

// Return the Sui address associated with this Ed25519 public key
func (base *BasePublicKey) ToSuiAddress() string {
	digest := blake2b.Sum256(base.ToSuiBytes())
	return utils.NormalizeAddress(hex.EncodeToString(digest[:])[:utils.SuiAddressLength*2])
}

// abstract: Return the byte array representation of the public key
// func (base *BasePublicKey) ToRawBytes() []byte

// abstract: Return signature scheme flag of the public key
// func (base *BasePublicKey) Flag() uint8

// // abstract: Verifies that the signature is valid for for the provided message
//
//	func (base *BasePublicKey) Verify(message []byte, signature SerializedSignature) (bool, error) {
//		return false, nil
//	}

// To ensure thread memory safety, this method needs to be removed. This matter is temporarily added to the todo list
func (base *BasePublicKey) SetSelf(pubkey PublicKey) {
	base.self = pubkey
}
