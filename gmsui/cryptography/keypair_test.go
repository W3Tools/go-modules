package cryptography

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestEncodeAndDecodeSuiPrivateKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate private key, msg: %v", err)
	}

	privateKeyBytes := privateKey.D.Bytes()
	if len(privateKeyBytes) != PrivateKeySize {
		t.Fatalf("expect private key size to be %d, got %d", privateKeyBytes, len(privateKeyBytes))
	}

	scheme := "ED25519"
	encoded, err := EncodeSuiPrivateKey(privateKeyBytes, scheme)
	if err != nil {
		t.Fatalf("failed to encode private key, msg: %v", err)
	}

	decoded, err := DecodeSuiPrivateKey(encoded)
	if err != nil {
		t.Fatalf("failed to decode sui private key, msg: %v", err)
	}

	if decoded.Scheme != scheme {
		t.Fatalf("expected scheme %s, got %s", scheme, decoded.Scheme)
	}

	if !bytes.Equal(decoded.SecretKey, privateKeyBytes) {
		t.Fatalf("decode private key does not match original")
	}
}
