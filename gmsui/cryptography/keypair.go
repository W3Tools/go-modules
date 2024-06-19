package cryptography

import (
	"fmt"

	"github.com/W3Tools/go-bcs/bcs"
	"github.com/W3Tools/go-modules/gmsui/b64"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"golang.org/x/crypto/blake2b"
)

const (
	PrivateKeySize      = 32
	SuiPrivateKeyPrefix = "suiprivkey"
)

type ParsedKeypair struct {
	Scheme    SignatureScheme
	SecretKey []byte
}

type SignatureWithBytes struct {
	Bytes     string
	Signature SerializedSignature
}

type Signer interface {
	Sign(bs []byte) ([]byte, error)
	SignWithIntent(bs []byte, intent IntentScope) (*SignatureWithBytes, error)
	SignTransactionBlock(bs []byte) (*SignatureWithBytes, error)
	SignPersonalMessage(bs []byte) (*SignatureWithBytes, error)
	ToSuiAddress() string
	SignData(data []byte) ([]byte, error)
	GetKeyScheme() SignatureScheme
	GetPublicKey() (PublicKey, error)
}

type Keypair interface {
	Signer
	GetSecretKey() (string, error)
}

type BaseSigner struct {
	self Signer
}

/**
 * Sign messages with a specific intent. By combining the message bytes with the intent before hashing and signing,
 * it ensures that a signed message is tied to a specific purpose and domain separator is provided
 */
func (signer *BaseSigner) SignWithIntent(bs []byte, intent IntentScope) (*SignatureWithBytes, error) {
	intentMessage := MessageWithIntent(intent, bs)
	digest := blake2b.Sum256(intentMessage)

	publicKey, err := signer.self.GetPublicKey()
	if err != nil {
		return nil, err
	}

	sign, err := signer.self.Sign(digest[:])
	if err != nil {
		return nil, err
	}

	signature, err := ToSerializedSignature(SerializeSignatureInput{
		SignatureScheme: signer.self.GetKeyScheme(),
		Signature:       sign,
		PublicKey:       publicKey,
	})
	if err != nil {
		return nil, err
	}

	return &SignatureWithBytes{Bytes: b64.ToBase64(bs), Signature: signature}, nil
}

/**
 * Signs provided transaction block by calling `signWithIntent()` with a `TransactionData` provided as intent scope
 */
func (signer *BaseSigner) SignTransactionBlock(bs []byte) (*SignatureWithBytes, error) {
	return signer.SignWithIntent(bs, TransactionData)
}

/**
 * Signs provided personal message by calling `signWithIntent()` with a `PersonalMessage` provided as intent scope
 */
func (signer *BaseSigner) SignPersonalMessage(bs []byte) (*SignatureWithBytes, error) {
	msg, err := bcs.Marshal(bs)
	if err != nil {
		return nil, err
	}

	return signer.SignWithIntent(msg, PersonalMessage)
}

func (signer *BaseSigner) ToSuiAddress() string {
	publicKey, _ := signer.self.GetPublicKey()
	return publicKey.ToSuiAddress()
}

// abstract
// func (signer *BaseSigner) Sign(bs []byte) []byte         { return nil }
// func (signer *BaseSigner) SignData(data []byte) []byte   { return nil }
// func (signer *BaseSigner) GetKeyScheme() SignatureScheme { return "" }
// func (signer *BaseSigner) GetPublicKey() PublicKey       { return nil }

type BaseKeypair struct {
	BaseSigner
}

func (base *BaseKeypair) SetSelf(signer Signer) {
	base.self = signer
}

/**
 * This returns an ParsedKeypair object based by validating the
 * 33-byte Bech32 encoded string starting with `suiprivkey`, and
 * parse out the signature scheme and the private key in bytes.
 */
func DecodeSuiPrivateKey(value string) (*ParsedKeypair, error) {
	prefix, words, err := bech32.Decode(value)
	if err != nil {
		return nil, err
	}
	if prefix != SuiPrivateKeyPrefix {
		return nil, fmt.Errorf("invalid private key prefix")
	}
	extendedSecretKey, err := bech32.ConvertBits(words, 5, 8, false)
	if err != nil {
		return nil, err
	}
	secretKey := extendedSecretKey[1:]
	signatureScheme := SignatureFlagToScheme[extendedSecretKey[0]]
	return &ParsedKeypair{Scheme: signatureScheme, SecretKey: secretKey}, nil
}

/**
 * This returns a Bech32 encoded string starting with `suiprivkey`,
 * encoding 33-byte `flag || bytes` for the given the 32-byte private
 * key and its signature scheme.
 */
func EncodeSuiPrivateKey(bytes []byte, scheme string) (string, error) {
	if len(bytes) != PrivateKeySize {
		return "", fmt.Errorf("invalid bytes length")
	}
	flag := SignatureSchemeToFlag[scheme]
	privKeyBytes := append([]byte{flag}, bytes...)
	words, err := bech32.ConvertBits(privKeyBytes, 8, 5, true)
	if err != nil {
		return "", err
	}
	encoded, err := bech32.Encode(SuiPrivateKeyPrefix, words)
	if err != nil {
		return "", err
	}
	return encoded, nil
}
