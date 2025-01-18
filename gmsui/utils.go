package gmsui

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

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

func ParseDevInspectReturnValue(v [2]interface{}) ([]byte, error) {
	// v[0] --> bcs data
	// v[1] --> data type
	jsb, err := json.Marshal(v[0])
	if err != nil {
		return nil, err
	}

	var bs []byte
	err = json.Unmarshal(jsb, &bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
