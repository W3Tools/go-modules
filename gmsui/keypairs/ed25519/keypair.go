package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
)

var (
	_ cryptography.Keypair = (*Ed25519Keypair)(nil)
)

const DefaultEd25519DerivationPath = "m/44'/784'/0'/0'/0'"

// Ed25519 Keypair data. The publickey is the 32-byte public key and the secretkey is 64-byte
// Where the first 32 bytes is the secret key and the last 32 bytes is the public key.
type Ed25519KeypairData struct {
	PublicKey []byte
	SecretKey []byte
}

// An Ed25519 Keypair used for signing transactions.
type Ed25519Keypair struct {
	keypair *Ed25519KeypairData
	cryptography.BaseKeypair
}

// Create or generate a new Ed25519 keypair instance.
func NewEd25519Keypair(keypair *Ed25519KeypairData) (*Ed25519Keypair, error) {
	if keypair == nil {
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		keypair = &Ed25519KeypairData{PublicKey: publicKey, SecretKey: privateKey}
	}

	kp := &Ed25519Keypair{keypair: keypair}
	kp.SetSelf(kp)
	return kp, nil
}

// Get the key scheme of the keypair ED25519
func (k *Ed25519Keypair) GetKeyScheme() cryptography.SignatureScheme {
	return cryptography.Ed25519Scheme
}

// The public key for this Ed25519 keypair
func (k *Ed25519Keypair) GetPublicKey() (cryptography.PublicKey, error) {
	return NewEd25519PublicKey(k.keypair.PublicKey)
}

// The Bech32 secret key string for this Ed25519 keypair
func (k *Ed25519Keypair) GetSecretKey() (string, error) {
	return cryptography.EncodeSuiPrivateKey(k.keypair.SecretKey[:cryptography.PrivateKeySize], k.GetKeyScheme())
}

func (k *Ed25519Keypair) Sign(data []byte) ([]byte, error) {
	return k.SignData(data)
}

// Return the signature for the provided data using Ed25519.
func (k *Ed25519Keypair) SignData(data []byte) ([]byte, error) {
	return ed25519.Sign(k.keypair.SecretKey, data), nil
}

// Generate a new random Ed25519 keypair
func GenerateEd25519Keypair() (*Ed25519Keypair, error) {
	return NewEd25519Keypair(nil)
}

// Create a Ed25519 keypair from a raw secret key byte array, also known as seed.
// This is NOT the private scalar which is result of hashing and bit clamping of the raw secret key.
func FromSecretKey(secretKey []byte, skipValidation bool) (*Ed25519Keypair, error) {
	secretKeyLength := len(secretKey)
	if secretKeyLength != cryptography.PrivateKeySize {
		return nil, fmt.Errorf("wrong secretKey size. expected %d bytes, got %d", cryptography.PrivateKeySize, secretKeyLength)
	}

	privateKey := ed25519.NewKeyFromSeed(secretKey)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	if !skipValidation {
		signData := []byte("sui validation")
		signature := ed25519.Sign(privateKey, signData)
		if !ed25519.Verify(publicKey, signData, signature) {
			return nil, errors.New("provided secretKey is invalid")
		}
	}

	return NewEd25519Keypair(&Ed25519KeypairData{PublicKey: publicKey, SecretKey: privateKey})
}

// Derive Ed25519 keypair from mnemonics and path
// The mnemonics must be normalized and validated against the english wordlist.
// If path is none, it will default to m/44'/784'/0'/0'/0'
// Otherwise the path must be compliant to SLIP-0010 in form m/44'/784'/{account_index}'/{change_index}'/{address_index}'.
func DeriveKeypair(mnemonics, path string) (*Ed25519Keypair, error) {
	if path == "" {
		path = DefaultEd25519DerivationPath
	}
	if !cryptography.IsValidHardenedPath(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}

	seedHex, err := cryptography.MnemonicToSeedHex(mnemonics)
	if err != nil {
		return nil, err
	}

	key, err := cryptography.DerivePath(path, seedHex)
	if err != nil {
		return nil, err
	}

	return FromSecretKey(key.Key, false)
}

// Derive Ed25519 keypair from mnemonicSeed and path.
// If path is none, it will default to m/44'/784'/0'/0'/0'
// Otherwise the path must be compliant to SLIP-0010 in form m/44'/784'/{account_index}'/{change_index}'/{address_index}'.
func DeriveKeypairFromSeed(seedHex, path string) (*Ed25519Keypair, error) {
	if path == "" {
		path = DefaultEd25519DerivationPath
	}

	if !cryptography.IsValidHardenedPath(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}

	key, err := cryptography.DerivePath(path, seedHex)
	if err != nil {
		return nil, err
	}
	return FromSecretKey(key.Key, false)
}
