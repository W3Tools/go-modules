package gmsui

import (
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"strings"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
)

type SuiSigner struct {
	Signer *signer.Signer
	Gas    *SuiGasObject
}

// Create New Signer
func NewSuiSignerFromBase64PrivateKey(b64PriKey string) (*SuiSigner, error) {
	priKey, err := B64ToSuiPrivateKey(b64PriKey)
	if err != nil {
		return nil, err
	}

	seed, err := hexutil.Decode(priKey)
	if err != nil {
		return nil, err
	}

	return NewSuiSignerFromSeed(seed), nil
}

func NewSuiSignerFromPrivateKey(priKey string) (*SuiSigner, error) {
	seed, err := hexutil.Decode(priKey)
	if err != nil {
		return nil, err
	}

	return NewSuiSignerFromSeed(seed), nil
}

func NewSuiSignerFromMnemonic(mnemonic string, derivePath string) (*SuiSigner, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	if strings.EqualFold(derivePath, "") {
		derivePath = "m/44'/784'/0'/0'/0'"
	}

	key, err := signer.DeriveForPath(derivePath, seed)
	if err != nil {
		return nil, err
	}

	return NewSuiSignerFromSeed(key.Key), nil
}

func NewSuiSignerFromSeed(seed []byte) *SuiSigner {
	return &SuiSigner{
		Signer: signer.NewSigner(seed),
		Gas:    &SuiGasObject{},
	}
}

// Instance Function
func (s *SuiSigner) GetPulbicKey() string {
	return Ed25519PublicKeyToB64PublicKey(s.Signer.PubKey)
}

func (s *SuiSigner) GetAddress() string {
	return s.Signer.Address
}

func (s *SuiSigner) SignTransaction(b64TxBytes string) (*SuiSignedTransactionRet, error) {
	return s.SignMessage(b64TxBytes, TransactionDataIntentScope)
}

func (s *SuiSigner) SignPersonalMessage(message string) (*SuiSignedMessageRet, error) {
	b64Message := base64.StdEncoding.EncodeToString([]byte(message))
	return s.SignMessage(b64Message, PersonalMessageIntentScope)
}

func (s *SuiSigner) SignMessage(data string, scope IntentScope) (*SuiSignedDataRet, error) {
	txBytes, _ := base64.StdEncoding.DecodeString(data)
	message := NewSuiMessageWithIntent(txBytes, scope)
	digest := blake2b.Sum256(message)
	var noHash crypto.Hash
	sigBytes, err := s.Signer.PriKey.Sign(nil, digest[:], noHash)
	if err != nil {
		return nil, err
	}

	ret := &SuiSignedDataRet{
		TxBytes:   data,
		Signature: toSerializedSignature(sigBytes, s.Signer.PriKey.Public().(ed25519.PublicKey)),
	}
	return ret, nil
}
