package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/crypto/blake2b"
)

var (
	_ cryptography.Keypair = (*Secp256r1Keypair)(nil)
)

const DefaultSecp256r1DerivationPath = "m/74'/784'/0'/0/0"

// Secp256r1 Keypair data
type Secp256r1KeypairData struct {
	PublicKey []byte
	SecretKey []byte
}

// An Secp256r1 Keypair used for signing transactions.
type Secp256r1Keypair struct {
	keypair *Secp256r1KeypairData
	cryptography.BaseKeypair
}

// Create or generate random keypair instance.
func NewSecp256r1Keypair(keypair *Secp256r1KeypairData) (*Secp256r1Keypair, error) {
	if keypair == nil {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		publicKey := elliptic.MarshalCompressed(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)
		keypair = &Secp256r1KeypairData{PublicKey: publicKey, SecretKey: privateKey.D.Bytes()}
	}

	kp := &Secp256r1Keypair{keypair: keypair}
	kp.SetSelf(kp)
	return kp, nil
}

// Get the key scheme of the keypair Secp256r1
func (k *Secp256r1Keypair) GetKeyScheme() cryptography.SignatureScheme {
	return cryptography.Secp256r1Scheme
}

// The public key for this keypair
func (k *Secp256r1Keypair) GetPublicKey() (cryptography.PublicKey, error) {
	return NewSecp256r1PublicKey(k.keypair.PublicKey)
}

// The Bech32 secret key string for this Secp256r1 keypair
func (k *Secp256r1Keypair) GetSecretKey() (string, error) {
	return cryptography.EncodeSuiPrivateKey(k.keypair.SecretKey, k.GetKeyScheme())
}

func (k *Secp256r1Keypair) Sign(data []byte) ([]byte, error) {
	return k.SignData(data)
}

// Return the signature for the provided data.
func (k *Secp256r1Keypair) SignData(data []byte) ([]byte, error) {
	hexMessage := sha256.Sum256(data)

	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.D = new(big.Int).SetBytes(k.keypair.SecretKey)
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privKey.D.Bytes())
	r, s, err := deterministicSign(privKey, hexMessage[:])
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
}

// Generate a new random keypair
func GenerateSecp256r1Keypair() (*Secp256r1Keypair, error) {
	return NewSecp256r1Keypair(nil)
}

// Create a keypair from a raw secret key byte array.
// This method should only be used to recreate a keypair from a previously generated secret key.
// Generating keypairs from a random seed should be done with the {@link Keypair.fromSeed} method.
func FromSecretKey(secretKey []byte, skipValidation bool) (*Secp256r1Keypair, error) {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(secretKey)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())
	publicKey := elliptic.MarshalCompressed(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)

	if !skipValidation {
		signData := []byte("sui validation")
		hash := blake2b.Sum256(signData)
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
		if err != nil {
			return nil, err
		}

		if !ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s) {
			return nil, errors.New("provided secretKey is invalid")
		}
	}
	return NewSecp256r1Keypair(&Secp256r1KeypairData{PublicKey: publicKey, SecretKey: secretKey})
}

// Generate a keypair from a 32 byte seed.
func FromSeed(seed []byte) (*Secp256r1Keypair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), strings.NewReader(string(seed)))
	if err != nil {
		return nil, err
	}
	publicKey := elliptic.MarshalCompressed(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	secretKey := privateKey.D.Bytes()

	return NewSecp256r1Keypair(&Secp256r1KeypairData{PublicKey: publicKey, SecretKey: secretKey})
}

// Derive Secp256r1 keypair from mnemonics and path. The mnemonics must be normalized and validated against the english wordlist.
// If path is none, it will default to m/74'/784'/0'/0/0
// Otherwise the path must be compliant to BIP-32 in form m/74'/784'/{account_index}'/{change_index}/{address_index}.
func DeriveKeypair(mnemonics string, path string) (*Secp256r1Keypair, error) {
	if path == "" {
		path = DefaultSecp256r1DerivationPath
	}

	if !cryptography.IsValidBIP32Path(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}
	seed, err := cryptography.MnemonicToSeed(mnemonics)
	if err != nil {
		return nil, err
	}

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	key, err := deriveChildKeyFromPath(masterKey, path)
	if err != nil {
		return nil, err
	}
	return FromSecretKey(key.Key, false)
}

func deriveChildKeyFromPath(masterKey *bip32.Key, path string) (*bip32.Key, error) {
	segments := strings.Split(path, "/")[1:]
	key := masterKey

	for _, segment := range segments {
		if strings.HasSuffix(segment, "'") {
			index, err := strconv.Atoi(strings.TrimSuffix(segment, "'"))
			if err != nil {
				return nil, err
			}
			key, err = key.NewChildKey(bip32.FirstHardenedChild + uint32(index))
			if err != nil {
				return nil, err
			}
		} else {
			index, err := strconv.Atoi(segment)
			if err != nil {
				return nil, err
			}
			key, err = key.NewChildKey(uint32(index))
			if err != nil {
				return nil, err
			}
		}
	}

	return key, nil
}
