package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/tyler-smith/go-bip39"
)

func TestGenerateAndVerifyEd25519Keypair(t *testing.T) {
	keypair, err := GenerateEd25519Keypair()
	if err != nil {
		t.Fatalf("unable to generate Ed25519 keypair, msg: %v", err)
	}

	if !reflect.DeepEqual(len(keypair.keypair.PublicKey), Ed25519PublicKeySize) {
		t.Errorf("expected public key size to be %d, but got %d", ed25519.PublicKeySize, len(keypair.keypair.PublicKey))
	}

	if !reflect.DeepEqual(len(keypair.keypair.SecretKey), ed25519.PrivateKeySize) {
		t.Errorf("expected private key size to be %d, but got %d", ed25519.PublicKeySize, len(keypair.keypair.PublicKey))
	}

	if !reflect.DeepEqual(keypair.GetKeyScheme(), cryptography.Ed25519Scheme) {
		t.Errorf("expected key scheme %v, but got %v", cryptography.Ed25519Scheme, keypair.GetKeyScheme())
	}

	publicKey, err := NewEd25519PublicKey(keypair.keypair.PublicKey)
	if err != nil {
		t.Fatalf("unable to create Ed25519 public key, msg: %v", err)
	}

	if !reflect.DeepEqual(publicKey.Flag(), cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme]) {
		t.Errorf("expected public key flag %v, but got %v", cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme], publicKey.Flag())
	}

	message := []byte("Hello, Go Modules!")

	t.Run("SignMessage", func(t *testing.T) {
		signature, _ := keypair.SignData(message)

		serializedSignature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{SignatureScheme: cryptography.Ed25519Scheme, PublicKey: publicKey, Signature: signature})
		if err != nil {
			t.Fatalf("unable to serialized signature, msg: %v", err)
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

func TestFromSecretKeyAndVerify(t *testing.T) {
	seed := make([]byte, ed25519.SeedSize)
	_, err := rand.Read(seed)
	if err != nil {
		t.Fatalf("error generating random secret key: %v", err)
	}

	keypair, err := FromSecretKey(seed, false)
	if err != nil {
		t.Fatalf("unable to create Ed25519 keypair from seed, msg: %v", err)
	}

	if !reflect.DeepEqual(len(keypair.keypair.PublicKey), ed25519.PublicKeySize) {
		t.Errorf("expected public key size to be %d, but got %d", ed25519.PublicKeySize, len(keypair.keypair.PublicKey))
	}

	if !reflect.DeepEqual(len(keypair.keypair.SecretKey), ed25519.PrivateKeySize) {
		t.Errorf("expected private key size to be %d, but got %d", ed25519.PublicKeySize, len(keypair.keypair.PublicKey))
	}

	if !reflect.DeepEqual(keypair.GetKeyScheme(), cryptography.Ed25519Scheme) {
		t.Errorf("expected key scheme %v, but got %v", cryptography.Ed25519Scheme, keypair.GetKeyScheme())
	}

	publicKey, err := NewEd25519PublicKey(keypair.keypair.PublicKey)
	if err != nil {
		t.Fatalf("unable to create Ed25519 public key, msg: %v", err)
	}

	if !reflect.DeepEqual(publicKey.Flag(), cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme]) {
		t.Errorf("expected public key flag %v, but got %v", cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme], publicKey.Flag())
	}

	message := []byte("Hello, Go Modules!")

	t.Run("SignMessage", func(t *testing.T) {
		signature, _ := keypair.SignData(message)

		serializedSignature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{SignatureScheme: cryptography.Ed25519Scheme, PublicKey: publicKey, Signature: signature})
		if err != nil {
			t.Fatalf("unable to serialized signature, msg: %v", err)
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

func TestDeriveKeypairFromMnemonic(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		t.Fatalf("failed to new entropy, msg: %v", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Fatalf("unable to generate mnemonic, msg: %v", err)
	}

	keypair, err := DeriveKeypair(mnemonic, DefaultEd25519DerivationPath)
	if err != nil {
		t.Fatalf("unable to derive keypair, msg: %v", err)
	}

	if !reflect.DeepEqual(len(keypair.keypair.PublicKey), ed25519.PublicKeySize) {
		t.Errorf("expected public key size to be %d, but got %d", ed25519.PublicKeySize, len(keypair.keypair.PublicKey))
	}

	if !reflect.DeepEqual(len(keypair.keypair.SecretKey), ed25519.PrivateKeySize) {
		t.Errorf("expected private key size to be %d, but got %d", ed25519.PublicKeySize, len(keypair.keypair.PublicKey))
	}

	if !reflect.DeepEqual(keypair.GetKeyScheme(), cryptography.Ed25519Scheme) {
		t.Errorf("expected key scheme %v, but got %v", cryptography.Ed25519Scheme, keypair.GetKeyScheme())
	}

	publicKey, err := NewEd25519PublicKey(keypair.keypair.PublicKey)
	if err != nil {
		t.Fatalf("unable to create Ed25519 public key, msg: %v", err)
	}

	if !reflect.DeepEqual(publicKey.Flag(), cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme]) {
		t.Errorf("expected public key flag %v, but got %v", cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme], publicKey.Flag())
	}

	message := []byte("Hello, Go Modules!")

	t.Run("SignMessage", func(t *testing.T) {
		signature, _ := keypair.SignData(message)

		serializedSignature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{SignatureScheme: cryptography.Ed25519Scheme, PublicKey: publicKey, Signature: signature})
		if err != nil {
			t.Fatalf("unable to serialized signature, msg: %v", err)
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
