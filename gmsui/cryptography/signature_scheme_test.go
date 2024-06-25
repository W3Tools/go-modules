package cryptography

import (
	"testing"
)

func TestSignatureSchemeToFlag(t *testing.T) {
	tests := []struct {
		name          string
		scheme        SignatureScheme
		expectedFlag  SignatureFlag
		expectingFail bool
	}{
		{"Ed25519", Ed25519Scheme, 0x00, false},
		{"Secp256k1", Secp256k1Scheme, 0x01, false},
		{"Secp256r1", Secp256r1Scheme, 0x02, false},
		{"MultiSig", MultiSigScheme, 0x03, false},
		{"ZkLogin", ZkLoginScheme, 0x05, false},
		{"UnknownScheme", "UnknownScheme", 0x00, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag, exists := SignatureSchemeToFlag[tt.scheme]
			if tt.expectingFail {
				if exists {
					t.Errorf("expected failure, but got flag %v", flag)
				}
			} else {
				if !exists || flag != tt.expectedFlag {
					t.Errorf("expected flag %v, but got %v", tt.expectedFlag, flag)
				}
			}
		})
	}
}

func TestSignatureSchemeToSize(t *testing.T) {
	tests := []struct {
		name          string
		scheme        SignatureScheme
		expectedSize  int
		expectingFail bool
	}{
		{"Ed25519", Ed25519Scheme, 32, false},
		{"Secp256k1", Secp256k1Scheme, 33, false},
		{"Secp256r1", Secp256r1Scheme, 33, false},
		{"UnknownScheme", "UnknownScheme", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, exists := SignatureSchemeToSize[tt.scheme]
			if tt.expectingFail {
				if exists {
					t.Errorf("expected failure, but got size %v", size)
				}
			} else {
				if !exists || size != tt.expectedSize {
					t.Errorf("expected size %v, but got %v", tt.expectedSize, size)
				}
			}
		})
	}
}

func TestSignatureFlagToScheme(t *testing.T) {
	tests := []struct {
		name           string
		flag           SignatureFlag
		expectedScheme SignatureScheme
		expectingFail  bool
	}{
		{"Ed25519", 0x00, Ed25519Scheme, false},
		{"Secp256k1", 0x01, Secp256k1Scheme, false},
		{"Secp256r1", 0x02, Secp256r1Scheme, false},
		{"MultiSig", 0x03, MultiSigScheme, false},
		{"ZkLogin", 0x05, ZkLoginScheme, false},
		{"UnknownFlag", 0x04, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, exists := SignatureFlagToScheme[tt.flag]
			if tt.expectingFail {
				if exists {
					t.Errorf("expected failure, but got scheme %v", scheme)
				}
			} else {
				if !exists || scheme != tt.expectedScheme {
					t.Errorf("expected scheme %v, but got %v", tt.expectedScheme, scheme)
				}
			}
		})
	}
}
