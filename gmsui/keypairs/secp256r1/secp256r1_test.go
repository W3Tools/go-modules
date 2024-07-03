package secp256r1_test

import (
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/W3Tools/go-modules/gmsui/keypairs/secp256r1"
	"github.com/tyler-smith/go-bip39"
)

func TestGenerateAndVerifySecp256r1Keypair(t *testing.T) {
	keypair, err := secp256r1.GenerateSecp256r1Keypair()
	if err != nil {
		t.Fatalf("unable to generate Secp256r1 keypair, msg: %v", err)
	}

	if !reflect.DeepEqual(len(keypair.PublicKey()), cryptography.Secp256r1PublicKeySize) {
		t.Errorf("expected public key size to be %d, but got %d", cryptography.Secp256r1PublicKeySize, len(keypair.PublicKey()))
	}

	if !reflect.DeepEqual(len(keypair.SecretKey()), 32) {
		t.Errorf("expected private key size to be %d, but got %d", 32, len(keypair.SecretKey()))
	}

	if !reflect.DeepEqual(keypair.GetKeyScheme(), cryptography.Secp256r1Scheme) {
		t.Errorf("expected key scheme %v, but got %v", cryptography.Secp256r1Scheme, keypair.GetKeyScheme())
	}

	publicKey, err := keypair.GetPublicKey()
	if err != nil {
		t.Fatalf("unable to get Secp256r1 public key, msg: %v", err)
	}

	if !reflect.DeepEqual(publicKey.Flag(), cryptography.SignatureSchemeToFlag[cryptography.Secp256r1Scheme]) {
		t.Errorf("expected public key flag %v, but got %v", cryptography.SignatureSchemeToFlag[cryptography.Secp256r1Scheme], publicKey.Flag())
	}

	message := []byte("Hello, Go Modules!")

	t.Run("SignMessage", func(t *testing.T) {
		signature, _ := keypair.SignData(message)

		serializedSignature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{SignatureScheme: cryptography.Secp256r1Scheme, PublicKey: publicKey, Signature: signature})
		if err != nil {
			t.Fatalf("unable to serialize signature, msg: %v", err)
		}

		valid, err := publicKey.Verify(message, serializedSignature)
		if err != nil {
			t.Fatalf("unable to verify signature, msg: %v", err)
		}
		if !valid {
			t.Errorf("signature verification failed")
		}
	})

	t.Run("SignPersonalMessage", func(t *testing.T) {
		signature, err := keypair.SignPersonalMessage(message)
		if err != nil {
			t.Fatalf("unable to sign personal message, msg: %v", err)
		}

		valid, err := publicKey.VerifyPersonalMessage(message, signature.Signature)
		if err != nil {
			t.Fatalf("unable to verify personal message, msg: %v", err)
		}

		if !valid {
			t.Errorf("signature verification failed")
		}
	})
}

func TestFromSecretKeyAndVerifySecp256r1(t *testing.T) {
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		t.Fatalf("error generating random secret key: %v", err)
	}

	keypair, err := secp256r1.FromSecretKey(secretKey, false)
	if err != nil {
		t.Fatalf("unable to create Secp256r1 keypair from secret key, msg: %v", err)
	}

	if !reflect.DeepEqual(len(keypair.PublicKey()), 33) {
		t.Errorf("expected public key size to be %d, but got %d", 33, len(keypair.PublicKey()))
	}

	if !reflect.DeepEqual(len(keypair.SecretKey()), 32) {
		t.Errorf("expected private key size to be %d, but got %d", 32, len(keypair.SecretKey()))
	}

	if !reflect.DeepEqual(keypair.GetKeyScheme(), cryptography.Secp256r1Scheme) {
		t.Errorf("expected key scheme %v, but got %v", cryptography.Secp256r1Scheme, keypair.GetKeyScheme())
	}

	publicKey, err := keypair.GetPublicKey()
	if err != nil {
		t.Fatalf("unable to get Secp256r1 public key, msg: %v", err)
	}

	if !reflect.DeepEqual(publicKey.Flag(), cryptography.SignatureSchemeToFlag[cryptography.Secp256r1Scheme]) {
		t.Errorf("expected public key flag %v, but got %v", cryptography.SignatureSchemeToFlag[cryptography.Secp256r1Scheme], publicKey.Flag())
	}

	message := []byte("Hello, Go Modules!")

	t.Run("SignMessage", func(t *testing.T) {
		signature, _ := keypair.SignData(message)

		serializedSignature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{SignatureScheme: cryptography.Secp256r1Scheme, PublicKey: publicKey, Signature: signature})
		if err != nil {
			t.Fatalf("unable to serialize signature, msg: %v", err)
		}

		valid, err := publicKey.Verify(message, serializedSignature)
		if err != nil {
			t.Fatalf("unable to verify signature, msg: %v", err)
		}
		if !valid {
			t.Errorf("signature verification failed")
		}
	})
}

func TestDeriveSecp256r1KeypairFromMnemonic(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		t.Fatalf("failed to new entropy, msg: %v", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Fatalf("unable to generate mnemonic, msg: %v", err)
	}

	keypair, err := secp256r1.DeriveKeypair(mnemonic, secp256r1.DefaultSecp256r1DerivationPath)
	if err != nil {
		t.Fatalf("unable to derive Secp256r1 keypair, msg: %v", err)
	}

	if !reflect.DeepEqual(len(keypair.PublicKey()), 33) {
		t.Errorf("expected public key size to be %d, but got %d", 33, len(keypair.PublicKey()))
	}

	if !reflect.DeepEqual(len(keypair.SecretKey()), 32) {
		t.Errorf("expected private key size to be %d, but got %d", 32, len(keypair.SecretKey()))
	}

	if !reflect.DeepEqual(keypair.GetKeyScheme(), cryptography.Secp256r1Scheme) {
		t.Errorf("expected key scheme %v, but got %v", cryptography.Secp256r1Scheme, keypair.GetKeyScheme())
	}

	publicKey, err := keypair.GetPublicKey()
	if err != nil {
		t.Fatalf("unable to get Secp256r1 public key, msg: %v", err)
	}

	if !reflect.DeepEqual(publicKey.Flag(), cryptography.SignatureSchemeToFlag[cryptography.Secp256r1Scheme]) {
		t.Errorf("expected public key flag %v, but got %v", cryptography.SignatureSchemeToFlag[cryptography.Secp256r1Scheme], publicKey.Flag())
	}

	message := []byte("Hello, Go Modules!")

	t.Run("SignMessage", func(t *testing.T) {
		signature, _ := keypair.SignData(message)

		serializedSignature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{SignatureScheme: cryptography.Secp256r1Scheme, PublicKey: publicKey, Signature: signature})
		if err != nil {
			t.Fatalf("unable to serialize signature, msg: %v", err)
		}

		valid, err := publicKey.Verify(message, serializedSignature)
		if err != nil {
			t.Fatalf("unable to verify signature, msg: %v", err)
		}
		if !valid {
			t.Errorf("signature verification failed")
		}
	})
}
