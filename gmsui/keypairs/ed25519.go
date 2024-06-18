package keypairs

import (
	"crypto/ed25519"
	"fmt"

	gm "github.com/W3Tools/go-modules"
	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/W3Tools/go-modules/gmsui/cryptography"
)

var (
	_ cryptography.PublicKey = (*Ed25519PublicKey)(nil)
)

type Ed25519PublicKey struct {
	data []byte
	cryptography.BasePublicKey
}

const Ed25519PublicKeySize = 32

func NewEd25519PublicKey[T string | []byte](value T) (publicKey *Ed25519PublicKey, err error) {
	publicKey = new(Ed25519PublicKey)
	switch v := any(value).(type) {
	case string:
		publicKey.data, err = b64.FromBase64(v)
		if err != nil {
			return nil, err
		}
	case []byte:
		publicKey.data = v
	}

	if len(publicKey.data) != Ed25519PublicKeySize {
		return nil, fmt.Errorf("invalid public key input. expected %v bytes, got %v", Ed25519PublicKeySize, len(publicKey.data))
	}
	publicKey.SetSelf(publicKey)
	return
}

// Return the byte array representation of the Ed25519 public key
func (key *Ed25519PublicKey) ToRawBytes() []byte {
	return key.data
}

// Return the Sui address associated with this Ed25519 public key
func (key *Ed25519PublicKey) Flag() uint8 {
	return cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme]
}

// Verifies that the signature is valid for for the provided message
func (key *Ed25519PublicKey) Verify(message []byte, signature cryptography.SerializedSignature) (bool, error) {
	parsed, err := cryptography.ParseSerializedSignature(signature)
	if err != nil {
		return false, err
	}

	if parsed.SignatureScheme != cryptography.Ed25519Scheme {
		return false, err
	}

	if !gm.BytesEqual(key.ToRawBytes(), parsed.PubKey) {
		fmt.Printf("123\n")
		return false, err
	}

	return ed25519.Verify(parsed.PubKey, message, parsed.Signature), nil
}
