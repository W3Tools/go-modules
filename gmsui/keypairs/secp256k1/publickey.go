package secp256k1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/btcsuite/btcd/btcec/v2"
)

var (
	_ cryptography.PublicKey = (*Secp256k1PublicKey)(nil)
)

// A Secp256k1 public key
type Secp256k1PublicKey struct {
	data []byte
	cryptography.BasePublicKey
}

// Create a new Secp256k1PublicKey object
func NewSecp256k1PublicKey[T string | []byte](value T) (publicKey *Secp256k1PublicKey, err error) {
	publicKey = new(Secp256k1PublicKey)
	switch v := any(value).(type) {
	case string:
		publicKey.data, err = b64.FromBase64(v)
		if err != nil {
			return nil, err
		}
	case []byte:
		publicKey.data = v
	}
	if len(publicKey.data) != cryptography.Secp256k1PublicKeySize {
		return nil, fmt.Errorf("invalid public key input, expected %d bytes, got %d", cryptography.Secp256k1PublicKeySize, len(publicKey.data))
	}
	publicKey.SetSelf(publicKey)
	return
}

// Checks if two Secp256k1 public keys are equal
func (k *Secp256k1PublicKey) Equals(publicKey cryptography.PublicKey) bool {
	return k.BasePublicKey.Equals(publicKey)
}

// Return the byte array representation of the Secp256k1 public key
func (k *Secp256k1PublicKey) ToRawBytes() []byte {
	return k.data
}

// Return the Sui address associated with this Secp256k1 public key
func (k *Secp256k1PublicKey) Flag() uint8 {
	return cryptography.SignatureSchemeToFlag[cryptography.Secp256k1Scheme]
}

// Verifies that the signature is valid for for the provided message
func (k *Secp256k1PublicKey) Verify(message []byte, signature cryptography.SerializedSignature) (bool, error) {
	parsed, err := cryptography.ParseSerializedSignature(signature)
	if err != nil {
		return false, err
	}

	if parsed.SignatureScheme != cryptography.Secp256k1Scheme {
		return false, fmt.Errorf("invalid signature scheme")
	}

	if !bytes.Equal(k.ToRawBytes(), parsed.PubKey) {
		return false, fmt.Errorf("signature does not match public key")
	}

	// Parse the signature into r and s components
	r := new(big.Int).SetBytes(parsed.Signature[:32])
	s := new(big.Int).SetBytes(parsed.Signature[32:])

	pubKey, err := btcec.ParsePubKey(k.ToRawBytes())
	if err != nil {
		return false, err
	}

	// Verify the signature
	hash := sha256.Sum256(message)
	verified := ecdsa.Verify(pubKey.ToECDSA(), hash[:], r, s)

	return verified, nil
}
