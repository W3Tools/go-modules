package verify

import (
	"fmt"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/W3Tools/go-modules/gmsui/keypairs/ed25519"
)

func PublicKeyFromRawBytes(signatureScheme cryptography.SignatureScheme, bs []byte) (cryptography.PublicKey, error) {
	switch signatureScheme {
	case cryptography.Ed25519Scheme:
		return ed25519.NewEd25519PublicKey(bs)
	case cryptography.Secp256k1Scheme:
		return nil, fmt.Errorf("unimplemented %v", signatureScheme)
	case cryptography.Secp256r1Scheme:
		return nil, fmt.Errorf("unimplemented %v", signatureScheme)
	case cryptography.MultiSigScheme:
		return nil, fmt.Errorf("unimplemented %v", signatureScheme)
	case cryptography.ZkLoginScheme:
		return nil, fmt.Errorf("unimplemented %v", signatureScheme)
	default:
		return nil, fmt.Errorf("unsupported signature scheme %v", signatureScheme)
	}
}
