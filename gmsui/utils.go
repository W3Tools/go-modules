package gmsui

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/blake2b"
)

// Utils
func B64ToSuiPrivateKey(b64 string) (string, error) {
	b64Decode, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	hexPriKey := hexutil.Encode(b64Decode)
	if len(hexPriKey) != 68 {
		return "", fmt.Errorf("unknown base64. %s", b64)
	}
	return fmt.Sprintf("0x%s", hexPriKey[4:]), nil
}

func SuiPrivateKeyToB64(pk string) (string, error) {
	if len(pk) != 66 {
		return "", fmt.Errorf("unknown private key. %s", pk)
	}

	pk = fmt.Sprintf("00%s", pk[2:])
	byteKey, err := hex.DecodeString(pk)
	if err != nil {
		return "", fmt.Errorf("private key decode err %v", err)
	}

	return base64.StdEncoding.EncodeToString(byteKey), nil
}

func B64PublicKeyToSuiAddress(b64 string) (string, error) {
	b64Decode, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("unknown base64. %s", b64)
	}
	addrBytes := blake2b.Sum256(b64Decode)
	return fmt.Sprintf("0x%s", hex.EncodeToString(addrBytes[:])[:64]), nil
}

func Ed25519PublicKeyToB64PublicKey(ed25519PubKey ed25519.PublicKey) string {
	newPubkey := []byte{0}
	newPubkey = append(newPubkey, ed25519PubKey...)
	return base64.StdEncoding.EncodeToString(newPubkey)
}

func uint16ToBuffer(num uint16) (*bytes.Buffer, error) {
	if num > 65535 {
		return nil, fmt.Errorf("invalid uint16 [%d]", num)
	}

	numBuffer := &bytes.Buffer{}
	err := binary.Write(numBuffer, binary.LittleEndian, num)

	return numBuffer, err
}

type IntentScope = uint8

const (
	TransactionDataIntentScope IntentScope = 0
	PersonalMessageIntentScope IntentScope = 3
)

func NewSuiMessageWithIntent(message []byte, scope IntentScope) []byte {
	intent := []byte{scope, 0, 0}
	intentMessage := make([]byte, len(intent)+len(message))
	copy(intentMessage, intent)
	copy(intentMessage[len(intent):], message)
	return intentMessage
}

func ToSerializedSignature(signature, pubKey []byte) string {
	signatureLen := len(signature)
	pubKeyLen := len(pubKey)
	serializedSignature := make([]byte, 1+signatureLen+pubKeyLen)
	serializedSignature[0] = byte(0x00)
	copy(serializedSignature[1:], signature)
	copy(serializedSignature[1+signatureLen:], pubKey)
	return base64.StdEncoding.EncodeToString(serializedSignature)
}

func FromSerializedSignature(serializedSignature string) (*SignaturePubkeyPair, error) {
	_bytes, err := base64.StdEncoding.DecodeString(serializedSignature)
	if err != nil {
		return nil, err
	}
	signatureScheme := ParseSignatureScheme(_bytes[0])
	if strings.EqualFold(serializedSignature, "") {
		return nil, fmt.Errorf("multiSig is not supported")
	}

	signature := _bytes[1 : len(_bytes)-32]
	pubKeyBytes := _bytes[1+len(signature):]

	keyPair := &SignaturePubkeyPair{
		SignatureScheme: signatureScheme,
		Signature:       signature,
		PubKey:          pubKeyBytes,
	}
	return keyPair, nil
}

func ParseSignatureScheme(scheme byte) string {
	switch scheme {
	case 0:
		return "ED25519"
	case 1:
		return "Secp256k1"
	case 2:
		return "Secp256r1"
	case 3:
		return "MultiSig"
	default:
		return "ED25519"
	}
}

func NormalizeShortAddress(address string) string {
	return fmt.Sprintf("0x%s", strings.TrimLeft(address, "0x"))
}

func NormalizeShortCoinType(coinType string) string {
	types := strings.Split(coinType, "::")
	if len(types) != 3 {
		return coinType
	}

	return fmt.Sprintf("%s::%s::%s", NormalizeShortAddress(types[0]), types[1], types[2])
}
