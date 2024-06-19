package cryptography

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var pathRegex = regexp.MustCompile(`^m(\/[0-9]+')+$`)

const (
	Ed25519CURVE   = "ed25519 seed"
	HardenedOffset = 0x80000000
)

type Key struct {
	Key       []byte
	ChainCode []byte
}

func ReplaceDerive(val string) string {
	return strings.Replace(val, "'", "", -1)
}

func IsValidPath(path string) bool {
	if !pathRegex.MatchString(path) {
		return false
	}
	segments := strings.Split(path, "/")[1:]
	for _, seg := range segments {
		seg = ReplaceDerive(seg)
		if _, err := strconv.Atoi(seg); err != nil {
			return false
		}
	}
	return true
}

func GetMasterKeyFromSeed(seed string) (*Key, error) {
	bs, err := hex.DecodeString(seed)
	if err != nil {
		return nil, err
	}

	h := hmac.New(sha512.New, []byte(Ed25519CURVE))
	_, err = h.Write(bs)
	if err != nil {
		return nil, err
	}

	sum := h.Sum(nil)
	return &Key{Key: sum[:32], ChainCode: sum[32:]}, nil
}

func CKDPriv(parent *Key, index uint32) (*Key, error) {
	indexBuffer := make([]byte, 4)
	indexBuffer[0] = byte(index >> 24)
	indexBuffer[1] = byte(index >> 16)
	indexBuffer[2] = byte(index >> 8)
	indexBuffer[3] = byte(index)

	data := make([]byte, 1+len(parent.Key)+len(indexBuffer))
	data[0] = 0
	copy(data[1:], parent.Key)
	copy(data[1+len(parent.Key):], indexBuffer)

	h := hmac.New(sha512.New, parent.ChainCode)
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}

	sum := h.Sum(nil)
	return &Key{Key: sum[:32], ChainCode: sum[32:]}, nil
}

func DerivePath(path string, seed string) (*Key, error) {
	if !IsValidPath(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}

	key, err := GetMasterKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}

	segments := strings.Split(path, "/")[1:]
	for _, seg := range segments {
		seg = ReplaceDerive(seg)
		index, err := strconv.Atoi(seg)
		if err != nil {
			return nil, fmt.Errorf("invalid segment %s in path %s", seg, path)
		}

		key, err = CKDPriv(key, uint32(index)+HardenedOffset)
		if err != nil {
			return nil, err
		}
	}

	return key, nil
}
