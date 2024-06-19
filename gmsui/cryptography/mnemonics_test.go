package cryptography

import (
	"testing"
)

func TestIsValidHardenedPath(t *testing.T) {
	validPaths := []string{
		"m/44'/784'/0'/0'/0'",
		"m/44'/784'/123'/456'/789'",
	}

	invalidPaths := []string{
		"m/44'/784'/0'/0/0'",
		"m/54'/784'/0'/0'/0'",
		"m/44'/784'/0'/0'/",
		"m/44/784'/0'/0'/0'",
		"m/44'/785'/0'/0'/0'",
	}

	for _, path := range validPaths {
		if !IsValidHardenedPath(path) {
			t.Errorf("Expected path %s to be valid", path)
		}
	}

	for _, path := range invalidPaths {
		if IsValidHardenedPath(path) {
			t.Errorf("Expected path %s to be invalid", path)
		}
	}
}

func TestIsValidBIP32Path(t *testing.T) {
	validPaths := []string{
		"m/54'/784'/0'/0/0",
		"m/74'/784'/123'/456/789",
	}

	invalidPaths := []string{
		"m/54'/784'/0'/0'/0'",
		"m/44'/784'/0'/0/0",
		"m/54'/784'/0'/0/",
		"m/54/784'/0'/0/0",
	}

	for _, path := range validPaths {
		if !IsValidBIP32Path(path) {
			t.Errorf("Expected path %s to be valid", path)
		}
	}

	for _, path := range invalidPaths {
		if IsValidBIP32Path(path) {
			t.Errorf("Expected path %s to be invalid", path)
		}
	}
}

func TestMnemonicToSeed(t *testing.T) {
	validMnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("failed to generate mnemonic: %v", err)
	}

	seed, err := MnemonicToSeed(validMnemonic)
	if err != nil {
		t.Errorf("expected mnemonic to be valid, got error: %v", err)
	}
	if len(seed) == 0 {
		t.Errorf("expected non-empty seed for valid mnemonic")
	}

	// invalidMnemonic := "invalid mnemonic phrase that does not conform to BIP39"
	// _, err = MnemonicToSeed(invalidMnemonic)
	// if err == nil {
	// 	t.Fatalf("mnemonic unable to seed, msg: %v", err)
	// }
}
