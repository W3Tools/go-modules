package secp256k1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/blake2b"
)

var (
	_ cryptography.Keypair = (*Secp256k1Keypair)(nil)
)

const DefaultSecp256k1DerivationPath = "m/54'/784'/0'/0/0"

// Secp256k1 Keypair data
type Secp256k1KeypairData struct {
	PublicKey []byte
	SecretKey []byte
}

// An Secp256k1 Keypair used for signing transactions.
type Secp256k1Keypair struct {
	keypair *Secp256k1KeypairData
	cryptography.BaseKeypair
}

// Create or generate random keypair instance.
func NewSecp256k1Keypair(keypair *Secp256k1KeypairData) (*Secp256k1Keypair, error) {
	if keypair == nil {
		privateKey, err := btcec.NewPrivateKey()
		if err != nil {
			return nil, err
		}
		publicKey := privateKey.PubKey()
		keypair = &Secp256k1KeypairData{PublicKey: publicKey.SerializeCompressed(), SecretKey: privateKey.Serialize()}
	}

	kp := &Secp256k1Keypair{keypair: keypair}
	kp.SetSelf(kp)
	return kp, nil
}

// Get the key scheme of the keypair Secp256k1
func (k *Secp256k1Keypair) GetKeyScheme() cryptography.SignatureScheme {
	return cryptography.Secp256k1Scheme
}

// The public key for this keypair
func (k *Secp256k1Keypair) GetPublicKey() (cryptography.PublicKey, error) {
	return NewSecp256k1PublicKey(k.keypair.PublicKey)
}

// The Bech32 secret key string for this Secp256k1 keypair
func (k *Secp256k1Keypair) GetSecretKey() (string, error) {
	return cryptography.EncodeSuiPrivateKey(k.keypair.SecretKey, k.GetKeyScheme())
}

func (k *Secp256k1Keypair) Sign(data []byte) ([]byte, error) {
	return k.SignData(data)
}

// Return the signature for the provided data.
func (k *Secp256k1Keypair) SignData(data []byte) ([]byte, error) {
	hexMessage := sha256.Sum256(data)

	privKey, _ := btcec.PrivKeyFromBytes(k.keypair.SecretKey)

	r, s, err := deterministicSign(privKey, hexMessage[:])
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
}

// Generate a new random keypair
func GenerateSecp256k1Keypair() (*Secp256k1Keypair, error) {
	return NewSecp256k1Keypair(nil)
}

// Create a keypair from a raw secret key byte array.
// This method should only be used to recreate a keypair from a previously generated secret key.
// Generating keypairs from a random seed should be done with the {@link Keypair.fromSeed} method.
func FromSecretKey(secretKey []byte, skipValidation bool) (*Secp256k1Keypair, error) {
	privKey, _ := btcec.PrivKeyFromBytes(secretKey)
	pubKey := privKey.PubKey().SerializeCompressed()

	if !skipValidation {
		signData := []byte("sui validation")
		hash := blake2b.Sum256(signData)

		r, s, err := ecdsa.Sign(rand.Reader, privKey.ToECDSA(), hash[:])
		if err != nil {
			return nil, fmt.Errorf("failed to sign message hash: %v", err)
		}
		if !ecdsa.Verify(privKey.PubKey().ToECDSA(), hash[:], r, s) {
			return nil, errors.New("provided secretKey is invalid")
		}
	}

	return NewSecp256k1Keypair(&Secp256k1KeypairData{PublicKey: pubKey, SecretKey: secretKey})
}

// Generate a keypair from a 32 byte seed.
func FromSeed(seed []byte) (*Secp256k1Keypair, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	privKey.D = new(big.Int).SetBytes(seed)
	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)

	return NewSecp256k1Keypair(&Secp256k1KeypairData{PublicKey: pubKey, SecretKey: privKey.D.Bytes()})
}

// Derive Secp256k1 keypair from mnemonics and path. The mnemonics must be normalized and validated against the english wordlist.
// If path is none, it will default to m/54'/784'/0'/0/0
// Otherwise the path must be compliant to BIP-32 in form m/54'/784'/{account_index}'/{change_index}/{address_index}.
func DeriveKeypair(mnemonics string, path string) (*Secp256k1Keypair, error) {
	if path == "" {
		path = DefaultSecp256k1DerivationPath
	}

	if !cryptography.IsValidBIP32Path(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}
	seed, err := cryptography.MnemonicToSeed(mnemonics)
	if err != nil {
		return nil, err
	}

	// Derive master key from seed
	master, err := hdkeychain.NewMaster(seed, &chaincfg.Params{})
	if err != nil {
		return nil, fmt.Errorf("failed to derive master key: %v", err)
	}

	// Derive child key using the specified path
	child, err := deriveChildKey(master, path)
	if err != nil {
		return nil, fmt.Errorf("failed to derive child key: %v", err)
	}

	privKey, err := child.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %v", err)
	}

	return FromSecretKey(privKey.Serialize(), false)
}

func deriveChildKey(masterKey *hdkeychain.ExtendedKey, path string) (*hdkeychain.ExtendedKey, error) {
	// Parse BIP32 path
	segments := strings.Split(path, "/")[1:]
	key := masterKey
	for _, part := range segments {
		var index uint32
		var err error
		if part[len(part)-1] == '\'' {
			index, err = parseHardenedIndex(part[:len(part)-1])
		} else {
			index, err = parseIndex(part)
		}
		if err != nil {
			return nil, err
		}
		key, err = key.Derive(index)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

// Helper function to parse a normal index
func parseIndex(part string) (uint32, error) {
	var index uint32
	_, err := fmt.Sscanf(part, "%d", &index)
	return index, err
}

// Helper function to parse a hardened index
func parseHardenedIndex(part string) (uint32, error) {
	index, err := parseIndex(part)
	if err != nil {
		return 0, err
	}
	return index + hdkeychain.HardenedKeyStart, nil
}
