package gmsui

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/blake2b"
)

const (
	MaxSignerInMultisig = 10
)

type SuiMultiSig struct {
	Threshold    uint16
	Address      string
	PublicKeyMap []SuiPubkeyWeightPair
	Gas          *SuiGasObject
}

type SuiPubkeyWeightPair struct {
	PublicKey string
	Weight    uint8
}

func NewSuiMultiSig(pubKeys []string, weights []uint8, threshold uint16) (*SuiMultiSig, error) {
	if len(pubKeys) != len(weights) {
		return nil, fmt.Errorf("public key length mismatch")
	}

	if len(pubKeys) > MaxSignerInMultisig {
		return nil, fmt.Errorf("max number of signers in a multisig is %v", MaxSignerInMultisig)
	}

	var keyMap []SuiPubkeyWeightPair
	for idx, pubKey := range pubKeys {
		keyMap = append(keyMap, SuiPubkeyWeightPair{
			PublicKey: pubKey,
			Weight:    weights[idx],
		})
	}

	ret := &SuiMultiSig{
		Threshold:    threshold,
		PublicKeyMap: keyMap,
		Gas:          &SuiGasObject{},
	}

	address, err := ret.ToMultiSigAddress()
	if err != nil {
		return nil, err
	}
	ret.Address = address

	return ret, nil
}

func (m *SuiMultiSig) ToMultiSigAddress() (string, error) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0x03)

	threshold, err := uint16ToBuffer(m.Threshold)
	if err != nil {
		return "", err
	}
	buffer.Write(threshold.Bytes())

	for _, key := range m.PublicKeyMap {
		pubKeyBytes, err := base64.StdEncoding.DecodeString(key.PublicKey)
		if len(pubKeyBytes) != 33 || err != nil {
			return "", fmt.Errorf("public key length error")
		}
		buffer.Write(pubKeyBytes)
		buffer.WriteByte(key.Weight)
	}
	addressBytes := blake2b.Sum256(buffer.Bytes())
	return fmt.Sprintf("0x%s", hex.EncodeToString(addressBytes[:])[:64]), nil
}

func (m *SuiMultiSig) CombineSignatures(signatures []string) (string, error) {
	cli := NewSuiSignatureCombineClient(m.PublicKeyMap, m.Threshold)

	res, err := cli.TryGetCombineSignatures(signatures)
	if err != nil {
		return "", err
	}

	if !strings.EqualFold(strings.ToLower(m.Address), strings.ToLower(res.Address)) {
		return "", fmt.Errorf("signature address is inconsistent")
	}

	return res.Serialized, nil
}

type SuiMultiSigInfo struct {
	Address   string
	Threshold uint16
	Signers   []SuiMultiSigInfoSigner
}

type SuiMultiSigInfoSigner struct {
	Address      string
	B64PublicKey string
	HexPublicKey string
	Weight       uint8
}

func (m *SuiMultiSig) Info() *SuiMultiSigInfo {
	info := &SuiMultiSigInfo{
		Address:   m.Address,
		Threshold: m.Threshold,
	}

	for _, key := range m.PublicKeyMap {
		hexKey, _ := B64ToSuiPrivateKey(key.PublicKey)
		address, _ := B64PublicKeyToSuiAddress(key.PublicKey)
		info.Signers = append(info.Signers, SuiMultiSigInfoSigner{
			Weight:       key.Weight,
			B64PublicKey: key.PublicKey,
			HexPublicKey: hexKey,
			Address:      address,
		})
	}
	return info
}
