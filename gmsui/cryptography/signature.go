package cryptography

import (
	"bytes"
	"fmt"

	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/fardream/go-bcs/bcs"
)

/**
 * Pair of signature and corresponding public key
 */
type SerializeSignatureInput struct {
	SignatureScheme SignatureScheme
	Signature       []byte
	PublicKey       PublicKey
}

type SerializedSignature = string

/**
 * Takes in a signature, its associated signing scheme and a public key, then serializes this data
 */
func ToSerializedSignature(input SerializeSignatureInput) (string, error) {
	if input.PublicKey == nil {
		return "", fmt.Errorf("publicKey is required")
	}

	pubKeyBytes := input.PublicKey.ToRawBytes()
	serializedSignature := new(bytes.Buffer)
	serializedSignature.Write([]byte{SignatureSchemeToFlag[input.SignatureScheme]})
	serializedSignature.Write(input.Signature)
	serializedSignature.Write(pubKeyBytes)
	return b64.ToBase64(serializedSignature.Bytes()), nil
}

/**
 * Decodes a serialized signature into its constituent components: the signature scheme, the actual signature, and the public key
 */
func ParseSerializedSignature(serializedSignature SerializedSignature) (*SerializedSignatureParsedData, error) {
	bs, err := b64.FromBase64(serializedSignature)
	if err != nil {
		return nil, err
	}

	signatureScheme := SignatureFlagToScheme[bs[0]]

	switch signatureScheme {
	case Ed25519Scheme, Secp256k1Scheme, Secp256r1Scheme:
		size := SignatureSchemeToSize[signatureScheme]
		signature := bs[1 : len(bs)-size]
		publicKey := bs[1+len(signature):]

		data := &SerializedSignatureParsedData{
			SerializedSignature: serializedSignature,
			SignatureScheme:     signatureScheme,
			Signature:           signature,
			PubKey:              publicKey,
			Bytes:               bs,
		}

		return data, nil
	case MultiSigScheme:
		multisig := new(MultiSigStruct)
		_, err := bcs.Unmarshal(bs[1:], &multisig)
		if err != nil {
			return nil, err
		}

		data := &SerializedSignatureParsedData{
			SerializedSignature: serializedSignature,
			SignatureScheme:     signatureScheme,
			Multisig:            multisig,
			Bytes:               bs,
		}
		return data, nil
	case ZkLoginScheme:
		return nil, fmt.Errorf("unimplemented %v", ZkLoginScheme)
	default:
		return nil, fmt.Errorf("unsupported signature scheme")
	}
}

type SerializedSignatureParsedData struct {
	SerializedSignature SerializedSignature `json:"serializedSignature"`
	SignatureScheme     string              `json:"signatureScheme"`
	Signature           []byte              `json:"signature,omitempty"`
	PubKey              []byte              `json:"pubKey,omitempty"`
	Bytes               []byte              `json:"bytes"`
	Multisig            *MultiSigStruct     `json:"multisig,omitempty"`
}
