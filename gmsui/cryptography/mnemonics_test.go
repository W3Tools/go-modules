package cryptography_test

import (
	"encoding/hex"
	"testing"

	"github.com/W3Tools/go-modules/gmsui/cryptography"
	"github.com/tyler-smith/go-bip39"
)

func TestIsValidHardenedPath(t *testing.T) {
	paths := []struct {
		input    string
		expected bool
	}{
		{"m/44'/784'/0'/0'/0'", true},
		{"m/44'/784'/123'/456'/789'", true},

		{"m/44'/784'/0'/0/0'", false},
		{"m/54'/784'/0'/0'/0'", false},
		{"m/44'/784'/0'/0'/", false},
		{"m/44/784'/0'/0'/0'", false},
		{"m/44'/785'/0'/0'/0'", false},
		{"m/44'/784'/0'/0'/0", false},        // Missing single quote at the end
		{"m/44'/784'/0'/0'/0'/", false},      // Extra slash at the end
		{"m/44'/784'/a'/b'/c'", false},       // Non-numeric indexes
		{"m/44'/784'/0'/0'/0'/extra", false}, // Extra characters at the end
	}

	for _, path := range paths {
		result := cryptography.IsValidHardenedPath(path.input)
		if result != path.expected {
			t.Errorf("IsValidHardenedPath(%s) returned %t, expected %t", path.input, result, path.expected)
		}
	}
}

func TestIsValidBIP32Path(t *testing.T) {
	paths := []struct {
		input    string
		expected bool
	}{
		{"m/54'/784'/0'/0/0", true},
		{"m/74'/784'/123'/456/789", true},
		{"m/54'/784'/0'/0'/0'", false},
		{"m/44'/784'/0'/0/0", false},
		{"m/54'/784'/0'/0/", false},
		{"m/54/784'/0'/0/0", false},

		{"m/54'/784'/0'/0/0'", false},      // Extra single quote at the end
		{"m/54'/784'/0'/0/0/", false},      // Extra slash at the end
		{"m/54'/784'/a'/b/c", false},       // Non-numeric indexes
		{"m/54'/784'/0'/0/0/extra", false}, // Extra characters at the end
		{"m/74'/784'/0'/0/0", true},        // Secp256r1 path
	}

	for _, path := range paths {
		result := cryptography.IsValidBIP32Path(path.input)
		if result != path.expected {
			t.Errorf("IsValidBIP32Path(%s) returned %t, expected %t", path.input, result, path.expected)
		}
	}
}

func TestMnemonicToSeed(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		t.Fatalf("failed to new entropy, msg: %v", err)
	}

	validMnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Fatalf("failed to generate mnemonic: %v", err)
	}

	seed, err := cryptography.MnemonicToSeed(validMnemonic)
	if err != nil {
		t.Errorf("expected mnemonic to be valid, got error: %v", err)
	}
	if len(seed) == 0 {
		t.Errorf("expected non-empty seed for valid mnemonic")
	}

	invalidMnemonic := "invalid mnemonic phrase that does not conform to BIP39"
	_, err = cryptography.MnemonicToSeed(invalidMnemonic)
	if err == nil {
		t.Fatalf("mnemonic unable to seed, msg: %v", err)
	}
}

func TestMnemonicToSeedHex(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		t.Fatalf("failed to new entropy, msg: %v", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Fatalf("GenerateMnemonic returned error: %v", err)
	}

	expectedSeed, err := cryptography.MnemonicToSeed(mnemonic)
	if err != nil {
		t.Fatalf("MnemonicToSeed returned error for valid mnemonic: %v", err)
	}

	expectedHex := hex.EncodeToString(expectedSeed)

	hexSeed, err := cryptography.MnemonicToSeedHex(mnemonic)
	if err != nil {
		t.Errorf("MnemonicToSeedHex returned error for valid mnemonic: %v", err)
	}
	if hexSeed != expectedHex {
		t.Errorf("MnemonicToSeedHex returned %s, expected %s", hexSeed, expectedHex)
	}

	invalidMnemonic := "invalid mnemonic"
	_, err = cryptography.MnemonicToSeedHex(invalidMnemonic)
	if err == nil {
		t.Errorf("MnemonicToSeedHex did not return error for invalid mnemonic: %s", invalidMnemonic)
	}
}

func TestGenerateMnemonic(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		t.Fatalf("failed to new entropy, msg: %v", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Errorf("GenerateMnemonic returned error: %v", err)
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		t.Errorf("GenerateMnemonic generated an invalid mnemonic: %s", mnemonic)
	}
}
